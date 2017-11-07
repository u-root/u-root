// Copyright 2012-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pci

import (
	"bytes"
	"fmt"
)

type BusReader interface {
	Read() (Devices, error)
}

// PCI is a PCI device. We will fill this in as we add options.
// For now it just holds two uint16 per the PCI spec.
type PCI struct {
	Addr       string
	Vendor     string `pci:"vendor"`
	Device     string `pci:"device"`
	VendorName string
	DeviceName string
}

//Devices is a container for []*PCI and Numbers config option
type Devices struct {
	PCIs    []*PCI
	Numbers bool
}

// String concatenates PCI address, Vendor, and Device to make a useful
// display for the user.
func (d Devices) String() string {
	var buffer bytes.Buffer

	for _, pci := range d.PCIs {
		if d.Numbers {
			buffer.WriteString(fmt.Sprintf("%s: %s:%s\n", pci.Addr, pci.Vendor, pci.Device))
			continue
		}
		pci.VendorName, pci.DeviceName = lookup(pci.Vendor, pci.Device)
		buffer.WriteString(fmt.Sprintf("%s: %s %s\n", pci.Addr, pci.VendorName, pci.DeviceName))
	}

	return buffer.String()
}
