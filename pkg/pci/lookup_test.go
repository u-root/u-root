// Copyright 2012-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pci

import (
	"testing"
)

func TestLookup(t *testing.T) {
	idLookupTests := []*PCI{
		{Vendor: 0x1055, Device: 0xe420, VendorName: "Efar Microsystems", DeviceName: "LAN9420/LAN9420i"},
		{Vendor: 0x8086, Device: 0x1237, VendorName: "Intel Corporation", DeviceName: "440FX - 82441FX PMC [Natoma]"},
		{Vendor: 0x8086, Device: 0x7000, VendorName: "Intel Corporation", DeviceName: "82371SB PIIX3 ISA [Natoma/Triton II]"},
		{Vendor: 0x8086, Device: 0x7111, VendorName: "Intel Corporation", DeviceName: "82371AB/EB/MB PIIX4 IDE"},
		{Vendor: 0x80ee, Device: 0xbeef, VendorName: "InnoTek Systemberatung GmbH", DeviceName: "VirtualBox Graphics Adapter"},
		{Vendor: 0x8086, Device: 0x100e, VendorName: "Intel Corporation", DeviceName: "82540EM Gigabit Ethernet Controller"},
		{Vendor: 0x80ee, Device: 0xcafe, VendorName: "InnoTek Systemberatung GmbH", DeviceName: "VirtualBox Guest Service"},
		{Vendor: 0x8086, Device: 0x2415, VendorName: "Intel Corporation", DeviceName: "82801AA AC'97 Audio Controller"},
		{Vendor: 0x8086, Device: 0x7113, VendorName: "Intel Corporation", DeviceName: "82371AB/EB/MB PIIX4 ACPI"},
		{Vendor: 0x8086, Device: 0x100f, VendorName: "Intel Corporation", DeviceName: "82545EM Gigabit Ethernet Controller (Copper)"},
	}

	t.Run("Lookup Using IDs", func(t *testing.T) {
		for _, tt := range idLookupTests {
			tp := tt
			tp.SetVendorDeviceName()
			v, d := tp.VendorName, tp.DeviceName
			if v != tt.VendorName {
				t.Errorf("Vendor mismatch, found %s, expected %s\n", v, tt.VendorName)
			}
			if d != tt.DeviceName {
				t.Errorf("Device mismatch, found %s, expected %s\n", d, tt.DeviceName)
			}
		}
	})
}
