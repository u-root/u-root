// Copyright 2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

// From FreeBSD header: /usr/include/sys/types.h

func major(dev uint64) uint64 {
	return (dev >> 8) & 0xff
}

func minor(dev uint64) uint64 {
	return dev & 0xffff00ff
}
