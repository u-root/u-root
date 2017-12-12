// Copyright 2012-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package pci

import (
	"testing"
)

func TestLookup(t *testing.T) {

	var idLookupTests = []*PCI{
		&PCI{Vendor: "1055", Device: "e420", VendorName: "Efar Microsystems", DeviceName: "LAN9420/LAN9420i"},
		&PCI{Vendor: "8086", Device: "1237", VendorName: "Intel Corporation", DeviceName: "440FX - 82441FX PMC [Natoma]"},
		&PCI{Vendor: "8086", Device: "7000", VendorName: "Intel Corporation", DeviceName: "82371SB PIIX3 ISA [Natoma/Triton II]"},
		&PCI{Vendor: "8086", Device: "7111", VendorName: "Intel Corporation", DeviceName: "82371AB/EB/MB PIIX4 IDE"},
		&PCI{Vendor: "80ee", Device: "beef", VendorName: "InnoTek Systemberatung GmbH", DeviceName: "VirtualBox Graphics Adapter"},
		&PCI{Vendor: "8086", Device: "100e", VendorName: "Intel Corporation", DeviceName: "82540EM Gigabit Ethernet Controller"},
		&PCI{Vendor: "80ee", Device: "cafe", VendorName: "InnoTek Systemberatung GmbH", DeviceName: "VirtualBox Guest Service"},
		&PCI{Vendor: "8086", Device: "2415", VendorName: "Intel Corporation", DeviceName: "82801AA AC'97 Audio Controller"},
		&PCI{Vendor: "8086", Device: "7113", VendorName: "Intel Corporation", DeviceName: "82371AB/EB/MB PIIX4 ACPI"},
		&PCI{Vendor: "8086", Device: "100f", VendorName: "Intel Corporation", DeviceName: "82545EM Gigabit Ethernet Controller (Copper)"},
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
