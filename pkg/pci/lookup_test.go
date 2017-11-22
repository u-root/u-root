// Copyright 2012-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package pci

import (
	"testing"
)

func TestLookup(t *testing.T) {

	var idLookupTests = []struct {
		vendor     string
		device     string
		vendorName string
		deviceName string
	}{
		{"1055", "e420", "Efar Microsystems", "LAN9420/LAN9420i"},
		{"8086", "1237", "Intel Corporation", "440FX - 82441FX PMC [Natoma]"},
		{"8086", "7000", "Intel Corporation", "82371SB PIIX3 ISA [Natoma/Triton II]"},
		{"8086", "7111", "Intel Corporation", "82371AB/EB/MB PIIX4 IDE"},
		{"80ee", "beef", "InnoTek Systemberatung GmbH", "VirtualBox Graphics Adapter"},
		{"8086", "100e", "Intel Corporation", "82540EM Gigabit Ethernet Controller"},
		{"80ee", "cafe", "InnoTek Systemberatung GmbH", "VirtualBox Guest Service"},
		{"8086", "2415", "Intel Corporation", "82801AA AC'97 Audio Controller"},
		{"8086", "7113", "Intel Corporation", "82371AB/EB/MB PIIX4 ACPI"},
		{"8086", "100f", "Intel Corporation", "82545EM Gigabit Ethernet Controller (Copper)"},
	}

	t.Run("Lookup Using IDs", func(t *testing.T) {
		ids, err := NewIDs()
		if err != nil {
			t.Fatalf("NewIDs error:%s\n", err)
		}
		for _, tt := range idLookupTests {
			v, d := Lookup(ids, tt.vendor, tt.device)
			if v != tt.vendorName {
				t.Errorf("Vendor mismatch, found %s, expected %s\n", v, tt.vendorName)
			}
			if d != tt.deviceName {
				t.Errorf("Device mismatch, found %s, expected %s\n", d, tt.deviceName)
			}
		}
	})

}
