// Copyright 2015-2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package kexec

import (
	"fmt"
	"log"
	"os"
	"syscall"

	"github.com/u-root/u-root/pkg/watchdogd"
)

// Reboot executes a kernel previously loaded with FileInit.
func Reboot() error {
	// Optionally disarm the watchdog.
	if os.Getenv("UROOT_KEXEC_DISARM_WATCHDOG") == "1" {
		d, err := watchdogd.New()
		if err != nil {
			log.Printf("Error dialing watchdog daemon: %v", err)
		} else if err := d.Disarm(); err != nil {
			log.Printf("Error disarming watchdog: %v", err)
		}
	}

	if err := syscall.Reboot(syscall.LINUX_REBOOT_CMD_KEXEC); err != nil {
		return fmt.Errorf("sys_reboot(..., kexec) = %w", err)
	}
	return nil
}
