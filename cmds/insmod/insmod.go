// Copyright 2012-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"os"
	"io/ioutil"
	"syscall"
	"unsafe"
)

func main() {
	var options string

	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "insmod: ERROR: missing filename.\n")
		os.Exit(1)
	}

	// get filename from argv[1]
	filename := os.Args[1]

	// Everything else is module options
	for i := 2; i < len(os.Args); i++ {
		options = options + os.Args[i] + " "
	}

	// read file into memory
	file, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "insmod: can't read '%s': %v\n", filename, err)
		os.Exit(1)
	}

	// call SYS_INIT_MODULE with file, length, and options
	ret, _, err := syscall.Syscall(syscall.SYS_INIT_MODULE, uintptr(unsafe.Pointer(&file)), uintptr(len(file)), uintptr(unsafe.Pointer(&options)))
	if err != nil {
		fmt.Fprintf(os.Stderr, "insmod: error inserting '%s': %v %v\n", filename, ret, err)
	}

	os.Exit(0)
}
