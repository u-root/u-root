// Copyright 2015-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// console implements a basic console. It establishes a pair of files
// to read from, the default being a UART at 0x3f8, but an alternative
// being just stdin and stdout. It will also set up a root file system
// using uroot.Rootfs, although this can be disabled as well.
// Console uses a Go version of fork_pty to start up a shell, default
// /ubin/rush. Console runs until the shell exits and then exits itself.
package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strconv"
	"syscall"
	"unsafe"

	"github.com/u-root/u-root/uroot"
)

var (
	serial    = flag.Bool("serial", true, "use 0x3f8 for stdin")
	setupRoot = flag.Bool("setuproot", true, "Set up a root file system")
)

// pty support. We used to import github.com/kr/pty but what we need is not that complex.
// Thanks to keith rarick for these functions.
func ptsopen() (pty, tty *os.File, slavename string, err error) {
	p, err := os.OpenFile("/dev/ptmx", os.O_RDWR, 0)
	if err != nil {
		return
	}

	err = ptsunlock(p)
	if err != nil {
		return
	}

	slavename, err = ptsname(p)
	if err != nil {
		return
	}

	// It can take a non-zero time for a pts to appear, it seems. 
	// Ten tries is reported to be far more than enough.
	for i := 0; i < 10; i++ {
		fi, err := os.Stat(slavename)
		if err == nil {
			fmt.Printf("stat of %v ok after %d iterations: %v", slavename, i, fi)
			break
		}
		fmt.Printf("stat of %v: %v", slavename, err)
	}
	t, err := os.OpenFile(slavename, os.O_RDWR|syscall.O_NOCTTY, 0)
	if err != nil {
		return
	}
	return p, t, slavename, nil
}

func ptsname(f *os.File) (string, error) {
	var n uintptr
	_, _, err := syscall.Syscall(syscall.SYS_IOCTL, f.Fd(), syscall.TIOCGPTN, uintptr(unsafe.Pointer(&n)))
	if err != 0 {
		return "", err
	}
	return "/dev/pts/" + strconv.Itoa(int(n)), nil
}

func ptsunlock(f *os.File) error {
	var u uintptr
	// use TIOCSPTLCK with a zero valued arg to clear the slave pty lock
	_, _, err := syscall.Syscall(syscall.SYS_IOCTL, f.Fd(), syscall.TIOCSPTLCK, uintptr(unsafe.Pointer(&u)))
	if err != 0 {
		return err
	}
	return nil
}

func main() {
	fmt.Printf("console -- starting")
	flag.Parse()

	a := flag.Args()
	if len(a) < 1 {
		a = []string{"/ubin/rush"}
	}

	// Make a good faith effort to set up root. This being
	// a kind of init program, we do our best and keep going.
	if *setupRoot {
		uroot.Rootfs()
	}

	in, out := io.Reader(os.Stdin), io.Writer(os.Stdout)

	if *serial {
		if err := openUART(); err != nil {
			// can't happen but ... maybe.
			fmt.Printf("Sorry, can't get a uart: %v", err)
			os.Exit(1)
		}
		in, out = uart{}, uart{}
	}

	ptm, pts, sname, err := ptsopen()

	if err != nil {
		fmt.Printf("ptsopen: %v; are you root?", err)
		os.Exit(1)
	}
	fmt.Printf("console: ptm %v pts %v sname %v ", ptm, pts, sname)
	c := exec.Command(a[0], a[1:]...)
	c.Stdout = pts
	c.Stdin = pts
	c.Stderr = c.Stdout
	c.SysProcAttr = &syscall.SysProcAttr{}
	c.SysProcAttr.Setctty = true
	c.SysProcAttr.Setsid = true
	err = c.Start()
	if err != nil {
		fmt.Printf("Can't start %v: %v", a, err)
		os.Exit(1)
	}
	kid := c.Process.Pid
	fmt.Printf("Started %d\n", kid)

	raw()

	go func() {
		io.Copy(out, ptm)
		fmt.Printf("kid stdout: done\n")
		os.Exit(1)
	}()
	go func() {
		var data = make([]byte, 1)
		for {
			if _, err := in.Read(data); err != nil {
				fmt.Printf("kid stdin: done\n")
			}
			if data[0] == '\r' {
				if _, err := out.Write(data); err != nil {
					fmt.Printf("error on echo %v: %v", data, err)
				}
				data[0] = '\n'
			}
			if _, err := ptm.Write(data); err != nil {
				fmt.Printf("Error writing input to ptm: %v: give up\n", err)
				os.Exit(1)
			}
		}
	}()

	err = c.Wait()
	fmt.Printf("kid: done %v", err)
	os.Exit(1)
}
