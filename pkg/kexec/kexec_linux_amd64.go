// Copyright 2015-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package kexec

import (
	"fmt"
	"os"
	"syscall"
	"unsafe"
)

// Syscall number for kexec_file_load(2).
const _SYS_KEXEC_FILE_LOAD = 320

// kexec_file_load(2) syscall flags.
const (
	_KEXEC_FILE_UNLOAD       = 0x1
	_KEXEC_FILE_ON_CRASH     = 0x2
	_KEXEC_FILE_NO_INITRAMFS = 0x4
)

// FileLoad loads the given kernel as the new kernel with the given ramfs and
// cmdline.
//
// The kexec_file_load(2) syscall is x86-64 bit only.
func FileLoad(kernel, ramfs *os.File, cmdline string) error {
	var flags uintptr
	var ramfsfd uintptr
	if ramfs != nil {
		ramfsfd = ramfs.Fd()
	} else {
		flags |= _KEXEC_FILE_NO_INITRAMFS
	}

	cmdPtr, err := syscall.BytePtrFromString(cmdline)
	if err != nil {
		return fmt.Errorf("could not use cmdline %q: %v", cmdline, err)
	}

	if _, _, errno := syscall.Syscall6(
		_SYS_KEXEC_FILE_LOAD,
		kernel.Fd(),
		ramfsfd,
		uintptr(len(cmdline)),
		uintptr(unsafe.Pointer(cmdPtr)),
		flags,
		0); errno != 0 {
		return fmt.Errorf("sys_kexec(%d, %d, %s, %x) = %v", kernel.Fd(), ramfsfd, cmdline, flags, errno)
	}
	return nil
}
