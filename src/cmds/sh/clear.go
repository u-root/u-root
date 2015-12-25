// Copyright 2015 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// By Manoel Vilela <manoel_vilela@engineer.com>

package main

import (
	"fmt"
	"syscall"
	"unsafe"
)

func init() {
	addBuiltIn("clear", clear)
}

type winsize struct {
	Row    uint16
	Col    uint16
	Xpixel uint16
	Ypixel uint16
}

func getLines() uint {
	ws := &winsize{}
	retCode, _, errno := syscall.Syscall(
		syscall.SYS_IOCTL,
		uintptr(syscall.Stdin),
		uintptr(syscall.TIOCGWINSZ),
		uintptr(unsafe.Pointer(ws)),
	)

	if int(retCode) == -1 {
		panic(errno)
	}
	return uint(ws.Row)
}

func clear(c *Command) error {
	lines := getLines()
	var i uint
	for i = 0; i < lines; i++ {
		fmt.Printf("\n")
	}

	return nil
}
