// Copyright 2015-2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// On arm, not only does the Linux kernel not support kexec_file_load(2), it is
// not included in Go's unix package. !arm prevents a compile error.
//+build linux,!arm

package kexec

import (
	"errors"
	"fmt"

	"golang.org/x/sys/unix"
)

// KernelFileLoad uses the kernel's kexec_file_load(2) syscall to parse and
// load the kernel. Where supported, this is very convenient because Linux does
// nearly all the work for us.
func KernelFileLoad(opts *LinuxOpts) error {
	if opts.DTB != nil {
		return &ProbeError{
			errors.New("kexec_file_load(2) does not support a modified device tree"),
		}
	}

	var flags int
	var ramfsFd int
	if opts.Initramfs != nil {
		ramfsFd = int(opts.Initramfs.Fd())
	} else {
		flags |= unix.KEXEC_FILE_NO_INITRAMFS
	}

	if opts.FileLoadSyscall == nil {
		opts.FileLoadSyscall = RawFileLoad
	}
	err := opts.FileLoadSyscall(int(opts.Kernel.Fd()), ramfsFd, opts.CmdLine, flags)
	if errno, ok := err.(*unix.Errno); ok && *errno == unix.ENOSYS {
		// The kernel does not have a kexec_file_load syscall.
		return &ProbeError{err}
	}
	return nil
}

// RawFileLoad is a wrapper around the kexec_file_load(2) syscall.
func RawFileLoad(kernelFd int, initrd int, cmdline string, flags int) error {
	if err := unix.KexecFileLoad(kernelFd, initrd, cmdline, flags); err != nil {
		return fmt.Errorf("sys_kexec(%d, %d, %s, %x) = %v",
			kernelFd, initrd, cmdline, flags, err)
	}
	return nil
}
