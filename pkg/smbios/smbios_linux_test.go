// Copyright 2016-2021 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package smbios

import (
	"testing"
)

func TestGetSMBIOSEFISMBIOS2(t *testing.T) {
	systabPath = "testdata/smbios2_systab"
	base, size, err := GetSMBIOSBaseEFI()
	if err != nil {
		t.Fatal(err)
	}

	var want int64 = 0x12345678

	if base != want {
		t.Errorf("GetSMBIOSEFI() get 0x%x, want 0x%x", base, want)
	}
	if size != smbios2HeaderSize {
		t.Errorf("GetSMBIOSLegacy() get size 0x%x, want 0x%x ", size, smbios2HeaderSize)
	}
}

func TestGetSMBIOSEFISMBIOS3(t *testing.T) {
	systabPath = "testdata/smbios3_systab"
	base, size, err := GetSMBIOSBaseEFI()
	if err != nil {
		t.Fatal(err)
	}

	var want int64 = 0x12345678

	if base != want {
		t.Errorf("GetSMBIOSEFI() get 0x%x, want 0x%x", base, want)
	}
	if size != smbios3HeaderSize {
		t.Errorf("GetSMBIOSEFI() get size 0x%x, want 0x%x ", size, smbios3HeaderSize)
	}
}

func TestGetSMBIOSEFINotFound(t *testing.T) {
	systabPath = "testdata/systab_NOT_FOUND"
	_, _, err := GetSMBIOSBaseEFI()
	if err == nil {
		t.Fatal("systab should be not found")
	}
}

func TestGetSMBIOSEFIInvalid(t *testing.T) {
	systabPath = "testdata/invalid_systab"
	_, _, err := GetSMBIOSBaseEFI()
	if err == nil {
		t.Fatal("systab should be invalid")
	}
}
