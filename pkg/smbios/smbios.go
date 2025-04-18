// Copyright 2016-2021 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package smbios

const (
	smbios2HeaderSize = 0x1f
	smbios3HeaderSize = 0x18
)

// We need this for testing
type parseStructure func(t *Table, off int, complete bool, sp interface{}) (int, error)

// SMBIOSBase returns SMBIOS Table's base pointer.
func SMBIOSBase() (int64, int64, error) {
	base, size, err := SMBIOSBaseEFI()
	if err != nil {
		base, size, err = SMBIOSBaseLegacy()
		if err != nil {
			return 0, 0, err
		}
	}
	return base, size, nil
}

func SMBIOS3HeaderSize() int64 {
	return smbios3HeaderSize
}
