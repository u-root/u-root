// Copyright 2020 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"os"
	"syscall"
)

const (
	_IOC_NONE = 0x0

	_IOC_NRBITS   = 8
	_IOC_TYPEBITS = 8
	_IOC_SIZEBITS = 14
	_IOC_NRSHIFT  = 0

	_IOC_TYPESHIFT = _IOC_NRSHIFT + _IOC_NRBITS
	_IOC_SIZESHIFT = _IOC_TYPESHIFT + _IOC_TYPEBITS
	_IOC_DIRSHIFT  = _IOC_SIZESHIFT + _IOC_SIZEBITS

	_IOC_NVRAM = 'p'
)

func ioc(dir int, t int, nr int, size int) int {
	return (dir << _IOC_DIRSHIFT) | (t << _IOC_TYPESHIFT) |
		(nr << _IOC_NRSHIFT) | (size << _IOC_SIZESHIFT)
}

func io(t int, nr int) int {
	return ioc(_IOC_NONE, t, nr, 0)
}

// ioctl to clear whole CMOS/NV-RAM
var _NVRAM_INIT = io(_IOC_NVRAM, 0x40)

func cmosClear() error {
	f, err := os.OpenFile("/dev/nvram", os.O_RDWR, 0)
	if err != nil {
		return err
	}
	defer f.Close()

	if _, _, errno := syscall.Syscall(syscall.SYS_IOCTL, uintptr(f.Fd()), uintptr(_NVRAM_INIT), 0); errno != 0 {
		return errno
	}
	return nil
}

func reboot() error {
	return syscall.Reboot(syscall.LINUX_REBOOT_CMD_RESTART)
}
