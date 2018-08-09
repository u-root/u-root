// Copyright 2012-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build linux

package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	"golang.org/x/sys/unix"
)

var (
	help    = flag.Bool("h", false, "Help")
	version = flag.Bool("V", false, "Version")
)

func usage() string {
	return "switch_root [-h] [-V]\nswitch_root newroot init"
}

// getDev returns the device (as returned by the FSTAT syscall) for the given file descriptor.
func getDev(fd int) (dev uint64, err error) {
	var stat unix.Stat_t

	if err := unix.Fstat(fd, &stat); err != nil {
		return 0, err
	}

	return stat.Dev, nil
}

// recursiveDelete deletes a directory identified by `fd` and everything in it.
//
// This function allows deleting directories no longer referenceable by
// any file name. This function does not descend into mounts.
func recursiveDelete(fd int) error {
	parentDev, err := getDev(fd)
	if err != nil {
		return err
	}

	// The file descriptor is already open, but allocating a os.File
	// here makes reading the files in the dir so much nicer.
	dir := os.NewFile(uintptr(fd), "__ignored__")
	defer dir.Close()
	names, err := dir.Readdirnames(-1)
	if err != nil {
		return err
	}

	for _, name := range names {
		// Loop here, but handle loop in separate function to make defer work as expected.
		if err := recusiveDeleteInner(fd, parentDev, name); err != nil {
			return err
		}
	}
	return nil
}

// recusiveDeleteInner is called from recursiveDelete and either deletes
// or recurses into the given file or directory
//
// There should be no need to call this function directly.
func recusiveDeleteInner(parentFd int, parentDev uint64, childName string) error {
	// O_DIRECTORY and O_NOFOLLOW make this open fail for all files and all symlinks (even when pointing to a dir).
	// We need to filter out symlinks because getDev later follows them.
	childFd, err := unix.Openat(parentFd, childName, unix.O_DIRECTORY|unix.O_NOFOLLOW, unix.O_RDWR)
	if err != nil {
		// childName points to either a file or a symlink, delete in any case.
		if err := unix.Unlinkat(parentFd, childName, 0); err != nil {
			return err
		}
	} else {
		// Open succeeded, which means childName points to a real directory.
		defer unix.Close(childFd)

		// Don't descent into other file systems.
		if childFdDev, err := getDev(childFd); err != nil {
			return err
		} else if childFdDev != parentDev {
			// This means continue in recursiveDelete.
			return nil
		}

		if err := recursiveDelete(childFd); err != nil {
			return err
		}
		// Back from recursion, the directory is now empty, delete.
		if err := unix.Unlinkat(parentFd, childName, unix.AT_REMOVEDIR); err != nil {
			return err
		}
	}
	return nil
}

// execCommand execs into the given command.
//
// In order to preserve whatever PID this program is running with,
// the implementation does an actual EXEC syscall without forking.
func execCommand(path string) error {
	return unix.Exec(path, []string{path}, []string{})
}

// isEmpty returns true if the directory with the given path is empty.
func isEmpty(name string) (bool, error) {
	f, err := os.Open(name)
	if err != nil {
		return false, err
	}
	defer f.Close()

	if _, err := f.Readdirnames(1); err == io.EOF {
		return true, nil
	}
	return false, err
}

// moveMount moves mount
//
// This function is just a wrapper around the MOUNT syscall with the
// MOVE flag supplied.
func moveMount(oldPath string, newPath string) error {
	return unix.Mount(oldPath, newPath, "", unix.MS_MOVE, "")
}

// specialFS moves the 'special' mounts to the given target path
//
// 'special' in this context refers to the following non-blockdevice backed
// mounts that are almost always used: /dev, /proc, /sys, and /run.
// This function will create the target directories, if necessary.
// If the target directories already exists, they must be empty.
// This function skips missing mounts.
func specialFS(newRoot string) error {
	var mounts = []string{"/dev", "/proc", "/sys", "/run"}

	for _, mount := range mounts {
		path := filepath.Join(newRoot, mount)
		// Skip all mounting if the directory does not exists.
		if _, err := os.Stat(mount); os.IsNotExist(err) {
			fmt.Println("switch_root: Skipping", mount)
			continue
		} else if err != nil {
			return err
		}
		// Make sure the target dir exists and is empty.
		if _, err := os.Stat(path); os.IsNotExist(err) {
			if err := unix.Mkdir(path, 0); err != nil {
				return err
			}
		} else if err != nil {
			return err
		}
		if empty, err := isEmpty(path); err != nil {
			return err
		} else if !empty {
			return fmt.Errorf("%v must be empty", path)
		}
		if err := moveMount(mount, path); err != nil {
			return err
		}
	}
	return nil
}

// switchroot moves special mounts (dev, proc, sys, run) to the new directory,
// then does a chroot, moves the root mount to the new directory and finally
// DELETES EVERYTHING in the old root and execs the given init.
func switchRoot(newRoot string, init string) error {
	log.Printf("switch_root: moving mounts")
	if err := specialFS(newRoot); err != nil {
		return fmt.Errorf("switch_root: moving mounts failed %v", err)
	}

	log.Printf("switch_root: Changing directory")
	if err := unix.Chdir(newRoot); err != nil {
		return fmt.Errorf("switch_root: failed change directory to new_root %v", err)
	}

	// Open "/" now, we need the file descriptor later.
	oldRoot, err := os.Open("/")
	if err != nil {
		return err
	}
	defer oldRoot.Close()

	log.Printf("switch_root: Moving /")
	if err := moveMount(newRoot, "/"); err != nil {
		return err
	}

	log.Printf("switch_root: Changing root!")
	if err := unix.Chroot("."); err != nil {
		return fmt.Errorf("switch_root: fatal chroot error %v", err)
	}

	log.Printf("switch_root: Deleting old /")
	if err := recursiveDelete(int(oldRoot.Fd())); err != nil {
		panic(err)
	}

	log.Printf("switch_root: executing init")
	if err := execCommand(init); err != nil {
		return fmt.Errorf("switch_root: exec failed %v", err)
	}

	return nil
}

func main() {
	flag.Parse()

	if len(flag.Args()) == 0 {
		fmt.Println(usage())
		os.Exit(0)
	}

	if *help {
		fmt.Println(usage())
		os.Exit(0)
	}

	if *version {
		fmt.Println("Version XX")
		os.Exit(0)
	}

	newRoot := flag.Args()[0]
	init := flag.Args()[1]

	if err := switchRoot(newRoot, init); err != nil {
		log.Fatalf("switch_root failed %v\n", err)
	}
}
