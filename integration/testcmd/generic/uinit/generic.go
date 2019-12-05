// Copyright 2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/u-root/u-root/pkg/mount"
	"golang.org/x/sys/unix"
)

func main() {
	if err := os.MkdirAll("/testdata", 0755); err != nil {
		log.Fatalf("Couldn't create testdata: %v", err)
	}

	// Mount a vfat volume and run the tests within.
	var (
		mp  *mount.MountPoint
		err error
	)
	if os.Getenv("UROOT_USE_9P") == "1" {
		mp, err = mount.Mount("tmpdir", "/testdata", "9p", "", 0)
	} else {
		mp, err = mount.Mount("/dev/sda1", "/testdata", "vfat", "", unix.MS_RDONLY)
	}
	if err != nil {
		log.Fatalf("Failed to mount test directory: %v", err)
	}
	defer mp.Unmount(0) //nolint:errcheck

	// Run the test script test.elv
	test := filepath.Join("/testdata", "test.elv")
	if _, err := os.Stat(test); os.IsNotExist(err) {
		log.Fatalf("Could not find any test script to run.")
	}

	cmd := exec.Command("elvish", test)
	cmd.Stdin, cmd.Stdout, cmd.Stderr = os.Stdin, os.Stdout, os.Stderr

	if err := cmd.Run(); err != nil {
		log.Fatalf("test.elv ran unsuccessfully: %v", err)
	}

	if err := unix.Reboot(unix.LINUX_REBOOT_CMD_POWER_OFF); err != nil {
		log.Fatalf("Failed to reboot: %v", err)
	}
}
