// Copyright 2023 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build !linux && !plan9 && !windows

package main

import (
	"fmt"
	"syscall"

	"golang.org/x/sys/unix"
)

func sync(data, filesystem bool, args []string) error {
	switch {
	case data:
		return doSyscall(unix.SYS_FDATASYNC, args)
	case filesystem:
		return fmt.Errorf("-f is supported only on linux")
	default:
		syscall.Sync()
		return nil
	}
}
