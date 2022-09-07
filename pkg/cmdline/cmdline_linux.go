// Copyright 2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cmdline

import (
	"os"
)

const cmdLinePath = "/proc/cmdline"

var procCmdLine *CmdLine

func cmdLine(n string) *CmdLine {
	procCmdLine = &CmdLine{AsMap: map[string]string{}}
	r, err := os.Open(n)
	if err != nil {
		procCmdLine.Err = err
		return procCmdLine
	}

	defer r.Close()

	procCmdLine = parse(r)
	return procCmdLine
}

func getCmdLine() *CmdLine {
	return cmdLine(cmdLinePath)
}
