//
// Copyright 2014-2023 Cristian Maglie. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

//go:build linux || freebsd || openbsd

package serial

import "golang.org/x/sys/unix"

func (port *unixPort) ResetInputBuffer() error {
	return unix.IoctlSetInt(port.handle, ioctlTcflsh, unix.TCIFLUSH)
}

func (port *unixPort) ResetOutputBuffer() error {
	return unix.IoctlSetInt(port.handle, ioctlTcflsh, unix.TCOFLUSH)
}
