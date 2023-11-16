// Copyright 2023 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package main

import (
	"syscall"

	"golang.org/x/sys/unix"
)

func sync(data, filesystem bool, args []string) error {
	switch {
	case data:
		return doSyscall(unix.SYS_FDATASYNC, args)
	case filesystem:
		return doSyscall(unix.SYS_SYNCFS, args)
	default:
		syscall.Sync()
		return nil
	}
}
