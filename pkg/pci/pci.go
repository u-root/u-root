// Copyright 2012-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pci

import "fmt"

// PCI is a PCI device. We will fill this in as we add options.
// For now it just holds two uint16 per the PCI spec.
type PCI struct {
	Addr       string
	Vendor     string `pci:"vendor"`
	Device     string `pci:"device"`
	VendorName string
	DeviceName string
}

// ToString concatenates PCI address, Vendor, and Device to make a useful
// display for the user. Boolean argument toggles displaying numeric IDs or
// human readable labels.
func (p PCI) ToString(n bool) string {
	if n {
		return fmt.Sprintf("%s: %s:%s", p.Addr, p.Vendor, p.Device)
	}
	p.VendorName, p.DeviceName = lookup(p.Vendor, p.Device)
	return fmt.Sprintf("%s: %v %v", p.Addr, p.VendorName, p.DeviceName)
}

// String is a Stringer for fmt and others' convenience.
func (p PCI) String() string {
	return p.ToString(false)
}
