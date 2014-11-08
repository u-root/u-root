// Copyright 2012 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
       "flag"
	"log"
	"syscall"
)

var (
	RO = flag.Bool("r", false, "Read only mount")
	fsType = flag.String("t", "", "File system type")
)

func main() {
	flags := syscall.MS_MGC_VAL
	flag.Parse()
	a := flag.Args()
	if len(a) < 2 {
		log.Fatalf("Usage: mount [-r] [-t fstype] dev path")
	}
	dev := a[0]
	path := a[1]
	if *RO {
		flags |= syscall.MS_RDONLY
	}
	if err := syscall.Mount(a[0], a[1], *fsType, uintptr(flags), ""); err != nil {
		log.Fatalf("Mount :%s: on :%s: type :%s: flags %x: %v\n", dev, path, fsType, flags, err)
	}
}
