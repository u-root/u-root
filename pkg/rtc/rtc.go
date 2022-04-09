// Copyright 2021 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rtc

import (
	"errors"
	"os"
)

type RTC struct {
	file *os.File
	syscalls
}

// OpenRTC opens an RTC. It will typically only work on Linux, but since it
// uses a file API, it will be tried on all systems. Perhaps at some future
// time, other kernels will implement this API.
func OpenRTC() (*RTC, error) {
	devs := []string{
		"/dev/rtc",
		"/dev/rtc0",
		"/dev/misc/rtc0",
	}

	// This logic seems a bit odd, but here is the problem:
	// an error that is NOT IsNotExist may indicate some
	// deeper RTC problem, and should probably halt further efforts.
	// If all opens fail, and we drop out of the loop, then there is
	// no device.
	for _, dev := range devs {
		f, err := os.Open(dev)
		if err == nil {
			return &RTC{f, realSyscalls{}}, err
		}
		if !os.IsNotExist(err) {
			return nil, err
		}
	}

	return nil, errors.New("no RTC device found")
}

// Close closes the RTC
func (r *RTC) Close() error {
	return r.file.Close()
}
