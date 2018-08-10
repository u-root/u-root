// Copyright 2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"log"

	"github.com/u-root/u-root/pkg/cmdline"
	"github.com/u-root/u-root/pkg/sh"
)

// Mount a vfat volume and kexec the kernel within.
func main() {
	sh.RunOrDie("mkdir", "/testdata")
	sh.RunOrDie("mount", "-r", "-t", "vfat", "/dev/sda1", "/testdata")

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
			"/testdata/bzImage")
	}
}
