// Copyright 2015-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package termios

import (
	"os"

	"golang.org/x/sys/unix"
)

type TTY struct {
	f *os.File
}

func New() (*TTY, error) {
	f, err := os.OpenFile("/dev/tty", os.O_RDWR, 0)
	if err != nil {
		return nil, err
	}
	t := &TTY{f: f}

	return t, nil
}

func GetTermios(fd uintptr) (*unix.Termios, error) {
	return unix.IoctlGetTermios(int(fd), unix.TCGETS)
}

func (t *TTY) Get() (*unix.Termios, error) {
	return GetTermios(t.f.Fd())
}

func SetTermios(fd uintptr, ti *unix.Termios) error {
	return unix.IoctlSetTermios(int(fd), unix.TCSETS, ti)
}

func (t *TTY) Set(ti *unix.Termios) error {
	return SetTermios(t.f.Fd(), ti)
}

func MakeRaw(term *unix.Termios) *unix.Termios {
	raw := *term
	raw.Iflag &^= unix.IGNBRK | unix.BRKINT | unix.PARMRK | unix.ISTRIP | unix.INLCR | unix.IGNCR | unix.ICRNL | unix.IXON
	raw.Oflag &^= unix.OPOST
	raw.Lflag &^= unix.ECHO | unix.ECHONL | unix.ICANON | unix.ISIG | unix.IEXTEN
	raw.Cflag &^= unix.CSIZE | unix.PARENB
	raw.Cflag |= unix.CS8

	raw.Cc[unix.VMIN] = 1
	raw.Cc[unix.VTIME] = 0

	return &raw
}

func GetWinSize(fd uintptr) (*unix.Winsize, error) {
	return unix.IoctlGetWinsize(int(fd), unix.TIOCGWINSZ)
}

func (t *TTY) GetWinSize() (*unix.Winsize, error) {
	return GetWinSize(t.f.Fd())
}

func SetWinSize(fd uintptr, w *unix.Winsize) error {
	return unix.IoctlSetWinsize(int(fd), unix.TIOCSWINSZ, w)
}

func (t *TTY) SetWinSize(w *unix.Winsize) error {
	return SetWinSize(t.f.Fd(), w)
}
