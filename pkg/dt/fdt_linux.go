// Copyright 2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package dt contains utilities for device tree reading on Linux.
package dt

import (
	"io"
	"math"
	"os"
)

var (
	sysfsFDT = "/sys/firmware/fdt"
)

// LoadFDT loads a flattened device tree from current running system.
//
// It first tries to load it from given io.ReaderAt, then from
// /sys/firmware/fdt.
func LoadFDT(dtb io.ReaderAt) (*FDT, error) {
	if dtb != nil {
		fdt, err := ReadFDT(io.NewSectionReader(dtb, 0, math.MaxInt64))
		if err == nil {
			return fdt, nil
		}
	}

	// Fallback to load fdt from sysfs.
	sysFDTFile, err := os.Open(sysfsFDT)
	if err == nil {
		defer sysFDTFile.Close()
		fdt, err := ReadFDT(sysFDTFile)
		if err == nil {
			return fdt, nil
		}
	}

	return nil, errLoadFDT
}
