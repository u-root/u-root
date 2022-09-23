// Copyright 2012-2022 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build linux || darwin

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
		log.Fatalln("rmmod: missing module name.")
	}
	l, err := kmodule.New()
	if err != nil {
		log.Fatalln(err)
	}
	for _, modname := range os.Args[1:] {
		if err := l.Delete(modname, syscall.O_NONBLOCK); err != nil {
			log.Printf("rmmod %q: %v", modname, err)
		}
	}
}
