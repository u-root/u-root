// Copyright 2013-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build !plan9 && !windows

// mkfifo creates a named pipe.
//
// Synopsis:
//
//	mkfifo [OPTIONS] NAME...
//
// Options:
//
//	-m: mode (default 0600)
package main

import (
	"flag"
	"log"
	"os"

	"golang.org/x/sys/unix"
)

const (
	defaultMode = 0o660 | unix.S_IFIFO
	cmd         = "mkfifo [-m] NAME..."
)

var mode = flag.Int("mode", defaultMode, "Mode to create fifo")

func init() {
	defUsage := flag.Usage
	flag.Usage = func() {
		os.Args[0] = cmd
		defUsage()
	}
}

func main() {
	flag.Parse()

	if flag.NArg() < 1 {
		log.Fatal("please provide a path, or multiple, to create a fifo")
	}

	for _, path := range flag.Args() {
		if err := unix.Mkfifo(path, uint32(*mode)); err != nil {
			log.Fatalf("Error while creating fifo, %v", err)
		}
	}
}
