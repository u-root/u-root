// Copyright 2015-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package termios

import (
	"fmt"
	"os"
	"path/filepath"
)

const reset = "rawoff\nholdoff\n"

// Termios is used to manipulate the control channel of a kernel.
type Termios struct {
	mode string
	// ctl is used to set raw mode.
	// It may not exist, but IO to the tty
	// may still work. Hence, it is not opened
	// until needed. Once opened, it is left open,
	// as the modes reset once it is closed.
	*os.File
}

// Winsize holds the window size information, it is modeled on unix.Winsize.
type Winsize struct {
	Row    uint16
	Col    uint16
	Xpixel uint16
	Ypixel uint16
}

// TTYIO is a wrapper that only allows Read and Write.
type TTYIO struct {
	f *os.File
}

// New creates a new TTYIO using /dev/cons
func New() (*TTYIO, error) {
	return NewWithDev("/dev/cons")
}

// NewWithDev creates a new TTYIO with the specified device
func NewWithDev(device string) (*TTYIO, error) {
	f, err := os.OpenFile(device, os.O_RDWR, 0)
	if err != nil {
		return nil, err
	}
	return &TTYIO{f: f}, nil
}

// NewTTYS returns a new TTYIO.
func NewTTYS(port string) (*TTYIO, error) {
	f, err := os.OpenFile(filepath.Join("/dev", port), os.O_RDWR, 0)
	if err != nil {
		return nil, err
	}
	return &TTYIO{f: f}, nil
}

// GetTermios returns a filled-in Termios, from an fd.
// And, sorry, on Plan 9, there seems to be no way to
// find out if it is in hold/raw mode. Odd.
// Because the ctl file is separate, do not open
// it until it is needed.
func GetTermios(fd uintptr) (*Termios, error) {
	f, err := consctlFile("/", fd)
	if err != nil {
		return nil, err
	}
	return &Termios{mode: reset, File: f}, nil
}

// Get a Termios from a TTYIO.
func (t *TTYIO) Get() (*Termios, error) {
	return GetTermios(t.f.Fd())
}

// SetTermios sets tty parameters for an fd from a Termios.
// The only thing we can do is write the mode to the ctl.
func SetTermios(_ uintptr, t *Termios) error {
	if t.File == nil {
		return fmt.Errorf("termios ctl is not set up:%w", os.ErrInvalid)
	}
	if _, err := t.Write([]byte(t.mode)); err != nil {
		return fmt.Errorf("writing %q to %v: %w", t.mode, t.Name(), err)
	}
	return nil
}

// Set sets tty parameters for a TTYIO from a Termios.
func (*TTYIO) Set(ti *Termios) error {
	return SetTermios(0, ti)
}

// GetWinSize gets window size from an fd.
func GetWinSize(_ uintptr) (*Winsize, error) {
	r, c, err := readWinSize("/dev/wctl")
	if err != nil {
		return nil, err
	}
	return &Winsize{Row: r, Col: c}, nil
}

// GetWinSize gets window size from a TTYIO.
func (t *TTYIO) GetWinSize() (*Winsize, error) {
	return GetWinSize(t.f.Fd())
}

// SetWinSize sets window size for an fd from a Winsize.
func SetWinSize(_ uintptr, _ *Winsize) error {
	return fmt.Errorf("plan 9: not yet")
}

// SetWinSize sets window size for a TTYIO from a Winsize.
func (t *TTYIO) SetWinSize(w *Winsize) error {
	return SetWinSize(t.f.Fd(), w)
}

// MakeRaw modifies Termio state so, if it used for an fd or tty, it will set it to raw mode.
func MakeRaw(term *Termios) *Termios {
	raw := *term
	raw.mode = "rawon"
	return &raw
}

// MakeSerialBaud updates the Termios to set the baudrate
func MakeSerialBaud(term *Termios, baud int) (*Termios, error) {
	t := *term
	return &t, nil
}

// MakeSerialDefault updates the Termios to typical serial configuration:
//   - Ignore all flow control (modem, hardware, software...)
//   - Translate carriage return to newline on input
//   - Enable canonical mode: Input is available line by line, with line editing
//     enabled (ERASE, KILL are supported)
//   - Local ECHO is added (and handled by line editing)
//   - Map newline to carriage return newline on output
func MakeSerialDefault(term *Termios) *Termios {
	t := *term

	return &t
}
