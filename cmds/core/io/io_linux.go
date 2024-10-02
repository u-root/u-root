// Copyright 2010-2020 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//go:build !tinygo || tinygo.enable

package main

import (
	"github.com/u-root/u-root/pkg/memio"
)

func init() {
	addCmd(readCmds, "rb", &cmd{memio.Read, 64, 8})
	addCmd(readCmds, "rw", &cmd{memio.Read, 64, 16})
	addCmd(readCmds, "rl", &cmd{memio.Read, 64, 32})
	addCmd(readCmds, "rq", &cmd{memio.Read, 64, 64})

	addCmd(writeCmds, "wb", &cmd{memio.Write, 64, 8})
	addCmd(writeCmds, "ww", &cmd{memio.Write, 64, 16})
	addCmd(writeCmds, "wl", &cmd{memio.Write, 64, 32})
	addCmd(writeCmds, "wq", &cmd{memio.Write, 64, 64})

	usageMsg += `io (r{b,w,l,q} address)...
io (w{b,w,l,q} address value)...
`
}
