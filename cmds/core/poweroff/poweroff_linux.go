// Copyright 2022 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// poweroff turns the system off, without delay. There are no options.
//
// Synopsis:
//
//	poweroff
//
// Description:
//
//	poweroff calls the kernel to power off the systems.
package main

import (
	"log"

	"golang.org/x/sys/unix"
)

func main() {
	if err := unix.Reboot(unix.LINUX_REBOOT_CMD_POWER_OFF); err != nil {
		log.Fatal(err)
	}
}
