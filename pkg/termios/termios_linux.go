// Copyright 2015-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package termios

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"

	"golang.org/x/sys/unix"
)

// TTY is a wrapper that only allows Read and Write.
type TTY struct {
	f *os.File
}

func New() (*TTY, error) {
	f, err := os.OpenFile("/dev/tty", os.O_RDWR, 0)
	if err != nil {
		return nil, err
	}
	return &TTY{f: f}, nil
}

func NewTTYS(port string) (*TTY, error) {
	f, err := os.OpenFile(filepath.Join("/dev", port), unix.O_RDWR|unix.O_NOCTTY|unix.O_NONBLOCK, 0620)
	if err != nil {
		return nil, err
	}
	return &TTY{f: f}, nil
}

func (t *TTY) Read(b []byte) (int, error) {
	return t.f.Read(b)
}

func (t *TTY) Write(b []byte) (int, error) {
	return t.f.Write(b)
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

func (t *TTY) Ctty(c *exec.Cmd) {
	c.Stdin, c.Stdout, c.Stderr = t.f, t.f, t.f
	if c.SysProcAttr == nil {
		c.SysProcAttr = &syscall.SysProcAttr{}
	}
	c.SysProcAttr.Setctty = true
	c.SysProcAttr.Setsid = true
	c.SysProcAttr.Ctty = int(t.f.Fd())
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

// MakeSerialBaud updates the Termios to set the baudrate
func MakeSerialBaud(term *unix.Termios, baud int) (*unix.Termios, error) {
	t := *term
	rate, ok := baud2unixB[baud]
	if !ok {
		return nil, fmt.Errorf("%d: Unrecognized baud rate", baud)
	}

	t.Cflag &^= unix.CBAUD
	t.Cflag |= rate
	t.Ispeed = rate
	t.Ospeed = rate

	return &t, nil
}

// MakeSerialDefault updates the Termios to typical serial configuration:
// - Ignore all flow control (modem, hardware, software...)
// - Translate carriage return to newline on input
// - Enable canonical mode: Input is available line by line, with line editing
//   enabled (ERASE, KILL are supported)
// - Local ECHO is added (and handled by line editing)
// - Map newline to carriage return newline on output
func MakeSerialDefault(term *unix.Termios) *unix.Termios {
	t := *term
	/* Clear all except baud, stop bit and parity settings */
	t.Cflag &= unix.CBAUD | unix.CSTOPB | unix.PARENB | unix.PARODD
	/* Set: 8 bits; ignore Carrier Detect; enable receive */
	t.Cflag |= unix.CS8 | unix.CLOCAL | unix.CREAD
	t.Iflag = unix.ICRNL
	t.Lflag = unix.ICANON | unix.ISIG | unix.ECHO | unix.ECHOE | unix.ECHOK | unix.ECHOKE | unix.ECHOCTL
	/* non-raw output; add CR to each NL */
	t.Oflag = unix.OPOST | unix.ONLCR
	/* reads will block only if < 1 char is available */
	t.Cc[unix.VMIN] = 1
	/* no timeout (reads block forever) */
	t.Cc[unix.VTIME] = 0
	t.Line = 0

	return &t
}
