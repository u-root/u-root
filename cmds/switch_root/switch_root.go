// Copyright 2012-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// This is a basic init script.
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
	help             = flag.Bool("h", false, "Help")
	version          = flag.Bool("V", false, "Version")
	DEFAULT_ROOT_DEV = "/dev/mmcblk0p1"
	DEVICE_PARAM     = "uroot.rootdevice"
)

func usage() string {
	// Return the usage string
	return "switch_root [-h] [-V]\nswitch_root newroot init"
}

func littleDoctor(path string, fs *syscall.Statfs_t) error {
	// Recursively deletes everything at slash
	// Does not continue down other filesystems i.e.
	// new_root, devtmpfs, profs and sysfs

	pathFS := syscall.Statfs_t{}

	if err := syscall.Statfs(path, &pathFS); err != nil {
		return err
	}

	if pathFS.Type != fs.Type {
		return nil
	}

	file, err := os.Open(path)

	if err != nil {
		return fmt.Errorf("Could not open %s: %v", path, err)
	}

	if fileStat, _ := file.Stat(); fileStat.IsDir() {

		names, err := file.Readdirnames(-1)
		file.Close()

		if err != nil {
			return err
		}

		for _, fileName := range names {

			if fileName == "." || fileName == ".." {
				return nil
			}

			littleDoctor(filepath.Join(path, fileName), fs)
			os.Remove(path)
		}

	} else {
		os.Remove(path)
	}

	return nil
}

func execCommand(path string) error {
	// Will exec and dup a command at path
	var fd int
	cmd := exec.Command(path)
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	cmd.SysProcAttr = &syscall.SysProcAttr{Ctty: fd, Setctty: true, Setsid: true, Cloneflags: uintptr(0)}
	log.Printf("Run %v", cmd)

	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}

func specialFS() {

	syscall.Mkdir("/path", 0)
	syscall.Mkdir("/sys", 0)
	syscall.Mkdir("/dev", 0)

	syscall.Mount("proc", "/proc", "proc", syscall.MS_MGC_VAL, "")
	syscall.Mount("sys", "/sys", "sysfs", syscall.MS_MGC_VAL, "")
	syscall.Mount("none", "/dev", "devtmpfs", syscall.MS_MGC_VAL, "")
}

func start(new_root string, init string) {
	// This getpid adds a bit of cost to each invocation (not much really)
	// but it allows us to merge init and sh. The 600K we save is worth it.
	// Figure out which init to run. We must always do this.

	// log.Printf("init: os is %v, initMap %v", filepath.Base(os.Args[0]), initMap)
	// we use filepath.Base in case they type something like ./cmd
	log.Printf("switch_root: Changing directory")

	syscall.Chdir(new_root)

	rootFS := syscall.Statfs_t{}

	if err := syscall.Statfs("/", &rootFS); err != nil {
		log.Fatalf("switch_root: failed Stat %v", err)
	}

	if err := littleDoctor("/", &rootFS); err != nil {
		log.Fatalf("switch_root: failed Deletion of rootfs: %v", err)
	}

	log.Printf("switch_root: Overmounting on /")

	if err := syscall.Mount(".", "/", "ext4", syscall.MS_MOVE, ""); err != nil {
		log.Fatalf("switch_root: fatal mount error %v", err)
	}

	log.Printf("switch_root: Changing root!")

	if err := syscall.Chroot("."); err != nil {
		log.Fatalf("switch_root: fatal chroot error %v", err)
	}

	log.Printf("switch_root: returning to slash")
	syscall.Chdir("/")

	log.Printf("switch_root: creating proc, dev and sys")

	specialFS()

	log.Printf("switch_root: executing init")
	if err := execCommand(init); err != nil {
		log.Printf("switch_root: returning to ramfs")
	}

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

	new_root := flag.Args()[0]
	init := flag.Args()[1]

	start(new_root, init)
	log.Printf("switch_root failed")

}
