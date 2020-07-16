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

// Vendor is a PCI vendor human readable label. It contains a map of one or
// more Devices keyed by hex ID.
type Vendor struct {
	Name    string
	Devices map[string]Device
}

// Device is a PCI device human readable label
type Device string
