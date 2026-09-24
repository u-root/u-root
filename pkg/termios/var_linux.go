// Copyright 2015-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package termios

import (
	"maps"
	"syscall"
)

// with BOTHER, Linux, sensibly, does not even recognize bauds below 57600.
// This begs the question of just how we get that information out:
// turns out there is another new system call.
// For bauds below 57600, we punt, for now, because 115200 is the new universal
// constant.
// These appear not to be define in Go at present.

const (
	B0 uint32 = iota
	B50
	B75
	B110
	B134
	B150
	B200
	B300
	B600
	B1200
	B1800
	B2400
	B4800
	B9600
	B19200
	B38400
)

const (
	BOTHER uint32 = iota + 0x00001000
	B57600
	B115200
	B230400
	B460800
	B500000
	B576000
	B921600
	B1000000
	B1152000
	B1500000
	B2000000
	B2500000
	B3000000
	B3500000
	B4000000
)

const (
	CBAUD  = 0x0000100f
	CIBAUD = 0x100f0000 // input baud rate
)

func init() {
	if B4000000 == BOTHER {
		panic("fuck")
	}
}

// baud2unixB converts an integer baud rate to a Linux-style B* variable.
var baud2unixB = map[int]uint32{
	0:     B0,
	50:    B50,
	75:    B75,
	110:   B110,
	134:   B134,
	150:   B150,
	200:   B200,
	300:   B300,
	600:   B600,
	1200:  B1200,
	1800:  B1800,
	2400:  B2400,
	4800:  B4800,
	9600:  B9600,
	19200: B19200,
	38400: B38400,

	57600:   B57600,
	115200:  B115200,
	230400:  B230400,
	460800:  B460800,
	500000:  B500000,
	576000:  B576000,
	921600:  B921600,
	1000000: B1000000,
	1152000: B1152000,
	1500000: B1500000,
	2000000: B2000000,
	2500000: B2500000,
	3000000: B3000000,
	3500000: B3500000,
	4000000: B4000000,
}

var bother2baud = map[uint32]int{
	B0:     0,
	B50:    50,
	B75:    75,
	B110:   110,
	B134:   134,
	B150:   150,
	B200:   200,
	B300:   300,
	B600:   600,
	B1200:  1200,
	B1800:  1800,
	B2400:  2400,
	B4800:  4800,
	B9600:  9600,
	B19200: 19200,
	B38400: 38400,

	B57600:   57600,
	B115200:  115200,
	B230400:  230400,
	B460800:  460800,
	B500000:  500000,
	B576000:  576000,
	B921600:  921600,
	B1000000: 1000000,
	B1152000: 1152000,
	B1500000: 1500000,
	B2000000: 2000000,
	B2500000: 2500000,
	B3000000: 3000000,
	B3500000: 3500000,
	B4000000: 4000000,
}

// init adds constants that are linux-specific
func init() {
	extra := map[string]*bit{
		"iuclc": {word: I, mask: syscall.IUCLC},
		"olcuc": {word: O, mask: syscall.OLCUC},
		"xcase": {word: L, mask: syscall.XCASE},
		// not in FreeBSD
		"iutf8": {word: I, mask: syscall.IUTF8},
		"ofill": {word: O, mask: syscall.OFILL},
		"ofdel": {word: O, mask: syscall.OFDEL},
	}
	maps.Copy(boolFields, extra)
}
