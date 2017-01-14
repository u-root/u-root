// Copyright 2012-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Setup loop devices.
//
// Synopsis:
//     losetup [-Ad] FILE
//     losetup [-Ad] DEV FILE
//
// Options:
//     -A: pick any device
//     -d: detach the device
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"syscall"
)

const (
	/*
	 * IOCTL commands --- we will commandeer 0x4C ('L')
	 */
	LOOP_SET_CAPACITY = 0x4C07
	LOOP_CHANGE_FD    = 0x4C06
	LOOP_GET_STATUS64 = 0x4C05
	LOOP_SET_STATUS64 = 0x4C04
	LOOP_GET_STATUS   = 0x4C03
	LOOP_SET_STATUS   = 0x4C02
	LOOP_CLR_FD       = 0x4C01
	LOOP_SET_FD       = 0x4C00
	LO_NAME_SIZE      = 64
	LO_KEY_SIZE       = 32
	/* /dev/loop-control interface */
	LOOP_CTL_ADD      = 0x4C80
	LOOP_CTL_REMOVE   = 0x4C81
	LOOP_CTL_GET_FREE = 0x4C82
	SYS_ioctl         = 16
)

var (
	anyLoop = flag.Bool("A", true, "Pick any device")
	detach  = flag.Bool("d", false, "Detach the device")
	l       = log.New(os.Stdout, "tcz: ", 0)
)

// consider making this a goroutine which pushes the string down the channel.
func findloop() (name string, err error) {
	cfd, err := syscall.Open("/dev/loop-control", syscall.O_RDWR, 0)
	if err != nil {
		log.Fatalf("/dev/loop-control: %v", err)
	}
	defer syscall.Close(cfd)
	a, b, errno := syscall.Syscall(SYS_ioctl, uintptr(cfd), LOOP_CTL_GET_FREE, 0)
	if errno != 0 {
		log.Fatalf("ioctl: %v\n", err)
	}
	log.Printf("a %v b %v err %v\n", a, b, err)
	name = fmt.Sprintf("/dev/loop%d", a)
	return name, nil
}

func main() {
	flag.Parse()
	args := flag.Args()
	if *detach {
		l.Fatalf("detach: not yet")
		os.Exit(1)
	}
	var file, dev string
	var err error
	if len(args) == 1 {
		dev, err = findloop()
		if err != nil {
			l.Fatalf("can't find a loop: %v\n", err)
			os.Exit(1)
		}
		file = args[0]
	} else if len(args) == 2 {
		dev = args[0]
		file = args[1]
	} else {
		l.Fatalf("usage\n")
		os.Exit(1)
	}

	ffd, err := syscall.Open(file, syscall.O_RDONLY, 0)
	if err != nil {
		l.Fatalf("file: %v, %v\n", file, err)
	}
	lfd, err := syscall.Open(dev, syscall.O_RDONLY, 0)
	if err != nil {
		l.Fatalf("dev: %v, %v\n", dev, err)
	}
	l.Printf("ffd %v lfd %v\n", ffd, lfd)
	a, b, errno := syscall.Syscall(SYS_ioctl, uintptr(lfd), LOOP_SET_FD, uintptr(ffd))
	if errno != 0 {
		l.Fatalf("loop set fd ioctl: %v, %v, %v\n", a, b, errno)
	}
}
