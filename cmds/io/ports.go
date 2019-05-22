// Copyright 2012-2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build amd64 386

package main

import (
	"github.com/u-root/u-root/pkg/memio"
)

func init() {
	usageMsg += `io (in{b,w,l} address)...
io (out{b,w,l} address value)...
`
	readCmds["inb"] = cmd{in, 16, 8}
	readCmds["inw"] = cmd{in, 16, 16}
	readCmds["inl"] = cmd{in, 16, 32}
	writeCmds["outb"] = cmd{out, 16, 8}
	writeCmds["outw"] = cmd{out, 16, 16}
	writeCmds["outl"] = cmd{out, 16, 32}
}

func in(addr int64, data memio.UintN) error {
	return memio.In(uint16(addr), data)
}

func out(addr int64, data memio.UintN) error {
	return memio.Out(uint16(addr), data)
}
