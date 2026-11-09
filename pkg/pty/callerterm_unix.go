// Copyright 2026 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build !plan9 && !windows

package pty

import (
	"errors"
	"syscall"

	"github.com/u-root/u-root/pkg/termios"
)

// callerTerminal returns the settings of the terminal this process is attached
// to, or nils when it has none.
//
// The distinction matters because "no controlling terminal" is not a failure to
// allocate a pty: it means there is nothing to copy, so the pty keeps its
// defaults and there is nothing to restore when the command exits. It is a
// failure of the terminal, not of the pty, that is an error.
func callerTerminal(dev string) (*termios.TTYIO, *termios.Termios, error) {
	tty, err := termios.NewWithDev(dev)
	if err != nil {
		if errors.Is(err, syscall.ENXIO) || errors.Is(err, syscall.ENOENT) || errors.Is(err, syscall.ENODEV) {
			return nil, nil, nil
		}
		return nil, nil, err
	}
	restorer, err := tty.Get()
	if err != nil {
		return nil, nil, err
	}
	return tty, restorer, nil
}
