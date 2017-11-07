// Copyright 2012-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pci

import "fmt"

type BusReader interface {
	Read() ([]*PCI, error)
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

// String concatenates PCI address, Vendor, and Device to make a useful 
// display for the user.
func (p *PCI) String() string {
	if *numbers {
		return fmt.Sprintf("%s: %s:%s", p.Addr, p.Vendor, p.Device)
	} 
	return fmt.Sprintf("%s: %v %v", p.Addr, p.VendorName, p.DeviceName)
}
