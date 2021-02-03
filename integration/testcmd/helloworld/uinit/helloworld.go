// Copyright 2021 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"log"

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

	fmt.Println("HELLO WORLD")
	return nil
}

// The most trivial init script.
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
