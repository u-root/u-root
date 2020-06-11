// Copyright 2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mtd

import (
	"fmt"
	"testing"
)

func TestFindVendor(t *testing.T) {
	var tests = []struct {
		id VendorID
		v  VendorName
		e  error
	}{
		{0xba, "ZETTADEVICE", nil},
		{0x123451234, "", fmt.Errorf("%v: not a known vendor", 0x123451234)},
	}

	for _, tt := range tests {
		v, err := VendorFromID(tt.id)
		if fmt.Sprintf("%v", err) != fmt.Sprintf("%v", tt.e) {
			t.Errorf("%v: got %v want %v", tt.id, err, tt.e)
		}
		if tt.e != nil {
			continue
		}
		if v.Name() != tt.v {
			t.Errorf("%v: got (%q) want (%q)", tt.id, v.Name(), tt.v)
		}

	}
}

func TestFindDevice(t *testing.T) {
	var tests = []struct {
		v  VendorName
		id ChipID
		d  ChipName
		e  error
	}{
		{"WINBOND", 0x32, "W49V002FA", nil},
		// Test a synonym
		{"AMD", 0x0212, "S25FL004A", nil},
		{"ZETTADEVICE", 0xaa66aa44, "", fmt.Errorf("no chip with id 0xaa66aa44 for vendor [\"Zetta\"]")},
	}

	for _, tt := range tests {
		v, err := VendorFromName(tt.v)
		t.Logf("vformname %v", v)
		if err != nil {
			t.Errorf("VendorFromName(%v): got %v, want nil", tt.v, err)
		}
		d, err := v.Chip(tt.id)
		if fmt.Sprintf("%v", err) != fmt.Sprintf("%v", tt.e) {
			t.Errorf("(%q,%v): got %v want %v", tt.v, tt.id, err, tt.e)
		}
		if tt.e != nil {
			continue
		}
		t.Logf("%s", d.Name())
		if d.Name() != tt.d {
			t.Errorf("(%q, %#x): got (%q) want (%q)", tt.v, tt.id, d, tt.d)
		}

	}
}
