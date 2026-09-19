// Copyright 2026  the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// slattach is a u-root implementation of the slattach command, for SLIP.
// usage: slattach ttyX command...
// ttyX is typically something like ttyUSB0 or ttyS0.
// The commands are run once the tty is attached and set in raw mode.
// They are particularly useful if the serial console is to be used as the
// SLIP device.
// Example:
// on an embedded system with u-root in firmware:
// slattach ttyS0 'ip addr add 192.168.4.1/32 peer 192.168.4.2 dev sl0' 'ip link set dev sl0 up'
// on a host system connected to the embeeded with USB uarts.
// slattach ttyUSB0 'ip addr add 192.168.4.2/32 peer 192.168.4.1 dev sl0' 'ip link set dev sl0 up'
// You can also arrange the commands to timeout slattach, which is handy for debugging:
// sudo slattach ttyUSB0 'ip addr add 192.168.0.1/24 dev sl0' 'ip link set dev sl0 up' 'sleep 50' 'blarg'
//
// You can skip the commands if you don't need them.
package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"
	"unsafe"

	"github.com/u-root/u-root/pkg/termios"
	"golang.org/x/sys/unix"
)

const (
	UNIX = iota
	SLIP
)

func run(stdin io.Reader, stdout, stderr io.Writer, args []string) error {
	f := flag.NewFlagSet(args[0], flag.ContinueOnError)
	f.SetOutput(stderr)
	if err := f.Parse(args[1:]); err != nil {
		return err
	}

	if f.NArg() < 1 {
		return fmt.Errorf("usage: %v ttyX", os.Args[0])
	}

	disc := SLIP

	dev := f.Arg(0)
	tty, err := termios.NewTTYS(dev)
	if err != nil {
		return fmt.Errorf("failed to open TTY: %w", err)
	}

	oldTermios, err := tty.Raw()
	if err != nil {
		return fmt.Errorf("setting raw: %w", err)
	}
	defer tty.Set(oldTermios)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	fd := int(tty.Fd())
	if _, _, errno := unix.Syscall(unix.SYS_IOCTL, uintptr(fd), uintptr(unix.TIOCSETD), uintptr(unsafe.Pointer(&disc))); errno != 0 {
		return fmt.Errorf("set disc: %w", errno)
	}

	defer func() {
		disc = UNIX
		_, _, _ = unix.Syscall(unix.SYS_IOCTL, uintptr(fd), uintptr(unix.TIOCSETD), uintptr(unsafe.Pointer(&disc)))
	}()

	var ifname [16]byte
	if _, _, errno := unix.Syscall(unix.SYS_IOCTL, uintptr(fd), uintptr(unix.SIOCGIFNAME), uintptr(unsafe.Pointer(&ifname[0]))); errno != 0 {
		return fmt.Errorf("attached %v (SLIP), but failed to query interface name: %w", tty, errno)
	}

	n := bytes.TrimRight(ifname[:], "\x00")
	fmt.Fprintf(stderr, "Attached %s to network interface %s\n", dev, n)

	for _, c := range f.Args()[1:] {
		args := strings.Fields(c)
		cmd := exec.Command(args[0], args[1:]...)
		b, err := cmd.CombinedOutput()
		if err != nil {
			return fmt.Errorf("%v(%v):%s:%w", c, cmd, string(b), err)
		}
	}
	<-sigChan

	fmt.Fprintln(stderr, "\nDetached SLIP interface.")
	return nil
}

func main() {
	if err := run(os.Stdin, os.Stdout, os.Stderr, os.Args); err != nil {
		log.Fatalf("run:%v", err)
	}
}
