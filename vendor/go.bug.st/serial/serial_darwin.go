//
// Copyright 2014-2023 Cristian Maglie. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package serial

import "golang.org/x/sys/unix"

const devFolder = "/dev"
const regexFilter = "^(cu|tty)\\..*"

const ioctlTcgetattr = unix.TIOCGETA
const ioctlTcsetattr = unix.TIOCSETA
const ioctlTcflsh = unix.TIOCFLUSH
const ioctlTioccbrk = unix.TIOCCBRK
const ioctlTiocsbrk = unix.TIOCSBRK

func setTermSettingsBaudrate(speed int, settings *unix.Termios) (error, bool) {
	baudrate, ok := baudrateMap[speed]
	if !ok {
		return nil, true
	}
	settings.Ispeed = toTermiosSpeedType(baudrate)
	settings.Ospeed = toTermiosSpeedType(baudrate)
	return nil, false
}

func (port *unixPort) setSpecialBaudrate(speed uint32) error {
	const kIOSSIOSPEED = 0x80045402
	return unix.IoctlSetPointerInt(port.handle, kIOSSIOSPEED, int(speed))
}

func (port *unixPort) ResetInputBuffer() error {
	return unix.IoctlSetPointerInt(port.handle, ioctlTcflsh, unix.TCIFLUSH)
}

func (port *unixPort) ResetOutputBuffer() error {
	return unix.IoctlSetPointerInt(port.handle, ioctlTcflsh, unix.TCOFLUSH)
}
