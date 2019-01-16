// Copyright 2015-2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// On arm, not only does the Linux kernel not support kexec_file_load(2), it is
// not included in Go's unix package. !arm prevents a not-found compile error.
//+build linux,arm

package kexec

import (
	"errors"
)

// RawFileLoad is a wrapper around the kexec_file_load(2) syscall.
func RawFileLoad(kernelFd int, initrd int, cmdline string, flags int) error {
	return errors.New("kexec_file_load(2) not supported on arm")
}
