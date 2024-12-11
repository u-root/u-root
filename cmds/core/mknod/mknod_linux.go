// Copyright 2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build !windows || !plan9

package main

import (
	"errors"
	"fmt"
	"strconv"

	"golang.org/x/sys/unix"
)

const defaultPerms = 0o660

func parseDevices(args []string, devtype string) (int, error) {
	if len(args) != 4 {
		return 0, fmt.Errorf("device type %v requires a major and minor number", devtype)
	}
	major, err := strconv.ParseUint(args[2], 10, 12)
	if err != nil {
		return 0, err
	}
	minor, err := strconv.ParseUint(args[3], 10, 20)
	if err != nil {
		return 0, err
	}
	return int(unix.Mkdev(uint32(major), uint32(minor))), nil
}

func mknod(a []string) error {
	if len(a) != 2 && len(a) != 4 {
		return errors.New("usage: mknod path type [major minor]")
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
		return fmt.Errorf("%q: mode %x: %w", path, mode, err)
	}
	return nil
}
