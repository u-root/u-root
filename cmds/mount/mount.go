// Copyright 2012-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Mount a filesystem at the specified path.
//
// Synopsis:
//     mount [-r] [-o options] [-t FSTYPE] DEV PATH
//
// Options:
//     -r: read only
package main

import (
	"flag"
	"log"

	"golang.org/x/sys/unix"
)

var (
	ro     = flag.Bool("r", false, "Read only mount")
	fsType = flag.String("t", "", "File system type")
	data   = flag.String("o", "", "Specify mount options")
)

func main() {
	// The need for this conversion is not clear to me, but we get an overflow error
	// on ARM without it.
	flags := uintptr(unix.MS_MGC_VAL)
	flag.Parse()
	a := flag.Args()
	if len(a) < 2 {
		log.Fatalf("Usage: mount [-r] [-t fstype] dev path")
	}
	dev := a[0]
	path := a[1]
	if *ro {
		flags |= unix.MS_RDONLY
	}
	if err := unix.Mount(a[0], a[1], *fsType, flags, *data); err != nil {
		log.Fatalf("Mount :%s: on :%s: type :%s: flags %x: %v\n", dev, path, *fsType, flags, err)
	}
}
