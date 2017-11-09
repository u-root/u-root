// Copyright 2015-2017 the u-root Authors. All rights reserved
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
	"golang.org/x/sys/unix"
)

type Pty struct {
	C        *exec.Cmd
	Ptm      *os.File
	Pts      *os.File
	Sname    string
	Kid      int
	TTY      *termios.TTY
	WS       *unix.Winsize
	Restorer *unix.Termios
}

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

	go io.Copy(p.C.Stdout, p.Ptm)

	// The 1 byte for IO may seem weird, but ptys are for human interacxtion
	// and, let's face it, we don't all type fast.
	go func() {
		var data [1]byte
		for {
			if _, err := p.C.Stdin.Read(data[:]); err != nil {
				return
			}
			// TODO: should we really echo this? Or pass it to the
			// shell? I think we need to echo it but not pass it
			// on.
			if data[0] == '\r' {
				if _, err := p.C.Stdout.Write(data[:]); err != nil {
					fmt.Fprintf(p.C.Stderr, "error on echo %v: %v", data, err)
				}
				data[0] = '\n'
			}
			// Log the error but it may be transient.
			if _, err := p.Ptm.Write(data[:]); err != nil {
				fmt.Fprintf(p.C.Stderr, "Error writing input to ptm: %v: give up\n", err)
			}
		}
	}()
	return nil
}

func (p *Pty) Run() error {
	if err := p.Start(); err != nil {
		return err
	}

	return p.Wait()
}

func (p *Pty) Wait() error {
	defer p.TTY.Set(p.Restorer)
	return p.C.Wait()
}
