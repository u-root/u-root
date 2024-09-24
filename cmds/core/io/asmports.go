// Copyright 2012-2020 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build (!tinygo || tinygo.enable) && ((linux && amd64) || (linux && 386))

package main

import (
	"github.com/u-root/u-root/pkg/memio"
)

// The xin* and xout* commands use iopl, and hence
// must be run by root.
func init() {
	usageMsg += `io (xin{b,w,l} address)...
io (xout{b,w,l} address value)...
`
	addCmd(readCmds, "xinb", &cmd{xin, 16, 8})
	addCmd(readCmds, "xinw", &cmd{xin, 16, 16})
	addCmd(readCmds, "xinl", &cmd{xin, 16, 32})
	addCmd(writeCmds, "xoutb", &cmd{xout, 16, 8})
	addCmd(writeCmds, "xoutw", &cmd{xout, 16, 16})
	addCmd(writeCmds, "xoutl", &cmd{xout, 16, 32})
}

func xin(addr int64, data memio.UintN) error {
	return memio.ArchIn(uint16(addr), data)
}

func xout(addr int64, data memio.UintN) error {
	return memio.ArchOut(uint16(addr), data)
}
