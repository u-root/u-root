// Copyright 2016-2021 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package smbios

import (
	"bytes"
	"fmt"

	"github.com/u-root/u-root/pkg/memio"
)

var memioRead = memio.Read

// getSMBIOSBase searches _SM_ or _SM3_ tag in the given memory range.
func getSMBIOSBase(start, end int64) (int64, int64, error) {
	for base := start; base < end; base++ {
		dat := memio.ByteSlice(make([]byte, 5))
		if err := memioRead(int64(base), &dat); err != nil {
			return 0, 0, err
		}
		if bytes.Equal(dat[:4], []byte("_SM_")) {
			return base, smbios2HeaderSize, nil
		}
		if bytes.Equal(dat[:], []byte("_SM3_")) {
			return base, smbios3HeaderSize, nil
		}
	}
	return 0, 0, fmt.Errorf("could not find _SM_ or _SM3_ via /dev/mem from %#08x to %#08x", start, end)
}

// SMBIOSBaseLegacy searches in SMBIOS entry point address in F0000 segment.
// NOTE: Legacy BIOS will store their SMBIOS in this region.
func SMBIOSBaseLegacy() (int64, int64, error) {
	return getSMBIOSBase(0xf0000, 0x100000)
}
