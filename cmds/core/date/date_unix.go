// Copyright 2015-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build !plan9

package main

import (
	"log"
	"syscall"
	"time"
)

func setDate(d string, z *time.Location) error {
	t, err := getTime(z, d)
	if err != nil {
		log.Fatalf("%v: %v", d, err)
	}
	tv := syscall.NsecToTimeval(t.UnixNano())
	return syscall.Settimeofday(&tv)
}
