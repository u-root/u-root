// Copyright 2015-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package termios implements basic termios operations including getting
// a tty struct, termio struct, a winsize struct, and setting raw mode.
// To get a TTY, call termios.New.
// To get a Termios, call tty.Get(); to set it, call tty.Set(*Termios)
// To set raw mode and then restore, one can do:
// tty := termios.NewTTY()
// restorer, err := tty.Raw()
// do things
// tty.Set(restorer)
package termios

import (
	"golang.org/x/sys/unix"
)

func (t *TTY) Raw() (*unix.Termios, error) {
	restorer, err := t.Get()
	if err != nil {
		return nil, err
	}

	raw := MakeRaw(restorer)

	if err := t.Set(raw); err != nil {
		return nil, err
	}
	return restorer, nil
}
