//
// Copyright 2014-2023 Cristian Maglie. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

//go:build darwin || dragonfly || freebsd || netbsd || openbsd

package serial

import "golang.org/x/sys/unix"

func (port *unixPort) Drain() error {
	return unix.IoctlSetInt(port.handle, unix.TIOCDRAIN, 0)
}
