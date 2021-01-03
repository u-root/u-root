// Copyright 2020 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"log"
	"os"
)

var (
	ttyf    *os.File
	ttyc    *os.File
	newline = "\n"
)

// tty does whatever needs to be done to set up a tty for GOOS.
func tty() {
	var err error

	// N.B. We can continue to use this file, in the foreground function,
	// but the runtime closes it on exec for us.
	ttyf, err = os.OpenFile("/dev/cons", os.O_RDWR, 0)
	if err != nil {
		log.Fatalf("ash: Can't open a console; no job control in this session")
	}
	ttyc, err = os.OpenFile("/dev/consctl", os.O_WRONLY, 0)
	if err != nil {
		log.Fatalf("ash: Can't open a consctl; no job control in this session")
	}
}

func foreground() {
}

func preExec(c *Command) {
}

func unshare(c *Command) {
}
