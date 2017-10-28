// Copyright 2012-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Remove a module from the Linux kernel
//
// Synopsis:
//	rmmod name
//
// Description:
//	rmmod is a clone of rmmod(8)
//
// Author:
//     Roland Kammerer <dev.rck@gmail.com>
package main

import (
	"log"
	"os"
	"syscall"
	"unsafe"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatalf("rmmod: ERROR: missing module name.\n")
	}

	flags := syscall.O_NONBLOCK

	for _, modname := range os.Args[1:] {
		modnameptr, err := syscall.BytePtrFromString(modname)
		if err != nil {
			log.Fatalf("rmmod: %v\n", err)
		}
		ret, _, err := syscall.Syscall(syscall.SYS_DELETE_MODULE, uintptr(unsafe.Pointer(modnameptr)), uintptr(flags), 0)
		if ret != 0 {
			log.Fatalf("rmmod: error removing '%s': %v %v\n", modname, ret, err)
		}
	}
}
