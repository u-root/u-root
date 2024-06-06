// Copyright 2012-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pci

// Filter can be used to filter a device
type Filter func(p *PCI) bool

// BusReader is the interface for reading device names for a given bus.
type BusReader interface {
	// Read returns Devices, possibly filter by a provided ...Filter
	Read(...Filter) (Devices, error)
}

// Vendor is a PCI vendor, with an ID, a Name, and a possibly empty []Device.
type Vendor struct {
	ID      uint16
	Name    string
	Devices []Device
}

// Device is a PCI Device, with an ID and a Name.
type Device struct {
	ID   uint16
	Name string
}

// Control configures how the device responds to operations. It is the 3rd 16-bit word.
type Control uint16

// Status contains status bits for the PCI device. It is the 4th 16-bit word.
type Status uint16

// BAR is a base address register. It can be a 32- or 64-bit quantity.
// Do you know that PCI was designed by DEC, in the time of Alpha, a 64-bit
// machine, and yet it's still full of 32-bit isms?
// Here's the good news: we don't have to care, since it is
// present as an array of strings in sysfs!
type BAR struct {
	// Index is the index of this resource in the resource list.
	Index int
	// Base is the base, derived (usually) from the resource
	Base uint64
	// Lim is the limit.
	Lim uint64
	// Attr are attributes of this BAR
	Attr uint64
}

// ROM is the expansion ROM type. 32-bit by design.
type ROM uint32

// BridgeCtl is the Bridge Control register.
type BridgeCtl uint16

// BridgeStatus is the Bridge Status register.
type BridgeStatus uint16
