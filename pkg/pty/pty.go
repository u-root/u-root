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
	"io"
	"os"
	"os/exec"

	"github.com/u-root/u-root/pkg/termios"
)

// Pty contains all the bits and pieces needed to start and control
// a process at the other end of a pty, as well as whatever is needed
// to control the pty, and restore modes on exit.
type Pty struct {
	C        *exec.Cmd
	Ptm      *os.File
	Pts      *os.File
	Sname    string
	Kid      int
	TTY      *termios.TTYIO
	WS       *termios.Winsize
	Restorer *termios.Termios
}

var sys = func(p *Pty) {}

// Command sets up an exec.Command at the remote end of a Pty.
func (p *Pty) Command(cmd string, args ...string) {
	p.C = exec.Command(cmd, args...)
	p.C.Stdin, p.C.Stdout, p.C.Stderr = p.Pts, p.Pts, p.Pts
	sys(p)
}

// Start starts Command attached to a Pty. It sets
// window size and other variables as needed. It does not
// block.
func (p *Pty) Start() error {
	tty, err := termios.New()
	if err != nil {
		return err
	}

	if p.WS, err = tty.GetWinSize(); err != nil {
		return err
	}

	if p.Restorer, err = tty.Raw(); err != nil {
		return err
	}

	if err := p.C.Start(); err != nil {
		tty.Set(p.Restorer)
		return err
	}
	p.Kid = p.C.Process.Pid

	// We make a good faith effort to set the
	// WinSize of the Pts, but it's not a deal breaker
	// if we can't do it.
	if err := termios.SetWinSize(p.Pts.Fd(), p.WS); err != nil {
		fmt.Fprintf(p.C.Stderr, "SetWinSize of Pts: %v", err)
	}

	return nil
}

// Run runs a Command attached to a Pty, waiting for completion and
// managing stdio. It uses Wait to restore tty modes when it is done.
func (p *Pty) Run() error {
	if err := p.Start(); err != nil {
		return err
	}

	go io.Copy(p.TTY, p.Ptm)

	// The 1 byte for IO may seem weird, but ptys are for human interaction
	// and, let's face it, we don't all type fast.
	go func() {
		var data [1]byte
		for {
			if _, err := p.TTY.Read(data[:]); err != nil {
				return
			}
			// Log the error but it may be transient.
			if _, err := p.Ptm.Write(data[:]); err != nil {
				fmt.Fprintf(p.C.Stderr, "Error writing input to ptm: %v: give up\n", err)
			}
		}
	}()
	return p.Wait()
}

// Wait waits for a previously started command to finish, and restores the
// tty mode when it is done.
func (p *Pty) Wait() error {
	defer p.TTY.Set(p.Restorer)
	return p.C.Wait()
}
