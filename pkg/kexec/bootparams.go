// Copyright 2015-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package kexec

import (
	"os"
)

// BootParams reads boot params info and returns a []Segment
// for use in a later kexec operation.
// At some point we *may* want to use "/sys/kernel/boot_params/data"
// instead but that is a ton more mess to do.
func BootParams() ([]Segment, error) {
	m, err := os.Open("/dev/mem")
	if err != nil {
		return nil, err
	}

	b := make([]byte, 65536)
	if _, err := m.ReadAt(b, 0x90000); err != nil {
		return nil, err
	}
	segs := []Segment{NewSegment(b, Range{Start: 0x90000, Size: 65536})}
	return segs, nil
}
