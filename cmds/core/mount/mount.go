// Copyright 2012-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build !plan9

// mount mounts a filesystem at the specified path.
//
// Synopsis:
//     mount [-r] [-o options] [-t FSTYPE] DEV PATH
//
// Options:
//     -r: read only
package main

import (
	"flag"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/u-root/u-root/pkg/mount"
	"github.com/u-root/u-root/pkg/mount/loop"
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
	if err := loop.SetFile(loopDevice, filename); err != nil {
		return "", err
	}
	return loopDevice, nil
}

// extended from boot.go
func getSupportedFilesystem(originFS string) ([]string, bool, error) {
	var known bool
	var err error
	fs, err := ioutil.ReadFile("/proc/filesystems")
	if err != nil {
		return nil, known, err
	}
	var returnValue []string
	for _, f := range strings.Split(string(fs), "\n") {
		n := strings.Fields(f)
		last := len(n) - 1
		if last < 0 {
			continue
		}
		if n[last] == originFS {
			known = true
		}
		returnValue = append(returnValue, n[last])
	}
	return returnValue, known, err

}

func informIfUnknownFS(originFS string) {
	knownFS, known, err := getSupportedFilesystem(originFS)
	if err != nil {
		// just don't make things even worse...
		return
	}
	if !known {
		log.Printf("Hint: unknown filesystem %s. Known are: %v", originFS, knownFS)
	}
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
			if f, ok := opts[option]; ok {
				flags |= f
			} else {
				data = append(data, option)
			}
		}
	}
	if *ro {
		flags |= unix.MS_RDONLY
	}
	if *fsType == "" {
		if _, err := mount.TryMount(dev, path, strings.Join(data, ","), flags); err != nil {
			log.Fatalf("%v", err)
		}
	} else {
		if _, err := mount.Mount(dev, path, *fsType, strings.Join(data, ","), flags); err != nil {
			log.Printf("%v", err)
			informIfUnknownFS(*fsType)
			os.Exit(1)
		}
	}
}
