// Copyright 2012-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// dmesg reads the Linux system log.
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
	"fmt"
	"io"
	"log"
	"os"

	"golang.org/x/sys/unix"
)

func run(out io.Writer, args []string) error {
	var clr, readClr bool

	f := flag.NewFlagSet(args[0], flag.ContinueOnError)
	f.BoolVar(&clr, "clear", false, "Clear the log")
	f.BoolVar(&readClr, "read-clear", false, "Clear the log after printing")
	f.Parse(args[1:])

	if clr && readClr {
		return fmt.Errorf("cannot specify both -clear and -read-clear:%w", os.ErrInvalid)
	}

	level := unix.SYSLOG_ACTION_READ_ALL
	if clr {
		level = unix.SYSLOG_ACTION_CLEAR
	}
	if readClr {
		level = unix.SYSLOG_ACTION_READ_CLEAR
	}

	b := make([]byte, 256*1024)
	amt, err := unix.Klogctl(level, b)
	if err != nil {
		return fmt.Errorf("syslog failed: %w", err)
	}

	_, err = out.Write(b[:amt])
	return err
}

func main() {
	if err := run(os.Stdout, os.Args); err != nil {
		log.Fatal(err)
	}
}
