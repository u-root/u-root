// Copyright 2012-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"unsafe"
)

var (
	ttypgrp, pgrpid int
	ttyf            *os.File
)

// tty does whatever needs to be done to set up a tty for GOOS.
func tty() {
	var err error

	// N.B. We can continue to use this file, in the foreground function,
	// but the runtime closes it on exec for us.
	if ttyf, err = os.OpenFile("/dev/tty", os.O_RDWR, 0); err != nil {
		log.Printf("rush: Can't open a console; no job control in this session")
		return
	}
	sigs := make(chan os.Signal, 512)
	signal.Notify(sigs, os.Interrupt)
	go func() {
		for i := range sigs {
			fmt.Println(i)
		}
	}()

	// Get the current pgrp, and the pgrp on the tty.
	// get current pgrp
	r1, r2, errno := syscall.RawSyscall(syscall.SYS_IOCTL, ttyf.Fd(), uintptr(syscall.TIOCGPGRP), uintptr(unsafe.Pointer(&pgrpid)))
	if errno != 0 {
		log.Printf("Can't set foreground: %v, %v, %v", r1, r2, errno)
	}

	pgrpid, err = syscall.Getpgid(0)
	if err != nil {
		log.Printf("rush: can't get my own pgid, no job control")
		return
	}
}

func foreground() {
	// Place process group in foreground.
	if pgrpid != 0 {
		if err := syscall.Setpgid(0, pgrpid); err != nil {
			log.Printf("Warning: failed to set pgid: %v", err)
		}
	}
	if ttyf.Fd() >= 0 {
		r1, r2, errno := syscall.RawSyscall(syscall.SYS_IOCTL, ttyf.Fd(), uintptr(syscall.TIOCSPGRP), uintptr(unsafe.Pointer(&pgrpid)))
		if errno != 0 {
			log.Printf("Can't set foreground: %v, %v, %v", r1, r2, errno)
		}
	}

}
