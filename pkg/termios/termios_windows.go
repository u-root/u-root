// Copyright 2015-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package termios

import (
	"os"
	"os/exec"
	"path/filepath"
	"syscall"

	"golang.org/x/sys/windows"
)

type TTYIO struct {
	f *os.File
}

// Winsize does nothing yet.
type Winsize struct {
}

// New creates a new TTYIO using windows.STD_INPUT_HANDLE
func New() (*TTYIO, error) {
	h, err := windows.GetStdHandle(windows.STD_INPUT_HANDLE)
	if err != nil {
		return nil, err
	}

	return &TTYIO{f: os.NewFile(uintptr(h), "stdin")}, nil
}

// NewWithDev creates a new TTYIO with the specified device
func NewWithDev(device string) (*TTYIO, error) {
	return nil, syscall.ENOSYS
}

// NewTTYS returns a new TTYIO, using the passed-in name.
// It will not work on windows at present..
func NewTTYS(port string) (*TTYIO, error) {
	f, err := os.OpenFile(filepath.Join("/dev", port), os.O_RDWR, 0o620)
	if err != nil {
		return nil, err
	}
	return &TTYIO{f: f}, nil
}

// GetTermios returns a filled-in Termios, from an fd.
// The uintptr fd is typically derived from windows.GetStdHandle.
func GetTermios(fd uintptr) (*Termios, error) {
	tty, err := New()
	if err != nil {
		return nil, err
	}

	return &Termios{tty: tty}, nil
}

// Get terms a Termios from a TTYIO.
func (t *TTYIO) Get() (*Termios, error) {
	return GetTermios(t.f.Fd())
}

// SetTermios sets tty parameters for an fd from a Termios.
func SetTermios(fd uintptr, ti *Termios) error {
	return syscall.ENOSYS
}

// Set sets tty parameters for a TTYIO from a Termios.
func (t *TTYIO) Set(ti *Termios) error {
	return SetTermios(t.f.Fd(), ti)
}

// GetWinSize gets window size from an fd.
func GetWinSize(fd uintptr) (*Winsize, error) {
	return nil, syscall.ENOSYS
}

// GetWinSize gets window size from a TTYIO.
func (t *TTYIO) GetWinSize() (*Winsize, error) {
	return GetWinSize(t.f.Fd())
}

// SetWinSize sets window size for an fd from a Winsize.
func SetWinSize(fd uintptr, w *Winsize) error {
	return syscall.ENOSYS
}

// SetWinSize sets window size for a TTYIO from a Winsize.
func (t *TTYIO) SetWinSize(w *Winsize) error {
	return SetWinSize(t.f.Fd(), w)
}

// Ctty sets the control tty into a Cmd, from a TTYIO.
// This has no meaning on windows, just ignore it.
// You are always the Ctty.
func (t *TTYIO) Ctty(c *exec.Cmd) {
}

// MakeRaw modifies Termio state so, if it used for an fd or tty, it will set it to raw mode.
func MakeRaw(term *Termios) *Termios {
	raw := *term
	//syscall.MakeRaw(term.h)
	return &raw
}

// MakeSerialBaud does nothing and returns an error
func MakeSerialBaud(term *Termios, baud int) (*Termios, error) {
	return nil, syscall.ENOSYS
}

// MakeSerialDefault does nothing.
func MakeSerialDefault(term *Termios) *Termios {
	t := *term
	return &t
}
