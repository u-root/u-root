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
	"bufio"
	"errors"
	"flag"
	"log"
	"os"
	"strings"

	"golang.org/x/sys/unix"
)

var (
	ro     = flag.Bool("r", false, "Read only mount")
	fsType = flag.String("t", "", "File system type")
	opt    = flag.String("o", "", "Specify mount options")
)

func translateUnknownFS(originFS string, originErr error) error {
	file, err := os.Open("/proc/filesystems")
	if err != nil {
		// just don't make things even worse...
		return originErr
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		tokens := bufio.NewScanner(strings.NewReader(line))
		tokens.Split(bufio.ScanWords)
		var fs string
		for tokens.Scan() {
			fs = tokens.Text()
		}
		// just check the last token of the line
		if fs == originFS {
			return originErr
		}
	}
	return errors.New("Unknown filesystem, check /proc/filesystems for supported ones")
}

func main() {
	// The need for this conversion is not clear to me, but we get an overflow error
	// on ARM without it.
	flags := uintptr(unix.MS_MGC_VAL)
	flag.Parse()
	a := flag.Args()
	if len(a) < 2 {
		log.Fatalf("Usage: mount [-r] [-o mount options] -t fstype dev path")
	}
	dev := a[0]
	path := a[1]
	var data []string
	for _, o := range strings.Split(*opt, ",") {
		f, ok := opts[o]
		if !ok {
			data = append(data, o)
			continue
		}
		flags |= f
	}
	if *ro {
		flags |= unix.MS_RDONLY
	}
	if *fsType == "" {
		log.Printf("No file system type provided!\n")
		log.Fatalf("Usage: mount [-r] [-o mount options] -t fstype dev path")
	} else if err := unix.Mount(a[0], a[1], *fsType, flags, strings.Join(data, ",")); err != nil {
		err = translateUnknownFS(*fsType, err)
		log.Fatalf("Mount :%s: on :%s: type :%s: flags %x: %v\n", dev, path, *fsType, flags, err)
	}
}
