// Copyright 2021 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/u-root/u-root/integration/testcmd/common"
	"golang.org/x/sys/unix"
)

func runTest() error {
	cleanup, err := common.MountSharedDir()
	if err != nil {
		return err
	}
	defer cleanup()
	defer common.CollectKernelCoverage()

	// Run the test script test.elv
	test := "/testdata/test.elv"
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
		log.Print("Tests passed")
	}

	if err := unix.Reboot(unix.LINUX_REBOOT_CMD_POWER_OFF); err != nil {
		log.Fatalf("Failed to reboot: %v", err)
	}
}
