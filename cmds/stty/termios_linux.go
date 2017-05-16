// Copyright 2015-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"syscall"
	"unsafe"
)

type bit struct {
	word int
	mask uint32
}

var (
	boolFields = map[string]*bit{
		// Input processing
		"ignbrk":  &bit{word: I, mask: syscall.IGNBRK},
		"brkint":  &bit{word: I, mask: syscall.BRKINT},
		"ignpar":  &bit{word: I, mask: syscall.IGNPAR},
		"parmrk":  &bit{word: I, mask: syscall.PARMRK},
		"inpck":   &bit{word: I, mask: syscall.INPCK},
		"istrip":  &bit{word: I, mask: syscall.ISTRIP},
		"inlcr":   &bit{word: I, mask: syscall.INLCR},
		"igncr":   &bit{word: I, mask: syscall.IGNCR},
		"icrnl":   &bit{word: I, mask: syscall.ICRNL},
		"iuclc":   &bit{word: I, mask: syscall.IUCLC},
		"ixon":    &bit{word: I, mask: syscall.IXON},
		"ixany":   &bit{word: I, mask: syscall.IXANY},
		"ixoff":   &bit{word: I, mask: syscall.IXOFF},
		"imaxbel": &bit{word: I, mask: syscall.IMAXBEL},
		"iutf8":   &bit{word: I, mask: syscall.IUTF8},

		//Outputprocessing
		"opost":  &bit{word: O, mask: syscall.OPOST},
		"olcuc":  &bit{word: O, mask: syscall.OLCUC},
		"onlcr":  &bit{word: O, mask: syscall.ONLCR},
		"ocrnl":  &bit{word: O, mask: syscall.OCRNL},
		"onocr":  &bit{word: O, mask: syscall.ONOCR},
		"onlret": &bit{word: O, mask: syscall.ONLRET},
		"ofill":  &bit{word: O, mask: syscall.OFILL},
		"ofdel":  &bit{word: O, mask: syscall.OFDEL},

		//Localprocessing
		"isig":    &bit{word: L, mask: syscall.ISIG},
		"icanon":  &bit{word: L, mask: syscall.ICANON},
		"xcase":   &bit{word: L, mask: syscall.XCASE},
		"echo":    &bit{word: L, mask: syscall.ECHO},
		"echoe":   &bit{word: L, mask: syscall.ECHOE},
		"echok":   &bit{word: L, mask: syscall.ECHOK},
		"echonl":  &bit{word: L, mask: syscall.ECHONL},
		"noflsh":  &bit{word: L, mask: syscall.NOFLSH},
		"tostop":  &bit{word: L, mask: syscall.TOSTOP},
		"echoctl": &bit{word: L, mask: syscall.ECHOCTL},
		"echoprt": &bit{word: L, mask: syscall.ECHOPRT},
		"echoke":  &bit{word: L, mask: syscall.ECHOKE},
		"flusho":  &bit{word: L, mask: syscall.FLUSHO},
		"pendin":  &bit{word: L, mask: syscall.PENDIN},
		"iexten":  &bit{word: L, mask: syscall.IEXTEN},

		//Controlprocessing

		"cstopb": &bit{word: C, mask: syscall.CSTOPB},
		"cread":  &bit{word: C, mask: syscall.CREAD},
		"parenb": &bit{word: C, mask: syscall.PARENB},
		"parodd": &bit{word: C, mask: syscall.PARODD},
		"hupcl":  &bit{word: C, mask: syscall.HUPCL},
		"clocal": &bit{word: C, mask: syscall.CLOCAL},
	}
	cc = map[string]int{
		"min":   5,
		"time":  0,
		"lnext": syscall.VLNEXT,
		//"flush": syscall.VFLUSH,
		"intr":  syscall.VINTR,
		"quit":  syscall.VQUIT,
		"erase": syscall.VERASE,
		"kill":  syscall.VKILL,
		"eof":   syscall.VEOF,
		"eol":   syscall.VEOL,
		"eol2":  syscall.VEOL2,
		//"swtch": syscall.VSWTCH,
		"start": syscall.VSTART,
		"stop":  syscall.VSTOP,
		"susp":  syscall.VSUSP,
		//"rprnt": syscall.VRPRNT,
		"werase": syscall.VWERASE,
	}
)

// These consts describe the offsets into the termios struct of various elements.
const (
	I = iota // Input control
	O        // Output control
	C        // Control
	L        // Line control
)

// ioctl constants
const (
	syscallTIOCGWINSZ = 0x5413
	syscallTIOCSWINSZ = 0x5414
)

func tiGet(fd uintptr) (*syscall.Termios, error) {
	var term syscall.Termios
	r1, _, errno := syscall.Syscall(syscall.SYS_IOCTL,
		fd, uintptr(syscall.TCGETS),
		uintptr(unsafe.Pointer(&term)))

	if errno != 0 || r1 != 0 {
		return nil, fmt.Errorf("tiGet: r1 %v, errno %v", r1, errno)
	}
	return &term, nil
}

func tiSet(fd uintptr, term *syscall.Termios) error {
	r1, _, errno := syscall.Syscall(syscall.SYS_IOCTL,
		fd, uintptr(syscall.TCSETS),
		uintptr(unsafe.Pointer(term)))
	if errno != 0 || r1 != 0 {
		return fmt.Errorf("tiSet: r1 %v, errno %v", r1, errno)
	}

	return nil
}

func wsGet(fd uintptr) (*winsize, error) {
	var w winsize
	r1, _, errno := syscall.Syscall(syscall.SYS_IOCTL,
		fd, uintptr(syscallTIOCGWINSZ),
		uintptr(unsafe.Pointer(&w)))
	if errno != 0 || r1 != 0 {
		return nil, fmt.Errorf("wsGet: r1 %v, errno %v", r1, errno)
	}

	return &w, nil
}

func wsSet(fd uintptr, w *winsize) error {
	r1, _, errno := syscall.Syscall(syscall.SYS_IOCTL,
		fd, uintptr(syscallTIOCSWINSZ),
		uintptr(unsafe.Pointer(w)))
	if errno != 0 || r1 != 0 {
		return fmt.Errorf("wsSet: r1 %v, errno %v", r1, errno)
	}

	return nil
}
