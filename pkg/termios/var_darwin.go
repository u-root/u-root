// Copyright 2021 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package termios

import (
	"maps"
	"syscall"

	"golang.org/x/sys/unix"
)

// baud2unixB convert a baudrate to the corresponding unix const.
var baud2unixB = map[int]uint64{
	50:     unix.B50,
	75:     unix.B75,
	110:    unix.B110,
	134:    unix.B134,
	150:    unix.B150,
	200:    unix.B200,
	300:    unix.B300,
	600:    unix.B600,
	1200:   unix.B1200,
	1800:   unix.B1800,
	2400:   unix.B2400,
	4800:   unix.B4800,
	9600:   unix.B9600,
	19200:  unix.B19200,
	38400:  unix.B38400,
	57600:  unix.B57600,
	115200: unix.B115200,
	230400: unix.B230400,
}

// init adds constants that are darwin-specific
func init() {
	extra := map[string]*bit{
		// not in FreeBSD
		"iutf8": {word: I, mask: syscall.IUTF8},
		"ofill": {word: O, mask: syscall.OFILL},
		"ofdel": {word: O, mask: syscall.OFDEL},
	}
	maps.Copy(boolFields, extra)
}

func toTermiosCflag(r uint64) uint64 { return r }
