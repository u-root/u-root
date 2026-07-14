// Copyright 2013-2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build !freebsd && !plan9 && !windows

package cpio

import (
	"syscall"
)

func mknod(path string, mode uint32, dev int) (err error) {
	return syscall.Mknod(path, mode, dev)
}
