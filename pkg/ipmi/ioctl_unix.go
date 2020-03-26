// Copyright 2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Bits and pieces from asm-generic/ioctl.h
package ipmi

import (
	"syscall"
	"unsafe"
)

const (
	_IOC_NONE  = 0x0
	_IOC_WRITE = 0x1
	_IOC_READ  = 0x2

	_IOC_NRBITS   = 8
	_IOC_TYPEBITS = 8
	_IOC_SIZEBITS = 14
	_IOC_NRSHIFT  = 0

	_IOC_TYPESHIFT = _IOC_NRSHIFT + _IOC_NRBITS
	_IOC_SIZESHIFT = _IOC_TYPESHIFT + _IOC_TYPEBITS
	_IOC_DIRSHIFT  = _IOC_SIZESHIFT + _IOC_SIZEBITS
)

func ioc(dir int, t int, nr int, size int) int {
	return (dir << _IOC_DIRSHIFT) | (t << _IOC_TYPESHIFT) |
		(nr << _IOC_NRSHIFT) | (size << _IOC_SIZESHIFT)
}

func IO(t int, nr int) int {
	return ioc(_IOC_NONE, t, nr, 0)
}

func IOR(t int, nr int, size int) int {
	return ioc(_IOC_READ, t, nr, size)
}

func IOWR(t int, nr int, size int) int {
	return ioc(_IOC_READ|_IOC_WRITE, t, nr, size)
}

func Ioctl(fd uintptr, name int, data unsafe.Pointer) syscall.Errno {
	_, _, err := syscall.Syscall(syscall.SYS_IOCTL, fd, uintptr(name), uintptr(data))
	return err
}
