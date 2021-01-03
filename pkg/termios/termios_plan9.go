// Copyright 2015-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build plan9

package termios

import (
	"fmt"
	"os"
)

// Termios is used to manipulate the control channel of a kernel.
type Termios struct {
	// current state
	state string
}

// Winsize holds the window size information, it is modeled on unix.Winsize.
type Winsize struct {
	Row    uint16
	Col    uint16
	Xpixel uint16
	Ypixel uint16
}

// TTYIO is a wrapper that only allows Read and Write.
// Plan 9 inherited little of the 1900-teletype-style
// API from Unix. But so many programs we want to support
// seem to need this nonsense, we'll pretend to support
// it for now.
type TTYIO struct {
	f *os.File
	// The control channel is used for Raw() control only.
	c *os.File
	Termios
}

// New creates a new TTYIO using /dev/tty
func New() (*TTYIO, error) {
	return NewWithDev("/dev/cons")
}

// NewWithDev creates a new TTYIO with the specified device
func NewWithDev(device string) (*TTYIO, error) {
	f, err := os.OpenFile(device, os.O_RDWR, 0)
	if err != nil {
		return nil, err
	}
	c, err := os.OpenFile(device+"ctl", os.O_WRONLY, 0)
	if err != nil {
		return nil, err
	}
	return &TTYIO{f: f, c: c, Termios: Termios{state: ""}}, nil
}

// NewTTYS returns a new TTYIO.
func NewTTYS(port string) (*TTYIO, error) {
	return nil, fmt.Errorf("not yet")
}

// Get terms a Termios from a TTYIO.
func (t *TTYIO) Get() (*Termios, error) {
	nt := t.Termios
	return &nt, nil
}

// Set sets tty parameters for a TTYIO from a Termios.
func (t *TTYIO) Set(ti *Termios) error {
	_, err := t.c.Write([]byte(ti.state))
	return err
}

// GetWinSize gets window size from an fd.
func GetWinSize(fd uintptr) (*Winsize, error) {
	return nil, fmt.Errorf("not yet")
}

// GetWinSize gets window size from a TTYIO.
// TODO: pick it up from rio
func (t *TTYIO) GetWinSize() (*Winsize, error) {
	return &Winsize{24, 80, 1024, 768}, nil
}

// SetWinSize sets window size for an fd from a Winsize.
func SetWinSize(fd uintptr, w *Winsize) error {
	return fmt.Errorf("Not yet")
}

// SetWinSize sets window size for a TTYIO from a Winsize.
func (t *TTYIO) SetWinSize(w *Winsize) error {
	return fmt.Errorf("Plan 9: not yet")
}

// MakeRaw modifies Termio state so, when Set is called, it will set it to raw mode.
func MakeRaw(term *Termios) *Termios {
	// This is dumb but should always work.
	raw := *term
	raw.state = "rawon"
	return &raw
}

// MakeSerialBaud updates the Termios to set the baudrate
// not on plan 9 however.
// we'll need to rework a lot of this nonsense.
func MakeSerialBaud(term *Termios, baud int) (*Termios, error) {
	t := *term

	return &t, nil
}

// MakeSerialDefault updates the Termios to typical serial configuration:
// not know how to do this yet.
func MakeSerialDefault(term *Termios) *Termios {
	t := *term

	return &t
}
