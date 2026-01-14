// Copyright 2012-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// dmesg reads the system log.
//
// Synopsis:
//
//	dmesg [-clear|-read-clear]
//
// Options:
//
//	-clear: clear the log
//	-read-clear: clear the log after printing
package main

import (
	"flag"
	"log"
	"os"

	"golang.org/x/sys/unix"
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

	level := unix.SYSLOG_ACTION_READ_ALL
	if clear {
		level = unix.SYSLOG_ACTION_CLEAR
	}
	if readClear {
		level = unix.SYSLOG_ACTION_READ_CLEAR
	}

	b := make([]byte, 256*1024)
	amt, err := unix.Klogctl(level, b)
	if err != nil {
		log.Fatalf("syslog failed: %v", err)
	}

	os.Stdout.Write(b[:amt])
}
