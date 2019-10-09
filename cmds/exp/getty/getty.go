// Copyright 2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// getty Open a TTY and invoke a shell
// There are no special options and no login support
// Also getty exits after starting the shell so if one exits the shell, there
// is no more shell!
//
// Synopsis:
//     getty <port> <baud> [term]
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"syscall"

	"github.com/u-root/u-root/pkg/upath"
	"golang.org/x/sys/unix"
)

var (
	verbose = flag.Bool("v", false, "verbose log")
	debug   = func(string, ...interface{}) {}
	cmdList []string
	envs    []string
)

func init() {
	r := upath.UrootPath
	cmdList = []string{
		r("/bin/defaultsh"),
		r("/bin/sh"),
	}
}

func openPort(port string, baud int) (*os.File, error) {
	var bauds = map[int]uint32{
		0:       0,
		50:      unix.B50,
		75:      unix.B75,
		110:     unix.B110,
		134:     unix.B134,
		150:     unix.B150,
		200:     unix.B200,
		300:     unix.B300,
		600:     unix.B600,
		1200:    unix.B1200,
		1800:    unix.B1800,
		2400:    unix.B2400,
		4800:    unix.B4800,
		9600:    unix.B9600,
		19200:   unix.B19200,
		38400:   unix.B38400,
		57600:   unix.B57600,
		115200:  unix.B115200,
		230400:  unix.B230400,
		460800:  unix.B460800,
		500000:  unix.B500000,
		576000:  unix.B576000,
		921600:  unix.B921600,
		1000000: unix.B1000000,
		1152000: unix.B1152000,
		1500000: unix.B1500000,
		2000000: unix.B2000000,
		2500000: unix.B2500000,
		3000000: unix.B3000000,
		3500000: unix.B3500000,
		4000000: unix.B4000000,
	}
	rate, ok := bauds[baud]
	if !ok {
		return nil, fmt.Errorf("Unrecognized baud rate")
	}

	f, err := os.OpenFile("/dev/"+port, unix.O_RDWR|unix.O_NOCTTY|unix.O_NONBLOCK, 0620)
	if err != nil {
		return nil, err
	}

	fd := int(f.Fd())
	t, err := unix.IoctlGetTermios(fd, unix.TCGETS)
	if err != nil {
		return nil, err
	}

	if baud != 0 {
		t.Cflag &^= unix.CBAUD
		t.Cflag |= rate
		t.Ispeed = rate
		t.Ospeed = rate
	}

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

	err = unix.IoctlSetTermios(fd, unix.TCSETS, t)
	if err != nil {
		return nil, err
	}

	return f, nil
}

func main() {
	flag.Parse()

	if *verbose {
		debug = log.Printf
	}

	port := flag.Arg(0)
	baud, err := strconv.Atoi(flag.Arg(1))
	if err != nil {
		baud = 0
	}
	term := flag.Arg(2)

	fd, err := openPort(port, baud)
	if err != nil {
		debug("Unable to open port %s: %v", port, err)
	}

	// Output the u-root banner
	log.New(fd, "", log.LstdFlags).Printf("Welcome to u-root!")
	fmt.Fprintln(fd, `                              _`)
	fmt.Fprintln(fd, `   _   _      _ __ ___   ___ | |_`)
	fmt.Fprintln(fd, `  | | | |____| '__/ _ \ / _ \| __|`)
	fmt.Fprintln(fd, `  | |_| |____| | | (_) | (_) | |_`)
	fmt.Fprintln(fd, `   \__,_|    |_|  \___/ \___/ \__|`)
	fmt.Fprintln(fd)

	log.SetPrefix("getty: ")

	if term != "" {
		err = os.Setenv("TERM", term)
		if err != nil {
			debug("Unable to set 'TERM=%s': %v", port, err)
		}
	}
	envs = os.Environ()
	debug("envs %v", envs)

	for _, v := range cmdList {
		debug("Trying to run %v", v)
		if _, err := os.Stat(v); os.IsNotExist(err) {
			debug("%v", err)
			continue
		}

		cmd := exec.Command(v)
		cmd.Env = envs
		cmd.Stdin, cmd.Stdout, cmd.Stderr = fd, fd, fd
		cmd.SysProcAttr = &syscall.SysProcAttr{Setctty: true, Setsid: true, Ctty: int(fd.Fd())}
		debug("running %v", cmd)
		if err := cmd.Start(); err != nil {
			log.Printf("Error starting %v: %v", v, err)
			continue
		}
		if err := cmd.Process.Release(); err != nil {
			log.Printf("Error releasing process %v:%v", v, err)
		}
		// stop after first valid command
		return
	}
	log.Printf("No suitable executable found in %+v", cmdList)
}
