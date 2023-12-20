// Copyright 2021 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"log"

	"github.com/hugelgupf/vmtest/guest"
	"golang.org/x/sys/unix"
)

func runTest() {
	defer guest.CollectKernelCoverage()

	fmt.Println("HELLO WORLD")
}

// The most trivial init script.
func main() {
	runTest()

	if err := unix.Reboot(unix.LINUX_REBOOT_CMD_POWER_OFF); err != nil {
		log.Fatalf("Failed to reboot: %v", err)
	}
}
