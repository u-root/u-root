// Copyright 2012-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Insert a module into the Linux kernel
//
// Synopsis:
//	insmod [filename] [module options...]
//
// Description:
//	insmod is a clone of insmod(8)
package main

import (
	"io/ioutil"
	"log"
	"os"
	"strings"
	"syscall"
	"unsafe"
)

func main() {
	var options string

	if len(os.Args) < 2 {
		log.Fatalf("insmod: ERROR: missing filename.\n")
	}

	// get filename from argv[1]
	filename := os.Args[1]

	// Everything else is module options
	options = strings.Join(os.Args[2:], " ")

	file, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatalf("insmod: can't read '%s': %v\n", filename, err)
	}

	// call SYS_INIT_MODULE with file, length, and options
	ret, _, err := syscall.Syscall(syscall.SYS_INIT_MODULE, uintptr(unsafe.Pointer(&file[0])), uintptr(len(file)), uintptr(unsafe.Pointer(&[]byte(options)[0])))
	if ret != 0 {
		log.Fatalf("insmod: error inserting '%s': %v %v\n", filename, ret, err)
	}
}
