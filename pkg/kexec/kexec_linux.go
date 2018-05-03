// Copyright 2015-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package kexec

import (
	"fmt"
	"syscall"
)

// Reboot executes a kernel previously loaded with FileInit.
func Reboot() error {
	if err := syscall.Reboot(syscall.LINUX_REBOOT_CMD_KEXEC); err != nil {
		return fmt.Errorf("sys_reboot(..., kexec) = %v", err)
	}
	return nil
}
