// Copyright 2012-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build !tinygo

package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"unsafe"

	"golang.org/x/sys/unix"
)

var (
	ttypgrp int
	ttyf    *os.File
)

// tty does whatever needs to be done to set up a tty for GOOS.
func tty() {
	var err error

	sigs := make(chan os.Signal, 512)
	signal.Notify(sigs, os.Interrupt)
	signal.Ignore(unix.SIGTTOU)
	go func() {
		for i := range sigs {
			fmt.Println(i)
		}
	}()

	// N.B. We can continue to use this file, in the foreground function,
	// but the runtime closes it on exec for us.
	ttyf, err = os.OpenFile("/dev/tty", os.O_RDWR, 0)
	if err != nil {
		log.Printf("rush: Can't open a console; no job control in this session")
		return
	}
	// Get the current pgrp, and the pgrp on the tty.
	// get current pgrp
	ttypgrp, err = unix.IoctlGetInt(int(ttyf.Fd()), unix.TIOCGPGRP)
	if err != nil {
		log.Printf("Can't get foreground: %v", err)
		ttyf.Close()
		ttyf = nil
		ttypgrp = 0
	}
}

func foreground() {
	// Place process group in foreground.
	if ttypgrp != 0 {
		_, _, errno := unix.RawSyscall(unix.SYS_IOCTL, ttyf.Fd(), uintptr(unix.TIOCSPGRP), uintptr(unsafe.Pointer(&ttypgrp)))
		if errno != 0 {
			log.Printf("rush pid %v: Can't set foreground to %v: %v", os.Getpid(), ttypgrp, errno)
		}
	}
}
