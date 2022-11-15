// Copyright 2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package dt contains utilities for device tree reading on Linux.
package dt

import (
	"io"
)

const sysfsFDT = "/sys/firmware/fdt"

// LoadFDT loads a flattened device tree from current running system.
//
// It first tries to load it from given io.ReaderAt, then from
// that passed-in file name. If there are not passed-in file names,
// it will try sysfsFDT.
//
// BUGS:
// It is a bit clunky due to its origins; in the original version it
// even had a race. hopefully we can deprecate it in a future u-root
// release.
func LoadFDT(dtb io.ReaderAt, names ...string) (*FDT, error) {
	f := []FDTReader{}
	if dtb != nil {
		f = append(f, WithReaderAt(dtb))
	}
	if len(names) == 0 {
		f = append(f, WithFileName(sysfsFDT))
	} else {
		for _, n := range names {
			f = append(f, WithFileName(n))
		}
	}
	return New(f...)
}
