// Copyright 2026 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// reboot restarts the system, without delay. There are no options.
//
// Synopsis:
//
//	reboot
//
// Description:
//
//	reboot flushes filesystem buffers and restarts the system.
package main

import (
	"log"

	"golang.org/x/sys/unix"
)

func main() {
	unix.Sync()
	if err := unix.Reboot(unix.LINUX_REBOOT_CMD_RESTART); err != nil {
		log.Fatal(err)
	}
}
