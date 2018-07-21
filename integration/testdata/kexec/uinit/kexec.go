// Copyright 2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/u-root/u-root/pkg/cmdline"
)

func sh(arg0 string, args ...string) {
	cmd := exec.Command(arg0, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}
}

// Mount a vfat volume and kexec the kernel within.
func main() {
	sh("mkdir", "/testdata")
	sh("mount", "-r", "-t", "vfat", "/dev/sda1", "/testdata")

	// Get and increment the counter.
	kExecCounter, ok := cmdline.Flag("kexeccounter")
	if !ok {
		kExecCounter = "0"
	}
	fmt.Printf("KEXECCOUNTER=%s\n", kExecCounter)

	if kExecCounter == "0" {
		sh("kexec", "/testdata/bzImage",
			"-i", "/testdata/initramfs.cpio",
			"-cmdline", cmdline.FullCmdLine()+" kexeccounter=1")
	}
}
