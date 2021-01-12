// Copyright 2016-2021 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package smbios

const (
	smbios2HeaderSize = 0x1f
	smbios3HeaderSize = 0x18
)

// GetSMBIOSBase returns SMBIOS Table's base pointer.
func GetSMBIOSBase() (int64, int64, error) {
	base, size, err := GetSMBIOSBaseEFI()
	if err != nil {
		base, size, err = GetSMBIOSBaseLegacy()
		if err != nil {
			return 0, 0, err
		}
	}
	return base, size, nil
}
