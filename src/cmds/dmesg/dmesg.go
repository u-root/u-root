// Copyright 2012 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
Cat reads each file from its arguments in sequence and writes it on the standard output.
*/

package main

import (
	"flag"
	"log"
	"os"
	"syscall"
	"unsafe"
)

var clearSyslog = flag.Bool("c", false, "Clear the log")

func main() {
	flag.Parse()
	l := uintptr(3)
	if *clearSyslog {
		l = 4
	}
	b := make([]byte, 256*1024)
	if amt, _, err := syscall.Syscall(syscall.SYS_SYSLOG, l, uintptr(unsafe.Pointer(&b[0])), uintptr(len(b))); err == 0 {
		os.Stdout.Write(b[:amt])
	} else {
		log.Fatalf("syslog failed: %v", err)
	}
}
