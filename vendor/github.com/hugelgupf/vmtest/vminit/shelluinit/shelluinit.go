// Copyright 2021 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Command shelluinit runs commands from an elvish script.
package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/hugelgupf/vmtest/guest"
	"golang.org/x/sys/unix"
)

func runTest() error {
	cleanup, err := guest.MountSharedDir()
	if err != nil {
		return err
	}
	defer cleanup()
	defer guest.CollectKernelCoverage()

	mp, err := guest.Mount9PDir("/shelltestdata", "shelltest")
	if err != nil {
		return err
	}
	defer func() { _ = mp.Unmount(0) }()

	// Run the test script test.elv
	test := "/shelltestdata/test.elv"
	if _, err := os.Stat(test); os.IsNotExist(err) {
		return errors.New("could not find any test script to run")
	}
	cmd := exec.Command("elvish", test)
	cmd.Stdin, cmd.Stdout, cmd.Stderr = os.Stdin, os.Stdout, os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("test.elv ran unsuccessfully: %v", err)
	}
	return nil
}

func main() {
	if err := runTest(); err != nil {
		log.Printf("Tests failed: %v", err)
	} else {
		log.Print("TESTS PASSED MARKER")
	}

	if err := unix.Reboot(unix.LINUX_REBOOT_CMD_POWER_OFF); err != nil {
		log.Fatalf("Failed to reboot: %v", err)
	}
}
