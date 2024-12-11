// Copyright 2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build !plan9 && !windows

// Unmount a filesystem at the specified path.
//
// Synopsis:
//
//	umount [-f | -l] PATH
//
// Options:
//
//	-f: force unmount
//	-l: lazy unmount
package main

import "log"

func main() {
	if err := umount(); err != nil {
		log.Fatal(err)
	}
}
