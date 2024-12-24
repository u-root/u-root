//
// Copyright 2014-2023 Cristian Maglie. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package serial

import "golang.org/x/sys/unix"

// termios manipulation functions

var baudrateMap = map[int]uint32{
	0:      unix.B9600, // Default to 9600
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

var databitsMap = map[int]uint32{
	0: unix.CS8, // Default to 8 bits
	5: unix.CS5,
	6: unix.CS6,
	7: unix.CS7,
	8: unix.CS8,
}

const tcCMSPAR uint32 = 0 // may be CMSPAR or PAREXT
const tcIUCLC uint32 = 0

const tcCCTS_OFLOW uint32 = 0x00010000
const tcCRTS_IFLOW uint32 = 0x00020000

const tcCRTSCTS uint32 = (tcCCTS_OFLOW | tcCRTS_IFLOW)

func toTermiosSpeedType(speed uint32) uint32 {
	return speed
}
