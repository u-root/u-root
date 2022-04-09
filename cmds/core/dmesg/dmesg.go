// Copyright 2012-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build !plan9 && !windows

// dmesg reads the system log.
//
// Synopsis:
//     dmesg [-clear|-read-clear]
//
// Options:
//     -clear: clear the log
//     -read-clear: clear the log after printing
package main

import (
	"fmt"
	"io"
	"log"
	"os"

	flag "github.com/spf13/pflag"
	"golang.org/x/sys/unix"
)

var (
	clear     = flag.Bool("clear", false, "Clear the log")
	readClear = flag.BoolP("read-clear", "c", false, "Clear the log after printing")
)

func dmesg(writer io.Writer) error {
	if *clear && *readClear {
		return fmt.Errorf("cannot specify both -clear and -read-clear")
	}

	level := unix.SYSLOG_ACTION_READ_ALL
	if *clear {
		level = unix.SYSLOG_ACTION_CLEAR
	}
	if *readClear {
		level = unix.SYSLOG_ACTION_READ_CLEAR
	}

	b := make([]byte, 256*1024)
	amt, err := unix.Klogctl(level, b)
	if err != nil {
		return fmt.Errorf("syslog failed: %v", err)
	}

	writer.Write(b[:amt])
	return nil
}

func main() {
	flag.Parse()
	if err := dmesg(os.Stdout); err != nil {
		log.Fatal(err)
	}
}
