// Copyright 2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"log"
	"os"

	"golang.org/x/sys/unix"

	"github.com/u-root/u-root/pkg/cmdline"
	"github.com/u-root/u-root/pkg/sh"
)

// Mount a vfat volume and kexec the kernel within.
func main() {
	if err := os.MkdirAll("/testdata", 0755); err != nil {
		log.Fatal(err)
	}
	if os.Getenv("UROOT_USE_9P") == "1" {
		sh.RunOrDie("mount", "-t", "9p", "tmpdir", "/testdata")
	} else {
		sh.RunOrDie("mount", "-r", "-t", "vfat", "/dev/sda1", "/testdata")
	}

	// Get and increment the counter.
	kExecCounter, ok := cmdline.Flag("kexeccounter")
	if !ok {
		kExecCounter = "0"
	}
	fmt.Printf("KEXECCOUNTER=%s\n", kExecCounter)

	if kExecCounter == "0" {
		cmdLine := cmdline.FullCmdLine() + " kexeccounter=1"
		log.Print("cmdline: ", cmdLine)
		sh.RunOrDie("kexec",
			"-i", "/testdata/initramfs.cpio",
			"-c", cmdLine,
			"/testdata/kernel")
	} else {
		unix.Reboot(unix.LINUX_REBOOT_CMD_POWER_OFF)
	}
}
