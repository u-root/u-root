// Copyright 2015-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build linux,amd64 linux,arm64

package kexec

import (
	"fmt"
	"os"

	"golang.org/x/sys/unix"
)

// FileLoad loads the given kernel as the new kernel with the given ramfs and
// cmdline.
//
// The kexec_file_load(2) syscall is x86-64 and arm64 only.
func FileLoad(kernel, ramfs *os.File, cmdline string) error {
	var flags int
	var ramfsfd int
	if ramfs != nil {
		ramfsfd = int(ramfs.Fd())
	} else {
		flags |= unix.KEXEC_FILE_NO_INITRAMFS
	}

	if err := unix.KexecFileLoad(int(kernel.Fd()), ramfsfd, cmdline, flags); err != nil {
		return fmt.Errorf("sys_kexec(%d, %d, %s, %x) = %v", kernel.Fd(), ramfsfd, cmdline, flags, err)
	}
	return nil
}
