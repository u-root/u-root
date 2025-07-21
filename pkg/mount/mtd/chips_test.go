// Copyright 2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mtd

import (
	"errors"
	"fmt"
	"strings"
	"testing"
)

func TestFindVendor(t *testing.T) {
	tests := []struct {
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
	tests := []struct {
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

func TestChipFromVIDID(t *testing.T) {
	for _, tt := range []struct {
		vid            VendorID
		vname          VendorName
		cid            ChipID
		wantChipName   ChipName
		wantChipString string
		wantError      error
	}{
		{
			vid:            0xBA,
			vname:          "ZETTADEVICE",
			cid:            0x2012,
			wantChipName:   "ZD25D20",
			wantChipString: "ZETTADEVICE/ZD25D20: 0 pages, 0 pagesize, 0x0 bytes",
		},
		{
			vid:            0xDA,
			vname:          "WINBOND",
			cid:            0x7E1D01,
			wantChipName:   "W29GL032CHL",
			wantChipString: "WINBOND/W29GL032CHL: 0 pages, 0 pagesize, 0x0 bytes, remarks: 	/* Uniform Sectors, WP protects Top OR Bottom sector */",
		},
		{
			vid:       0xAA,
			vname:     "InvalidVendor",
			cid:       0x2012,
			wantError: fmt.Errorf("not a known vendor"),
		},
	} {
		t.Run("Lookup ChipFromVIDDID: "+fmt.Sprintf("VID:CID %d:%d", tt.vid, tt.cid), func(t *testing.T) {
			chip, err := ChipFromVIDDID(tt.vid, tt.cid)
			if !errors.Is(err, tt.wantError) {
				if !strings.Contains(err.Error(), tt.wantError.Error()) {
					t.Errorf("ChipFromVIDDID(tt.vid, tt.cid)=chip, %q, want chip, %q", err, tt.wantError)
				}
			}
			if err != nil {
				return
			}
			if chip.Name() != tt.wantChipName {
				t.Errorf("chip.Name() = %s, want %s", chip.Name(), tt.wantChipName)
			}
			if chip.String() != tt.wantChipString {
				t.Errorf("chip.String()=%s, want %s", chip.String(), tt.wantChipString)
			}
			if chip.ID() != tt.cid {
				t.Errorf("chip.ID()= %x, want %x", chip.ID(), tt.cid)
			}
		})
		t.Run("Lookup VendorFromName: "+string(tt.vname), func(t *testing.T) {
			vendor, err := VendorFromName(tt.vname)
			if !errors.Is(err, tt.wantError) {
				if !strings.Contains(err.Error(), tt.wantError.Error()) {
					t.Errorf("VendorFromName(tt.vname)=vendor, %q, want vendor, %q", err, tt.wantError)
				}
			}
			if err != nil {
				return
			}
			if vendor.Name() != tt.vname {
				t.Errorf("Vendor name does not match up.")
			}
			if vendor.ID() != tt.vid {
				t.Errorf("vendor.ID()= %x, want %x", vendor.ID(), tt.vid)
			}
		})
	}
}
