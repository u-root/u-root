// Copyright 2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package build

// From Linux header: /include/uapi/linux/kdev_t.h
const (
	minorBits = 8
	minorMask = (1 << minorBits) - 1
)

// dev returns the device number given the major and minor numbers.
func dev(major, minor uint64) uint64 {
	return major<<minorBits + minor
}

// major returns the device number's major number.
func major(dev uint64) uint64 {
	return dev >> minorBits
}

// minor returns the device number's minor number.
func minor(dev uint64) uint64 {
	return dev & minorMask
}
