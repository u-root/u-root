// Copyright 2012-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// This is a basic init script.
package main

import (
	"log"
	"os"
	"os/exec"
	"strings"
)

var (
	commands = []string{
		"/bin/bash",
	}
)

func start() {
	flag.Usage = usage
	// This getpid adds a bit of cost to each invocation (not much really)
	// but it allows us to merge init and sh. The 600K we save is worth it.
	// Figure out which init to run. We must always do this.

	// log.Printf("init: os is %v, initMap %v", filepath.Base(os.Args[0]), initMap)
	// we use filepath.Base in case they type something like ./cmd

	log.Printf("init: Making mount directory")

	if err := syscall.Mkdir("/mnt", 0777); err != nil {
		log.Printf("init: error %v", err)
	}

	log.Printf("init: Mounting filesystem")

	if err := syscall.Mount("/dev/mmcblk0p1", "/mnt", "ext4", 0, ""); err != nil {
		log.Printf("init: error %v", err)
	}

	log.Printf("init: Changing directory")

	syscall.Chdir("/mnt")

	log.Printf("init: Overmounting on /")

	if err := syscall.Mount(".", "/", "ext4", syscall.MS_MOVE, ""); err != nil {
		log.Printf("init: error %v", err)
	}

	log.Printf("init: Changing root!")

	if err := syscall.Chroot("."); err != nil {
		log.Printf("Change root: error %v", err)
	}

	log.Printf("init: returning to slash")
	syscall.Chdir("/")

	log.Printf("Exec init!")

	cmd := exec.Command("/bin/bash")
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	var fd int
	cmd.SysProcAttr = &syscall.SysProcAttr{Ctty: fd, Setctty: true, Setsid: true, Cloneflags: uintptr(0)}
	log.Printf("Run %v", cmd)

	if err := cmd.Run(); err != nil {
		log.Print(err)
	}

}

func main() {

	start()

}
