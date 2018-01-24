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

// check for the symlink bit
func isSymlink(mode uint32) bool {
	// numbers come from `man fstatat`
	return (mode & 0170000) == 0120000
}

// recursively deletes a directory and everything in it
// works with a file descriptor, can delete files not referenceable
// any more (e.g. after chroot)
// does not descent into mounts
func recursiveDelete(fd int) error {
	var rb syscall.Stat_t

	if err := syscall.Fstat(fd, &rb); err != nil {
		return err
	}

	// may need to call syscall.ReadDirent multiple times
	for {

		// get the subfiles in buff
		buff := make([]byte, 4096)
		var nbuff int
		nbuff, err := syscall.ReadDirent(fd, buff)
		if err != nil {
			return err
		}
		if nbuff <= 0 {
			break
		}

		_, _, names := syscall.ParseDirent(buff, nbuff, []string{})

		for _, name := range names {
			// try to open it with O_DIRECTORY to know if it's a directory
			// NOTE O_DIRECTORY is linux specific
			namefd, err := syscall.Openat(fd, name, syscall.O_DIRECTORY, syscall.O_RDWR)
			if err != nil {
				// just delete files
				if err := unix.Unlinkat(fd, name, 0); err != nil {
					return err
				}
			} else {
				defer syscall.Close(namefd)

				var nameStatT unix.Stat_t
				if err := unix.Fstatat(fd, name, &nameStatT, unix.AT_SYMLINK_NOFOLLOW); err != nil {
					return err
				}

				if fdDev, err := getDev(fd); err != nil {
					return err
				} else if fdDev != nameStatT.Dev {
					// on different filesystem, don't go there
					continue
				}

				// not actually a dir, but a symlink to a dir, treat like file and remove
				if isSymlink(nameStatT.Mode) {
					if err := unix.Unlinkat(fd, name, 0); err != nil {
						return err
					}
					continue
				}

				// actual recursion
				if err := recursiveDelete(namefd); err != nil {
					return err
				}
				// back from recursion, dir is not empty, delete
				if err := unix.Unlinkat(fd, name, unix.AT_REMOVEDIR); err != nil {
					return err
				}
			}
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
