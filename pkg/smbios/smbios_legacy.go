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

func isValidChecksum(data []byte) bool {
	var sum uint8
	for _, b := range data {
		sum += b
	}
	return sum == 0
}

// getSMBIOSBase searches _SM_ or _SM3_ tag in the given memory range and
// validates the checksum to prevent invalid/junk SMBIOS memory from being used.
func getSMBIOSBase(start, end int64) (int64, int64, error) {
	for base := start; base < end; base++ {
		dat := memio.ByteSlice(make([]byte, 5))
		if err := memioRead(int64(base), &dat); err != nil {
			return 0, 0, err
		}
		if bytes.Equal(dat[:4], []byte("_SM_")) {
			// Read the full 32-bit entry point structure (31 bytes)
			entryData := memio.ByteSlice(make([]byte, smbios2HeaderSize))
			if err := memioRead(int64(base), &entryData); err != nil {
				continue
			}
			if isValidChecksum(entryData) {
				return base, smbios2HeaderSize, nil
			}
		}
		if bytes.Equal(dat[:], []byte("_SM3_")) {
			// Read the full 64-bit entry point structure (24 bytes)
			entryData := memio.ByteSlice(make([]byte, smbios3HeaderSize))
			if err := memioRead(int64(base), &entryData); err != nil {
				continue
			}
			if isValidChecksum(entryData) {
				return base, smbios3HeaderSize, nil
			}
		}
	}
	return 0, 0, fmt.Errorf("could not find valid _SM_ or _SM3_ via /dev/mem from %#08x to %#08x", start, end)
}

// SMBIOSBaseLegacy searches in SMBIOS entry point address in F0000 segment.
// NOTE: Legacy BIOS will store their SMBIOS in this region.
func SMBIOSBaseLegacy() (int64, int64, error) {
	return getSMBIOSBase(0xf0000, 0x100000)
}
