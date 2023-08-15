// Copyright 2015-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package termios

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"

	"golang.org/x/sys/windows"
)

type TTYIO struct {
	f *os.File
	m master
}

// Winsize does nothing yet.
type Winsize struct {
	Row int
	Col int
}

// New creates a new TTYIO using windows.STD_INPUT_HANDLE
func New() (*TTYIO, error) {
	h, err := windows.GetStdHandle(windows.STD_INPUT_HANDLE)
	//h, err := syscall.GetStdHandle(syscall.STD_INPUT_HANDLE)
	if err != nil {
		return nil, err
	}

	t := &TTYIO{f: os.NewFile(uintptr(h), "stdin")}
	(&t.m).initStdios()

	return t, nil
}

// NewWithDev creates a new TTYIO with the specified device
func NewWithDev(device string) (*TTYIO, error) {
	return nil, syscall.ENOSYS
}

// NewTTYS returns a new TTYIO, using the passed-in name.
// It will not work on windows at present..
func NewTTYS(port string) (*TTYIO, error) {
	return nil, os.ErrNotExist
}

// GetTermios returns a filled-in Termios, from an fd.
// The uintptr fd is typically derived from windows.GetStdHandle.
func GetTermios(fd uintptr) (*Termios, error) {
	tty, err := New()
	if err != nil {
		return nil, err
	}

	return &Termios{tty: *tty}, nil
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
	(&raw.tty.m).SetRaw()
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

// Copyright The containerd Authors.

// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at

//     http://www.apache.org/licenses/LICENSE-2.0

// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

type master struct {
	in     windows.Handle
	inMode uint32

	out     windows.Handle
	outMode uint32

	err     windows.Handle
	errMode uint32

	vtInputSupported bool
}

func (m *master) initStdios() {
	// Note: We discard console mode warnings, because in/out can be redirected.
	//
	// TODO: Investigate opening CONOUT$/CONIN$ to handle this correctly

	m.in = windows.Handle(os.Stdin.Fd())
	if err := windows.GetConsoleMode(m.in, &m.inMode); err == nil {
		// Validate that windows.ENABLE_VIRTUAL_TERMINAL_INPUT is supported, but do not set it.
		if err = windows.SetConsoleMode(m.in, m.inMode|windows.ENABLE_VIRTUAL_TERMINAL_INPUT); err == nil {
			m.vtInputSupported = true
		}
		// Unconditionally set the console mode back even on failure because SetConsoleMode
		// remembers invalid bits on input handles.
		windows.SetConsoleMode(m.in, m.inMode)
	}

	m.out = windows.Handle(os.Stdout.Fd())
	if err := windows.GetConsoleMode(m.out, &m.outMode); err == nil {
		if err := windows.SetConsoleMode(m.out, m.outMode|windows.ENABLE_VIRTUAL_TERMINAL_PROCESSING); err == nil {
			m.outMode |= windows.ENABLE_VIRTUAL_TERMINAL_PROCESSING
		} else {
			windows.SetConsoleMode(m.out, m.outMode)
		}
	}

	m.err = windows.Handle(os.Stderr.Fd())
	if err := windows.GetConsoleMode(m.err, &m.errMode); err == nil {
		if err := windows.SetConsoleMode(m.err, m.errMode|windows.ENABLE_VIRTUAL_TERMINAL_PROCESSING); err == nil {
			m.errMode |= windows.ENABLE_VIRTUAL_TERMINAL_PROCESSING
		} else {
			windows.SetConsoleMode(m.err, m.errMode)
		}
	}
}

func (m *master) SetRaw() error {
	if err := makeInputRaw(m.in, m.inMode, m.vtInputSupported); err != nil {
		return err
	}

	// Set StdOut and StdErr to raw mode, we ignore failures since
	// windows.DISABLE_NEWLINE_AUTO_RETURN might not be supported on this version of
	// Windows.

	windows.SetConsoleMode(m.out, m.outMode|windows.DISABLE_NEWLINE_AUTO_RETURN)

	windows.SetConsoleMode(m.err, m.errMode|windows.DISABLE_NEWLINE_AUTO_RETURN)

	return nil
}

func (m *master) Reset() error {
	var errs []error

	for _, s := range []struct {
		fd   windows.Handle
		mode uint32
	}{
		{m.in, m.inMode},
		{m.out, m.outMode},
		{m.err, m.errMode},
	} {
		if err := windows.SetConsoleMode(s.fd, s.mode); err != nil {
			// we can't just abort on the first error, otherwise we might leave
			// the console in an unexpected state.
			errs = append(errs, fmt.Errorf("unable to restore console mode: %w", err))
		}
	}

	if len(errs) > 0 {
		return errs[0]
	}

	return nil
}

// func (m *master) Size() (WinSize, error) {
// 	var info windows.ConsoleScreenBufferInfo
// 	err := windows.GetConsoleScreenBufferInfo(m.out, &info)
// 	if err != nil {
// 		return WinSize{}, fmt.Errorf("unable to get console info: %w", err)
// 	}

// 	winsize := WinSize{
// 		Width:  uint16(info.Window.Right - info.Window.Left + 1),
// 		Height: uint16(info.Window.Bottom - info.Window.Top + 1),
// 	}

// 	return winsize, nil
// }

// func (m *master) Resize(ws WinSize) error {
// 	return ErrNotImplemented
// }

func (m *master) DisableEcho() error {
	mode := m.inMode &^ windows.ENABLE_ECHO_INPUT
	mode |= windows.ENABLE_PROCESSED_INPUT
	mode |= windows.ENABLE_LINE_INPUT

	if err := windows.SetConsoleMode(m.in, mode); err != nil {
		return fmt.Errorf("unable to set console to disable echo: %w", err)
	}

	return nil
}

func (m *master) Close() error {
	return nil
}

func (m *master) Read(b []byte) (int, error) {
	return os.Stdin.Read(b)
}

func (m *master) Write(b []byte) (int, error) {
	return os.Stdout.Write(b)
}

func (m *master) Fd() uintptr {
	return uintptr(m.in)
}

// on windows, console can only be made from os.Std{in,out,err}, hence there
// isnt a single name here we can use. Return a dummy "console" value in this
// case should be sufficient.
func (m *master) Name() string {
	return "console"
}

// makeInputRaw puts the terminal (Windows Console) connected to the given
// file descriptor into raw mode
func makeInputRaw(fd windows.Handle, mode uint32, vtInputSupported bool) error {
	// See
	// -- https://msdn.microsoft.com/en-us/library/windows/desktop/ms686033(v=vs.85).aspx
	// -- https://msdn.microsoft.com/en-us/library/windows/desktop/ms683462(v=vs.85).aspx

	// Disable these modes
	mode &^= windows.ENABLE_ECHO_INPUT
	mode &^= windows.ENABLE_LINE_INPUT
	mode &^= windows.ENABLE_MOUSE_INPUT
	mode &^= windows.ENABLE_WINDOW_INPUT
	mode &^= windows.ENABLE_PROCESSED_INPUT

	// Enable these modes
	mode |= windows.ENABLE_EXTENDED_FLAGS
	mode |= windows.ENABLE_INSERT_MODE
	mode |= windows.ENABLE_QUICK_EDIT_MODE

	if vtInputSupported {
		mode |= windows.ENABLE_VIRTUAL_TERMINAL_INPUT
	}

	if err := windows.SetConsoleMode(fd, mode); err != nil {
		return fmt.Errorf("unable to set console to raw mode: %w", err)
	}

	return nil
}

// func checkConsole(f File) error {
// 	var mode uint32
// 	if err := windows.GetConsoleMode(windows.Handle(f.Fd()), &mode); err != nil {
// 		return err
// 	}
// 	return nil
// }
