// Copyright 2012-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pci

import (
	"testing"
)

func TestLookup(t *testing.T) {
	for _, tt := range []struct {
		name           string
		pci            PCI
		VendorNameWant string
		DeviceNameWant string
	}{
		{
			name: "Lookup Using ID 1055 Device e420",
			pci: PCI{
				Vendor: 0x1055,
				Device: 0xe420,
			},
			VendorNameWant: "Microchip Technology / SMSC",
			DeviceNameWant: "LAN9420/LAN9420i",
		},
		{
			name: "Lookup Using ID 8086 Device 1237",
			pci: PCI{
				Vendor: 0x8086,
				Device: 0x1237,
			},
			VendorNameWant: "Intel Corporation",
			DeviceNameWant: "440FX - 82441FX PMC [Natoma]",
		},
		{
			name: "Lookup Using ID 8086 Device 7000",
			pci: PCI{
				Vendor: 0x8086,
				Device: 0x7000,
			},
			VendorNameWant: "Intel Corporation",
			DeviceNameWant: "82371SB PIIX3 ISA [Natoma/Triton II]",
		},
		{
			name: "Lookup Using ID 8086 Device 7111",
			pci: PCI{
				Vendor: 0x8086,
				Device: 0x7111,
			},
			VendorNameWant: "Intel Corporation",
			DeviceNameWant: "82371AB/EB/MB PIIX4 IDE",
		},
		{
			name: "Lookup Using ID 80ee Device beef",
			pci: PCI{
				Vendor: 0x80ee,
				Device: 0xbeef,
			},
			VendorNameWant: "InnoTek Systemberatung GmbH",
			DeviceNameWant: "VirtualBox Graphics Adapter",
		},
		{
			name: "Lookup Using ID 8086 Device 100e",
			pci: PCI{
				Vendor: 0x8086,
				Device: 0x100e,
			},
			VendorNameWant: "Intel Corporation",
			DeviceNameWant: "82540EM Gigabit Ethernet Controller",
		},
		{
			name: "Lookup Using ID 80ee Device cafe",
			pci: PCI{
				Vendor: 0x80ee,
				Device: 0xcafe,
			},
			VendorNameWant: "InnoTek Systemberatung GmbH",
			DeviceNameWant: "VirtualBox Guest Service",
		},
		{
			name: "Lookup Using ID 8086 Device 2415",
			pci: PCI{
				Vendor: 0x8086,
				Device: 0x2415,
			},
			VendorNameWant: "Intel Corporation",
			DeviceNameWant: "82801AA AC'97 Audio Controller",
		},
		{
			name: "Lookup Using ID 8086 Device 7113",
			pci: PCI{
				Vendor: 0x8086,
				Device: 0x7113,
			},
			VendorNameWant: "Intel Corporation",
			DeviceNameWant: "82371AB/EB/MB PIIX4 ACPI",
		},
		{
			name: "Lookup Using ID 8086 Device 100f",
			pci: PCI{
				Vendor: 0x8086,
				Device: 0x100f,
			},
			VendorNameWant: "Intel Corporation",
			DeviceNameWant: "82545EM Gigabit Ethernet Controller (Copper)",
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			tt.pci.SetVendorDeviceName(IDs)
			VendorNameGot, DeviceNameGot := tt.pci.VendorName, tt.pci.DeviceName
			if VendorNameGot != tt.VendorNameWant {
				t.Errorf("Vendor mismatch, got: '%s', want: '%s'\n", VendorNameGot, tt.VendorNameWant)
			}
			if DeviceNameGot != tt.DeviceNameWant {
				t.Errorf("Device mismatch, got: '%s', want: '%s'\n", DeviceNameGot, tt.DeviceNameWant)
			}
		})
	}
}
