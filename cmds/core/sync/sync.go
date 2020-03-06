// Copyright 2016-2020 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// sync command in Go.
//
// Synopsis:
//		sync [-df] [FILE]
//

package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"syscall"

	"golang.org/x/sys/unix"
)

var (
	data       = flag.Bool("data", false, "sync file data, no metadata")
	filesystem = flag.Bool("filesystem", false, "commit filesystem caches to disk")
)

func init() {
	flag.BoolVar(data, "d", false, "")
	flag.BoolVar(filesystem, "f", false, "")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [OPTION] [FILE]...\n", os.Args[0])
		flag.PrintDefaults()
	}
}

func doSyscall(syscallNum uintptr) {
	for _, fileName := range flag.Args() {
		f, err := os.OpenFile(fileName, syscall.O_RDONLY|syscall.O_NOCTTY|syscall.O_CLOEXEC, 0644)
		if err != nil {
			log.Fatal(err)
		}
		_, _, err = syscall.Syscall(syscallNum, uintptr(f.Fd()), 0, 0)
		if err != nil {
			log.Fatal(err)
		}
		f.Close()
	}
}

func main() {
	flag.Parse()
	switch {
	case *data:
		doSyscall(unix.SYS_FDATASYNC)
	case *filesystem:
		doSyscall(unix.SYS_SYNCFS)
	default:
		syscall.Sync()
		os.Exit(0)
	}
}
