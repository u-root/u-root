// Copyright 2016-2021 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package smbios

import "runtime"

const (
	smbios2HeaderSize = 0x1f
	smbios3HeaderSize = 0x18
)

// legacySMBIOSSupported reports whether this architecture stores the SMBIOS
// entry point in IO region. That scan is an x86 convention.
// Other architectures do not map the region and reading it leads to a fault.
func legacySMBIOSSupported() bool {
	switch runtime.GOARCH {
	case "386", "amd64":
		return true
	default:
		return false
	}
}

// We need this for testing
type parseStructure func(t *Table, off int, complete bool, sp any) (int, error)

// SMBIOSBase returns SMBIOS Table's base pointer.
func SMBIOSBase() (int64, int64, error) {
	base, size, err := SMBIOSBaseEFI()
	if err == nil {
		return base, size, nil
	}
	// Firmware without an EFI SMBIOS pointer has nothing further to search
	// unless the platform implements the legacy BIOS entry-point region.
	if !legacySMBIOSSupported() {
		return 0, 0, err
	}
	base, size, err = SMBIOSBaseLegacy()
	if err != nil {
		return 0, 0, err
	}
	return base, size, nil
}

func SMBIOS3HeaderSize() int64 {
	return smbios3HeaderSize
}
