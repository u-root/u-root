// Copyright 2015-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build !amd64 && !arm64 && !riscv64

package kexec

import (
	"os"
	"syscall"
)

// FileLoad is not implemented for platforms other than amd64, arm64 and riscv64.
func FileLoad(kernel, ramfs *os.File, cmdline string) error {
	return syscall.ENOSYS
}
