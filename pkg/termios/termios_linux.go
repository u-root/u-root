// Copyright 2015-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package termios

import (
	"fmt"
	"os"
	"syscall"
	"unsafe"
)

// NCCS is the size of the control character array, which is basically
// ^@ to ^Z.
const NCCS = 32

type (
	ControlChar byte
	Speed       uint32
	TCFlag      uint32

	Termios struct {
		Iflag, Oflag, Cflag, Lflag TCFlag
		Line                       uint8
		CC                         [NCCS]ControlChar
		Ispeed, Ospeed             Speed
	}

	WinSize struct {
		Row, Col       uint16
		Xpixel, Ypixel uint16
	}

	TTY struct {
		f *os.File
	}
)

// termios constants
const (
	IGNBRK = TCFlag(0000001)
	BRKINT = TCFlag(0000002)
	PARMRK = TCFlag(0000010)
	ISTRIP = TCFlag(0000040)
	INLCR  = TCFlag(0000100)
	IGNCR  = TCFlag(0000200)
	ICRNL  = TCFlag(0000400)
	IXON   = TCFlag(0002000)
	OPOST  = TCFlag(0000001)
	ECHO   = TCFlag(0000010)
	ECHONL = TCFlag(0000100)
	ICANON = TCFlag(0000002)
	ISIG   = TCFlag(0000001)
	IEXTEN = TCFlag(0100000)
	CSIZE  = TCFlag(0000060)
	CS8    = TCFlag(0000060)
	PARENB = TCFlag(0000400)
	VTIME  = 5
	VMIN   = 6
)

// ioctl constants
const (
	TCGETS     = 0x5401
	TCSETS     = 0x5402
	TIOCGWINSZ = 0x5413
	TIOCSWINSZ = 0x5414
)

func New() (*TTY, error) {
	f, err := os.OpenFile("/dev/tty", os.O_RDWR, 0)
	if err != nil {
		return nil, err
	}
	t := &TTY{f: f}

	return t, nil
}

func (t *TTY) Get() (*Termios, error) {
	var ti = &Termios{}
	r1, _, errno := syscall.Syscall(syscall.SYS_IOCTL, t.f.Fd(), uintptr(TCGETS), uintptr(unsafe.Pointer(ti)))
	if errno != 0 || r1 != 0 {
		return nil, fmt.Errorf("termios.get: r1 %v, errno %v", r1, errno)
	}
	return ti, nil
}

func (t *TTY) Set(ti *Termios) error {
	if r1, _, errno := syscall.Syscall(syscall.SYS_IOCTL, t.f.Fd(), uintptr(TCSETS), uintptr(unsafe.Pointer(ti))); errno != 0 || r1 != 0 {
		return fmt.Errorf("Termios.Set: r1 %v, errno %v", r1, errno)
	}

	return nil
}

func MakeRaw(term *Termios) *Termios {
	raw := *term
	raw.Iflag &= ^(IGNBRK | BRKINT | PARMRK | ISTRIP | INLCR | IGNCR | ICRNL | IXON)
	raw.Oflag &= ^OPOST
	raw.Lflag &= ^(ECHO | ECHONL | ICANON | ISIG | IEXTEN)
	raw.Cflag &= ^(CSIZE | PARENB)
	raw.Cflag |= CS8

	raw.CC[VMIN] = 1
	raw.CC[VTIME] = 0

	return &raw
}

func GetWinSize(fd uintptr) (*WinSize, error) {
	var w = &WinSize{}
	if r1, _, errno := syscall.Syscall(syscall.SYS_IOCTL, fd, uintptr(TIOCGWINSZ), uintptr(unsafe.Pointer(w))); errno != 0 || r1 != 0 {
		return nil, fmt.Errorf("WinSize.Get: r1 %v, errno %v", r1, errno)
	}

	return w, nil
}

func (t *TTY) GetWinSize() (*WinSize, error) {
	return GetWinSize(t.f.Fd())
}

func SetWinSize(fd uintptr, w *WinSize) error {
	if r1, _, errno := syscall.Syscall(syscall.SYS_IOCTL, fd, uintptr(TIOCSWINSZ), uintptr(unsafe.Pointer(w))); errno != 0 || r1 != 0 {
		return fmt.Errorf("WinSize.Set: r1 %v, errno %v", r1, errno)
	}

	return nil
}

func (t *TTY) SetWinSize(w *WinSize) error {
	return SetWinSize(t.f.Fd(), w)
}
