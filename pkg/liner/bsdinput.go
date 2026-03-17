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

//go:build openbsd || freebsd || netbsd
// +build openbsd freebsd netbsd

package liner

import "syscall"

const (
	getTermios = syscall.TIOCGETA
	setTermios = syscall.TIOCSETA
)

const (
	// Input flags
	inpck  = 0x010
	istrip = 0x020
	icrnl  = 0x100
	ixon   = 0x200

	// Output flags
	opost = 0x1

	// Control flags
	cs8 = 0x300

	// Local flags
	isig   = 0x080
	icanon = 0x100
	iexten = 0x400
)

type termios struct {
	Iflag  uint32
	Oflag  uint32
	Cflag  uint32
	Lflag  uint32
	Cc     [20]byte
	Ispeed int32
	Ospeed int32
}

const cursorColumn = false
