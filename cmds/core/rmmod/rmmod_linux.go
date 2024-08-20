// Copyright 2012-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Remove a module from the Linux kernel
//
// Synopsis:
//
//	rmmod name
//
// Description:
//
//	rmmod is a clone of rmmod(8)
//
// Author:
//
//	Roland Kammerer <dev.rck@gmail.com>
package main

import (
	"log"
	"os"
	"syscall"

	"github.com/u-root/u-root/pkg/kmodule"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatalf("rmmod: ERROR: missing module name.\n")
	}

	for _, modname := range os.Args[1:] {
		if err := kmodule.Delete(modname, syscall.O_NONBLOCK); err != nil {
			log.Fatalf("rmmod: %v", err)
		}
	}
}
