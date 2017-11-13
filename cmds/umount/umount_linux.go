// Copyright 2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"errors"
	"flag"
	"fmt"

	"golang.org/x/sys/unix"
)

var (
	force = flag.Bool("f", false, "Force unmount")
	lazy  = flag.Bool("l", false, "Lazy unmount")
)

func umount() error {
	var flags = unix.UMOUNT_NOFOLLOW
	flag.Parse()
	a := flag.Args()
	if len(a) != 1 {
		return errors.New("Usage: umount [-f | -l] path")
	}
	if *force && *lazy {
		return errors.New("force and lazy unmount cannot both be set")
	}
	path := a[0]
	if *force {
		flags |= unix.MNT_FORCE
	}
	if *lazy {
		flags |= unix.MNT_DETACH
	}
	if err := unix.Unmount(path, flags); err != nil {
		return fmt.Errorf("umount :%s: flags %x: %v", path, flags, err)
	}
	return nil
}
