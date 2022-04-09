// Copyright 2021 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build !plan9 && !windows && !darwin
// +build !plan9,!windows,!darwin

package termios

import "golang.org/x/sys/unix"

const (
	gets       = unix.TCGETS
	sets       = unix.TCSETS
	getWinSize = unix.TIOCGWINSZ
	setWinSize = unix.TIOCSWINSZ
)

func speed(speed int) uint32 {
	return uint32(speed)
}
