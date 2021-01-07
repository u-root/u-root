// Copyright 2020 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"os/exec"
	"syscall"
	"unsafe"
)

func exitWithStatus(err *exec.ExitError) {
	// DAMN. os.Exit ABI is an int :-(
	// This does not play nice with Plan 9.
	cp := unsafe.Pointer(&[]byte(err.Error())[0])
	syscall.Syscall(8, uintptr(cp), 0, 0)
}

func lowpriority() error {
	// TODO: write pri 0 to /proc/<me>/ctl
	return nil
}
