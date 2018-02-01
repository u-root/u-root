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
	"os"
	"strings"

	"github.com/u-root/u-root/pkg/losetup"
	"golang.org/x/sys/unix"
)

type mountOptions []string

func (o *mountOptions) String() string {
	return strings.Join(*o, ",")
}

func (o *mountOptions) Set(value string) error {
	for _, option := range strings.Split(value, ",") {
		*o = append(*o, option)
	}
	return nil
}

var (
	ro      = flag.Bool("r", false, "Read only mount")
	fsType  = flag.String("t", "", "File system type")
	data    = flag.String("d", "", "Specify fs options")
	options mountOptions
)

func init() {
	flag.Var(&options, "o", "Comma separated list of mount options")
}

func loop(filename string) (loopDevice string, err error) {
	loopDevice, err = losetup.FindLoopDevice()
	if err != nil {
		return "", err
	}
	if err := losetup.LoopSetFdFiles(loopDevice, filename); err != nil {
		return "", err
	}
	return loopDevice, nil
}

func main() {
	// The need for this conversion is not clear to me, but we get an overflow error
	// on ARM without it.

	var err error

	flags := uintptr(unix.MS_MGC_VAL)
	flag.Parse()
	a := flag.Args()
	if len(a) < 2 {
		flag.Usage()
		os.Exit(1)
	}

	dev := a[0]
	path := a[1]
	if *ro {
		flags |= unix.MS_RDONLY
	}

	for _, option := range options {
		switch option {
		case "loop":
			 dev, err = loop(dev)
			if err != nil {
				log.Fatal("Error setting loop device: ", err)
			}
		default:
			log.Fatal("Unrecognized option: ", option)
		}
	}

	if err := unix.Mount(dev, path, *fsType, flags, *data); err != nil {
		log.Fatalf("Mount :%s: on :%s: type :%s: flags %x: %v\n", dev, path, *fsType, flags, err)
	}
}
