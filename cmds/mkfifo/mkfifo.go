// Copyright 2013-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"log"
	"syscall"

	"github.com/u-root/u-root/pkg/mkfifo"
)

var (
	defaultMode = syscall.S_IRUSR | syscall.S_IWUSR | syscall.S_IRGRP | syscall.S_IWGRP | syscall.S_IROTH | syscall.S_IWOTH
	mode        = flag.Int("mode", defaultMode, "Mode to create fifo")
)

func main() {
	flag.Parse()

	if flag.NArg() < 1 {
		log.Fatal("please provide a path, or multiple, to create a fifo")
	}

	mk := mkfifo.Mkfifo{Paths: flag.Args(), Mode: uint32(*mode)}

	if err := mk.Exec(); err != nil {
		log.Fatalf("Error while creating fifo, %v", err)
	}
}
