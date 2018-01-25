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

	"github.com/u-root/u-root/pkg/loop"
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
	options mountOptions
)

func init() {
	flag.Var(&options, "o", "Comma separated list of mount options")
}

func loopSetup(filename string) (loopDevice string, err error) {
	loopDevice, err = loop.FindDevice()
	if err != nil {
		return "", err
	}
	if err := loop.SetFdFiles(loopDevice, filename); err != nil {
		return "", err
	}
	return loopDevice, nil
}

func main() {
	flag.Parse()
	a := flag.Args()
	if len(a) < 2 {
		flag.Usage()
		os.Exit(1)
	}

	dev := a[0]
	path := a[1]
	var flags uintptr
	var data []string
	var err error
	for _, option := range options {
		switch option {
		case "loop":
			dev, err = loopSetup(dev)
			if err != nil {
				log.Fatal("Error setting loop device:", err)
			}
		default:
			f, ok := opts[option]
			if !ok {
				data = append(data, option)
				continue
			}
			flags |= f
		}
	}
	if *ro {
		flags |= unix.MS_RDONLY
	}

	if err := unix.Mount(dev, path, *fsType, flags, strings.Join(data, ",")); err != nil {
		log.Fatalf("Mount :%s: on :%s: type :%s: flags %x: %v\n", dev, path, *fsType, flags, err)
	}
}
