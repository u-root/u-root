// Copyright 2016-2020 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build !plan9
// +build !plan9

// sync command in Go.
//
// Synopsis:
//		sync [-df] [FILE]
//

package main

import (
	"flag"
	"log"
	"os"
	"syscall"

	"github.com/u-root/u-root/pkg/uroot/util"
	"golang.org/x/sys/unix"
)

var (
	data       = flag.Bool("data", false, "sync file data, no metadata")
	filesystem = flag.Bool("filesystem", false, "commit filesystem caches to disk")
)

var usage = "Usage: %s [OPTION] [FILE]...\n"

func init() {
	util.Usage(usage)
}

func doSyscall(syscallNum uintptr, args []string) error {
	for _, fileName := range args {
		f, err := os.OpenFile(fileName, syscall.O_RDONLY|syscall.O_NOCTTY|syscall.O_CLOEXEC, 0o644)
		if err != nil {
			return err
		}
		if _, _, err = syscall.Syscall(syscallNum, uintptr(f.Fd()), 0, 0); err.(syscall.Errno) != 0 {
			return err
		}
		f.Close()
	}
	return nil
}

func sync(args []string) error {
	switch {
	case *data:
		return doSyscall(unix.SYS_FDATASYNC, args)
	case *filesystem:
		return doSyscall(unix.SYS_SYNCFS, args)
	default:
		syscall.Sync()
		return nil
	}
}

func main() {
	flag.Parse()
	if err := sync(flag.Args()); err != nil {
		log.Fatal(err)
	}
}
