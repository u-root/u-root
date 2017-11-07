// Copyright 2012-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Read the system log.
//
// Synopsis:
//     dmesg [-clear|-read-clear]
//
// Options:
//     -clear: clear the log
//     -read-clear: clear the log after printing
package main

import (
	"flag"
	"log"
	"os"
	"syscall"
	"unsafe"
)

const (
	_SYSLOG_ACTION_READ       = 2
	_SYSLOG_ACTION_READ_ALL   = 3
	_SYSLOG_ACTION_READ_CLEAR = 4
	_SYSLOG_ACTION_CLEAR      = 5
)

var (
	clear     bool
	readClear bool
)

func init() {
	flag.BoolVar(&clear, "clear", false, "Clear the log")
	flag.BoolVar(&readClear, "read-clear", false, "Clear the log after printing")
	flag.BoolVar(&readClear, "c", false, "Clear the log after printing")
}

func main() {
	flag.Parse()
	if clear && readClear {
		log.Fatalf("cannot specify both -clear and -read-clear")
	}

	level := uintptr(_SYSLOG_ACTION_READ_ALL)
	if clear {
		level = _SYSLOG_ACTION_CLEAR
	}
	if readClear {
		level = _SYSLOG_ACTION_READ_CLEAR
	}

	b := make([]byte, 256*1024)
	amt, _, err := syscall.Syscall(syscall.SYS_SYSLOG, level, uintptr(unsafe.Pointer(&b[0])), uintptr(len(b)))
	if err != 0 {
		log.Fatalf("syslog failed: %v", err)
	}

	os.Stdout.Write(b[:amt])
}
