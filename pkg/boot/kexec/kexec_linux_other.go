// Copyright 2015-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build linux,!amd64,!arm64

package kexec

import (
	"os"
	"syscall"
)

func FileLoad(kernel, ramfs *os.File, cmdline string) error {
	return syscall.ENOSYS
}
