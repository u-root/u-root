// Copyright 2022 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"os"
	"os/exec"
	"syscall"
	"testing"
)

func TestAttr(t *testing.T) {
	var err error
	ttyf, err = os.OpenFile("/dev/tty", os.O_RDWR, 0)
	if err != nil {
		t.Skipf("can not open /dev/tty, skipping test")
		return
	}

	c := &Command{Cmd: exec.Command("true")}

	forkAttr(c)

	if !c.Cmd.SysProcAttr.Foreground {
		t.Errorf("forkAttr(&c): c.Cmd.SysProcAttr.Foreground is false")
	}

	fd := int(ttyf.Fd())
	if c.Cmd.SysProcAttr.Ctty != fd {
		t.Errorf("forkAttr(&c): c.Cmd.SysProcAttr.Ctty %d != %d", c.Cmd.SysProcAttr.Ctty, fd)
	}

	c.BG = true
	forkAttr(c)
	if !c.Cmd.SysProcAttr.Setpgid {
		t.Errorf("forkAttr(&c): c.Cmd.SysProcAttr.Setpgid is not true although BG (BackGround) was set")
	}

	builtinAttr(c)
	// make sure the right struct member is set.
	if c.Cmd.SysProcAttr.Cloneflags&syscall.CLONE_NEWNS != syscall.CLONE_NEWNS {
		t.Errorf("builtinAttr(&c): c.Cmd.SysProcAttr.Cloneflags did not have syscall.CLONE_NEWNS set")
	}
}
