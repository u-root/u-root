// Copyright 2012-2020 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build amd64 386

package main

import (
	"github.com/u-root/u-root/pkg/cmos"
)

func init() {
	usageMsg += `io (cr index)... # read from CMOS register index
io (cw index value)... # write value to CMOS register index
`
	readCmds["cr"] = cmd{cmos.Read, 7, 8}
	writeCmds["cw"] = cmd{cmos.Write, 7, 8}
}
