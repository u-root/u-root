// Copyright 2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build !linux

package cmdline

import "os"

var procCmdLine *CmdLine

func cmdLine(f string) *CmdLine {
	procCmdLine = &CmdLine{AsMap: map[string]string{}, Err: os.ErrNotExist}
	return procCmdLine
}

func getCmdLine() *CmdLine {
	return cmdLine("")
}
