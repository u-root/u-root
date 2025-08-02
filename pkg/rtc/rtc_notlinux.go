// Copyright 2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build !linux

package rtc

import (
	"errors"
	"time"
)

type syscalls any

type realSyscalls struct{}

// Read implements Read for RTC, returning time.Now()
func (r *RTC) Read() (time.Time, error) {
	return time.Now(), nil
}

// Set returns an error for RTC
func (r *RTC) Set(tu time.Time) error {
	return errors.New("not supported")
}
