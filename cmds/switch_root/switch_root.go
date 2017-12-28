// Copyright 2012-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"fmt"
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

// littleDoctor recursively deletes everything at "path" that
// is in the same file system as "fs".
func littleDoctor(path string, fsType int64) error {
	var pathFS syscall.Statfs_t

	if err := syscall.Statfs(path, &pathFS); err != nil {
		return err
	}

	if int64(pathFS.Type) != fsType {
		return nil
	}

	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("Could not open %s: %v", path, err)
	}

	if fileStat, err := file.Stat(); fileStat.IsDir() {
		if err != nil {
			return err
		}

		names, err := file.Readdirnames(-1)
		defer file.Close()
		if err != nil {
			return err
		}

		for _, fileName := range names {

			if fileName == "." || fileName == ".." {
				continue
			}

			if err := littleDoctor(filepath.Join(path, fileName), fsType); err != nil {
				return err
			}
			if err := os.Remove(path); err != nil {
				return err
			}
		}

	} else {
		if err := os.Remove(path); err != nil {
			return err
		}
	}

	return nil
}

// execCommand will run the executable at "path" with PID 1.
// It returns an error if the command exits incorrectly.
func execCommand(path string) error {
	cmd := exec.Command(path)
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setctty:    true,
		Setsid:     true,
		Cloneflags: syscall.CLONE_THREAD | syscall.CLONE_NEWPID,
	}

	return cmd.Run()
}

// specialFS creates and mounts proc, sys and dev at the root level.
func specialFS() error {

	if err := syscall.Mkdir("/proc", 0); err != nil {
		return err
	}

	if err := syscall.Mkdir("/sys", 0); err != nil {
		return err
	}

	if err := syscall.Mkdir("/dev", 0); err != nil {
		return err
	}

	if err := syscall.Mount("", "/proc", "proc", 0, ""); err != nil {
		return err
	}

	if err := syscall.Mount("", "/sys", "sysfs", 0, ""); err != nil {
		return err
	}

	return syscall.Mount("", "/dev", "devtmpfs", 0, "")
}

// switchRoot will recursive deletes current root, switches the current root to
// the "newRoot", creates special filesystems (proc, sys and dev) in the new root
// and execs "init"
func switchRoot(newRoot string, init string) error {
	log.Printf("switch_root: Changing directory")
	var rootFS syscall.Statfs_t

	if err := syscall.Chdir(newRoot); err != nil {
		return fmt.Errorf("switch_root: failed change directory to new_root %v", err)
	}

	if err := syscall.Statfs("/", &rootFS); err != nil {
		return fmt.Errorf("switch_root: failed statfs %v", err)
	}

	if err := littleDoctor("/", int64(rootFS.Type)); err != nil {
		return fmt.Errorf("switch_root: failed Deletion of rootfs %v", err)
	}

	log.Printf("switch_root: Overmounting on /")

	if err := syscall.Mount(".", "/", "ext4", syscall.MS_MOVE, ""); err != nil {
		return fmt.Errorf("switch_root: fatal mount error %v", err)
	}

	log.Printf("switch_root: Changing root!")

	if err := syscall.Chroot("."); err != nil {
		return fmt.Errorf("switch_root: fatal chroot error %v", err)
	}

	log.Printf("switch_root: returning to slash")
	if err := syscall.Chdir("/"); err != nil {
		return fmt.Errorf("switch_root: failed change directory to '/' %v", err)
	}

	log.Printf("switch_root: creating proc, dev and sys")

	if err := specialFS(); err != nil {
		return fmt.Errorf("switch_root: failed to create special files %v", err)
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
