// Copyright 2015-2020 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package pty provides basic pty support.
// It implements much of exec.Command
// but the Start() function starts two goroutines that relay the
// data for Stdin, Stdout, and Stdout such that proper kernel pty
// processing is done. We did not simply embed an exec.Command
// as we can no guarantee that we can implement all aspects of it
// for all time to come.
package pty

import (
	"os"
	"strconv"
	"syscall"
	"time"

	"github.com/u-root/u-root/pkg/termios"
	"golang.org/x/sys/unix"
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
	for range 10 {
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
	n, err := unix.IoctlGetInt(int(f.Fd()), unix.TIOCGPTN)
	if err != nil {
		return "", err
	}
	return "/dev/pts/" + strconv.Itoa(n), nil
}

func ptsunlock(f *os.File) error {
	// use TIOCSPTLCK with a zero valued arg to clear the slave pty lock
	return unix.IoctlSetPointerInt(int(f.Fd()), unix.TIOCSPTLCK, 0)
}

func sysLinux(p *Pty) {
	p.C.SysProcAttr = &syscall.SysProcAttr{Setctty: true, Setsid: true}
}

func init() {
	sys = sysLinux
}

// Copied from another pty pkg
// TODO: Merge into New function
func NewPTMS() (*os.File, *os.File, error) {
	p, err := os.OpenFile("/dev/ptmx", os.O_RDWR, 0)
	if err != nil {
		time.Sleep(1 * time.Second)
		return nil, nil, err
	}

	// unlock
	// use TIOCSPTLCK with a zero valued arg to clear the lock.
	if err := unix.IoctlSetPointerInt(int(p.Fd()), unix.TIOCSPTLCK, 0); err != nil {
		return nil, nil, err
	}

	sname, err := ptsname(p)
	if err != nil {
		return nil, nil, err
	}

	t, err := os.OpenFile(sname, os.O_RDWR|syscall.O_NOCTTY|syscall.O_NONBLOCK, 0o620) //nolint:gosec // Expected Open from a variable.
	if err != nil {
		return nil, nil, err
	}

	return p, t, nil
}
