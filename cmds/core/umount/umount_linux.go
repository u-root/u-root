// Copyright 2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"errors"
	"flag"

	"github.com/u-root/u-root/pkg/mount"
)

var (
	force = flag.Bool("f", false, "Force unmount")
	lazy  = flag.Bool("l", false, "Lazy unmount")
)

func umount() error {
	flag.Parse()
	a := flag.Args()
	if len(a) != 1 {
		return errors.New("usage: umount [-f | -l] path")
	}
	path := a[0]
	return mount.Unmount(path, *force, *lazy)
}
