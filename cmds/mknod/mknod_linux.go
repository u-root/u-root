// Copyright 2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"errors"
	"flag"
	"fmt"
	"golang.org/x/sys/unix"
	"strconv"
)

const defaultPerms = 0660

func parseDevices(args []string, devtype string) (int, error) {
	if len(args) != 4 {
		return 0, fmt.Errorf("device type %v requires a major and minor number", devtype)
	}
	major, err := strconv.ParseUint(args[2], 10, 8)
	if err != nil {
		return 0, err
	}
	minor, err := strconv.ParseUint(args[3], 10, 8)
	if err != nil {
		return 0, err
	}
	return int((major << 8) | minor), nil
}

func mknod() error {
	flag.Parse()
	a := flag.Args()
	if len(a) != 2 && len(a) != 4 {
		return errors.New("Usage: mknod path type [major minor]")
	}
	path := a[0]
	devtype := a[1]

	var err error
	var mode uint32

	mode = defaultPerms
	var dev int

	switch devtype {
	case "b":
		// This is a block device. A major/minor number is needed.
		mode |= unix.S_IFBLK
		dev, err = parseDevices(a, devtype)
		if err != nil {
			return err
		}
	case "c", "u":
		// This is a character/unbuffered device. A major/minor number is needed.
		mode |= unix.S_IFCHR
		dev, err = parseDevices(a, devtype)
		if err != nil {
			return err
		}
	case "p":
		// This is a pipe. A major and minor number must not be supplied
		mode |= unix.S_IFIFO
		if len(a) != 2 {
			return fmt.Errorf("device type %v requires no other arguments", devtype)
		}
	default:
		return fmt.Errorf("device type not recognized: %v", devtype)
	}

	if err := unix.Mknod(path, mode, dev); err != nil {
		return fmt.Errorf("mknod :%s: mode %x: %v", path, mode, err)
	}
	return nil
}
