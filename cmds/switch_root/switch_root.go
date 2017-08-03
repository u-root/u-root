// Copyright 2012-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// This is a basic init script.
package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
	"syscall"

	"github.com/u-root/u-root/uroot"
)

var (
	commands = []string{
		"/bin/bash",
	}

	DEFAULT_ROOT_DEV = "/dev/mmcblk0p1"
	DEVICE_PARAM     = "uroot.rootdevice"
)

func usage() string {
	// Return the usage string
	return ""
}

func littleDoctor(path string) {
	// Recursively deletes everything at slash
	// Does not continue down other filesystems i.e.
	// new_root, devtmpfs, profs and sysfs
}

func exec_command(path string) error {
	// Will exec and dup a command at path
	cmd := exec.Command(path)
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	var fd int
	cmd.SysProcAttr = &syscall.SysProcAttr{Ctty: fd, Setctty: true, Setsid: true, Cloneflags: uintptr(0)}
	log.Printf("Run %v", cmd)

	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}

func getCmdline() (string, error) {
	// Reads the kernel cmdline at "/proc/cmdline to a string
	var contents []byte

	contents, err := ioutil.ReadFile("/proc/cmdline")

	if err != nil {
		return "", fmt.Errorf("Read command line failed: %v", err)
	}

	return string(contents), nil
}

func getDevice() (string, error) {
	// Given the kernel command line, this will select a device to mount
	var cmdline string

	cmdline, err := getCmdline()

	if err != nil {
		log.Printf("Could not get kernel cmdline")
		return "", err
	}

	for _, item := range strings.Split(cmdline, " ") {
		paramValue := strings.SplitN(item, "=", 2)
		if paramValue[0] == DEVICE_PARAM {
			return paramValue[1], nil
		}
	}

	return "", fmt.Errorf("no device specified")
}

func start() {
	// This getpid adds a bit of cost to each invocation (not much really)
	// but it allows us to merge init and sh. The 600K we save is worth it.
	// Figure out which init to run. We must always do this.

	// log.Printf("init: os is %v, initMap %v", filepath.Base(os.Args[0]), initMap)
	// we use filepath.Base in case they type something like ./cmd

	log.Printf("switch_root: Making mount directory")

	if err := syscall.Mkdir("/mnt", 0777); err != nil {
		log.Printf("init: error %v", err)
	}

	log.Printf("switch_root: Mounting filesystem")

	if rootDevice, err := getDevice(); err != nil {
		log.Printf("Using Device Default %v", DEFAULT_ROOT_DEV)
		rootDevice = DEFAULT_ROOT_DEV
	} else {
		log.Printf("Using Device %v", rootDevice)
	}

	if err := syscall.Mount("/dev/mmcblk0p1", "/mnt", "ext4", 0, ""); err != nil {
		log.Fatalf("init: fatal mount error %v", err)
	}

	// Copy everything over to mnt/

	log.Printf("switch_root: Excing bash")

	if err := exec_command("/bin/bash-static"); err != nil {
		log.Printf("switch_root: exit ramfs")
	}

	log.Printf("switch_root: Changing directory")

	syscall.Chdir("/mnt")

	log.Printf("switch_root: Overmounting on /")

	if err := syscall.Mount(".", "/", "ext4", syscall.MS_MOVE, ""); err != nil {
		log.Fatalf("switch_root: fatal mount error %v", err)
	}

	log.Printf("switch_root: Changing root!")

	if err := syscall.Chroot("."); err != nil {
		log.Fatalf("switch_root: fatal chroot error %v", err)
	}

	log.Printf("unit: returning to slash")
	syscall.Chdir("/")

	log.Printf("unit: creating Uroot filesystem")

	log.Printf("Exec init!")
	uroot.Rootfs()

	if err := exec_command("/bin/bash"); err != nil {
		log.Printf("switch_root: returning to ramfs")
	}

}

func main() {

	start()
	log.Printf("switch_root failed")

}
