// Copyright 2015-2020 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package pty provides basic pty support.
// It implments much of exec.Command
// but the Start() function starts two goroutines that relay the
// data for Stdin, Stdout, and Stdout such that proper kernel pty
// processing is done. We did not simply embed an exec.Command
// as we can no guarantee that we can implement all aspects of it
// for all time to come.
package pty

import (
	"fmt"
	"os"
	"syscall"
	"unsafe"

	"github.com/u-root/u-root/pkg/termios"
)

// New returns a new Pty.
func New() (*Pty, error) {
	tty, err := termios.New()
	if err != nil {
		return nil, err
	}
	restorer, err := tty.Get()
	if err != nil {
		return nil, err
	}

	ptm, err := os.OpenFile("/dev/ptmx", os.O_RDWR, 0)
	if err != nil {
		return nil, err
	}

	if err := ptsunlock(ptm); err != nil {
		return nil, err
	}

	sname, err := ptsname(ptm)
	if err != nil {
		return nil, err
	}

	// It can take a non-zero time for a pts to appear, it seems.
	// Ten tries is reported to be far more than enough.
	// We could consider something like inotify rather than polling?
	for i := 0; i < 10; i++ {
		_, err := os.Stat(sname)
		if err == nil {
			break
		}
	}

	pts, err := os.OpenFile(sname, os.O_RDWR|syscall.O_NOCTTY, 0)
	if err != nil {
		return nil, err
	}
	return &Pty{Ptm: ptm, Pts: pts, Sname: sname, Kid: -1, TTY: tty, Restorer: restorer}, nil
}

func ptsname(f *os.File) (string, error) {
	var n uintptr
	_, _, err := syscall.Syscall(syscall.SYS_IOCTL, f.Fd(), syscall.TIOCGPTN, uintptr(unsafe.Pointer(&n)))
	if err != 0 {
		return "", err
	}
	return fmt.Sprintf("/dev/pts/%d", n), nil
}

func ptsunlock(f *os.File) error {
	var u uintptr
	// freebsd people, what is the way to do this?
	// use TIOCSPTLCK with a zero valued arg to clear the slave pty lock
	_, _, err := syscall.Syscall(syscall.SYS_IOCTL, f.Fd() /*syscall.TIOCSPTLCK*/, 0, uintptr(unsafe.Pointer(&u)))
	if err != 0 {
		return err
	}
	return nil
}

func sysLinux(p *Pty) {
	p.C.SysProcAttr = &syscall.SysProcAttr{Setctty: true, Setsid: true}
}

func init() {
	sys = sysLinux
}
