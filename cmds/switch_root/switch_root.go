// Copyright 2012-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"fmt"
	"golang.org/x/sys/unix"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
)

var (
	help    = flag.Bool("h", false, "Help")
	version = flag.Bool("V", false, "Version")
)

func usage() string {
	return "switch_root [-h] [-V]\nswitch_root newroot init"
}

// get device for file descriptor
func getDev(fd int) (dev uint64, err error) {
	var stat syscall.Stat_t

	if err := syscall.Fstat(fd, &stat); err != nil {
		return 0, err
	}

	return stat.Dev, nil
}

// recursively deletes a directory and everything in it
// works with a file descriptor, can delete files not referenceable
// any more (e.g. after chroot)
// does not descent into mounts
func recursiveDelete(fd int) error {
	parentDev, err := getDev(fd)
	if err != nil {
		return err
	}

  // the file descriptor is already open, but allocating a os.File
	// here for it makes reading the files in the dir so much nicer
	dir := os.NewFile(uintptr(fd), "__ignored__") // filename is completely irrelevant
	defer dir.Close()
	names, err := dir.Readdirnames(-1)
	if err != nil {
		return err
	}

	for _, name := range names {
		// loop here, but handle loop in separate function to make defer work as expected
		if err := recusiveDeleteInner(fd, parentDev, name); err != nil {
			return err
		}
	}
	return nil
}

// don't call this directly, this function is the loop content of recursiveDelete
func recusiveDeleteInner(parentFd int, parentDev uint64, childName string) error {
	// O_DIRECTORY and O_NOFOLLOW make this open fail for all files and all symlinks (even when pointing to a dir)
	// we need to filter out symlinks because getDev later follows them
	childFd, err := syscall.Openat(parentFd, childName, syscall.O_DIRECTORY | syscall.O_NOFOLLOW, syscall.O_RDWR)
	if err != nil {
		// either file or symlink, delete in any case
		if err := unix.Unlinkat(parentFd, childName, 0); err != nil {
			return err
		}
	} else {
		// open succeeded, which means it is a real directory
		defer unix.Close(childFd)

		// don't descent into other file systems
		if childFdDev, err := getDev(childFd); err != nil {
			return err
		} else if childFdDev != parentDev {
			// this means continue in recursiveDelete
			return nil
		}

		if err:= recursiveDelete(childFd); err != nil {
			return err
		}
		// back from recursion, dir is now empty, delete
		if err := unix.Unlinkat(parentFd, childName, unix.AT_REMOVEDIR); err != nil {
			return err
		}
	}
	return nil
}

// do a proper exec to retain PID 1.
// It returns an error if the command exits incorrectly.
func execCommand(path string) error {
	return syscall.Exec(path, []string{path}, []string{})
}

// check if dir is empty
func isEmpty(name string) (bool, error) {
	f, err := os.Open(name)
	if err != nil {
		return false, err
	}
	defer f.Close()

	_, err = f.Readdirnames(1)
	if err == io.EOF { // this means emtpy
		return true, nil
	}
	return false, err
}

// moves a mountpoint
func moveMount(oldPath string, newPath string) error {
	return syscall.Mount(oldPath, newPath, "", syscall.MS_MOVE, "")
}

// moves common 'special' mounts if the exist
func specialFS(newRoot string) error {
	var mounts = []string{"/dev", "/proc", "/sys", "/run"}

	for _, mount := range mounts {
		path := filepath.Join(newRoot, mount)
		// skip all mounting if the dir does not exists
		if _, err := os.Stat(mount); os.IsNotExist(err) {
			fmt.Println("switch_root: Skipping", mount)
			continue
		} else if err != nil {
			return err
		}
		// make sure the target dir exist and is empty
		if _, err := os.Stat(path); os.IsNotExist(err) {
			if err := syscall.Mkdir(path, 0); err != nil {
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

// switchroot moves special mounts (dev, proc, sys, run) to the new dir,
// then does a chroot, moves the root mount to the new dir and finally
// DELETES EVERYTHING in the old dir and execs init
func switchRoot(newRoot string, init string) error {
	log.Printf("switch_root: moving mounts")
	if err := specialFS(newRoot); err != nil {
		return fmt.Errorf("switch_root: moving mounts failed %v", err)
	}

	log.Printf("switch_root: Changing directory")
	if err := syscall.Chdir(newRoot); err != nil {
		return fmt.Errorf("switch_root: failed change directory to new_root %v", err)
	}

	// Open "/" now, we need the file descriptor later
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
	if err := syscall.Chroot("."); err != nil {
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
