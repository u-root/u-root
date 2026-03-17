// SPDX-License-Identifier: MIT

// Copyright © 2012 Peter Harris
//
// Permission is hereby granted, free of charge, to any person obtaining a
// copy of this software and associated documentation files (the "Software"),
// to deal in the Software without restriction, including without limitation
// the rights to use, copy, modify, merge, publish, distribute, sublicense,
// and/or sell copies of the Software, and to permit persons to whom the
// Software is furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice (including the next
// paragraph) shall be included in all copies or substantial portions of the
// Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.  IN NO EVENT SHALL
// THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING
// FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER
// DEALINGS IN THE SOFTWARE.

//go:build linux || darwin || freebsd || openbsd || netbsd
// +build linux darwin freebsd openbsd netbsd

package liner

import (
	"syscall"
	"unsafe"
)

func (mode *termios) ApplyMode() error {
	_, _, errno := syscall.Syscall(syscall.SYS_IOCTL, uintptr(syscall.Stdin), setTermios, uintptr(unsafe.Pointer(mode)))

	if errno != 0 {
		return errno
	}
	return nil
}

// TerminalMode returns the current terminal input mode as an InputModeSetter.
//
// This function is provided for convenience, and should
// not be necessary for most users of liner.
func TerminalMode() (ModeApplier, error) {
	return getMode(syscall.Stdin)
}

func getMode(handle int) (*termios, error) {
	var mode termios
	var err error
	_, _, errno := syscall.Syscall(syscall.SYS_IOCTL, uintptr(handle), getTermios, uintptr(unsafe.Pointer(&mode)))
	if errno != 0 {
		err = errno
	}

	return &mode, err
}
