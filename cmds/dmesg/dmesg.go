// Copyright 2012,2016 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//dmesg reads the system log.
package main

import (
	"flag"
	"log"
	"os"
	"unsafe"

	"golang.org/x/sys/unix"
)

var clearSyslog = flag.Bool("c", false, "Clear the log")

func main() {
	flag.Parse()
	l := uintptr(3)
	if *clearSyslog {
		l = 4
	}
	b := make([]byte, 256*1024)
	if amt, _, err := unix.Syscall(unix.SYS_SYSLOG, l, uintptr(unsafe.Pointer(&b[0])), uintptr(len(b))); err == 0 {
		os.Stdout.Write(b[:amt])
	} else {
		log.Fatalf("syslog failed: %v", err)
	}
}
