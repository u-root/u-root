// Copyright 2012-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build !tinygo && linux

package main

import (
	"syscall"
)

func builtinAttr(c *Command) {
	c.Cmd.SysProcAttr.Cloneflags |= syscall.CLONE_NEWNS
}

func forkAttr(c *Command) {
	c.Cmd.SysProcAttr = &syscall.SysProcAttr{}
	if c.BG {
		c.Cmd.SysProcAttr.Setpgid = true
	} else {
		c.Cmd.SysProcAttr.Foreground = true
		c.Cmd.SysProcAttr.Ctty = int(ttyf.Fd())
	}
}
