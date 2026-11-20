// Copyright 2026 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pty

import (
	"github.com/u-root/u-root/pkg/termios"
)

// callerTerminal returns the settings of the terminal this process is attached
// to, or nils when it has none.
//
// Plan 9 has no controlling-terminal concept to fail on, and pty.New there is
// still unimplemented, so this only has to satisfy Start.
func callerTerminal(dev string) (*termios.TTYIO, *termios.Termios, error) {
	tty, err := termios.NewWithDev(dev)
	if err != nil {
		return nil, nil, nil
	}
	restorer, err := tty.Get()
	if err != nil {
		return nil, nil, err
	}
	return tty, restorer, nil
}
