// Copyright 2015-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build !plan9 && !windows

package main

import (
	"syscall"
	"time"
)

func setDate(d string, z *time.Location, clocksource Clock) error {
	t, err := getTime(z, d, clocksource)
	if err != nil {
		return err
	}
	tv := syscall.NsecToTimeval(t.UnixNano())
	return syscall.Settimeofday(&tv)
}
