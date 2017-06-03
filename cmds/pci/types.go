// Copyright 2012-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import "fmt"

type BusReader interface {
	Read() ([]*PCI, error)
}

// PCI is a PCI device. We will fill this in as we add options.
// For now it just holds two uint16 per the PCI spec.
type PCI struct {
	Addr   string
	Vendor string `pci:"vendor"`
	Device string `pci:"device"`
	VendorName string
	DeviceName string
}

func (p *PCI) String() string {
	return fmt.Sprintf("%s:", p.Addr) +
		fmt.Sprintf(" %v", p.VendorName) +
		fmt.Sprintf(" %v", p.DeviceName)
}
