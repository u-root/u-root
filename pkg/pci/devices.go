// Copyright 2012-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pci

import (
	"bytes"
)

//Devices contains a slice of one or more PCI devices
type Devices []*PCI

// String stringifies the PCI devices. Currently it just calls the device String().
func (d Devices) String() string {
	var buffer bytes.Buffer
	for _, pci := range d {
		buffer.WriteString(pci.String())
		buffer.WriteString("\n")
	}
	return buffer.String()
}

// SetVendorDeviceName sets all numeric IDs of all the devices
// using the pci device SetVendorDeviceName.
func (d Devices) SetVendorDeviceName() {
	for _, p := range d {
		p.SetVendorDeviceName()
	}
}

// ReadConfig reads the config info for all the devices.
func (d Devices) ReadConfig() error {
	for _, p := range d {
		if err := p.ReadConfig(); err != nil {
			return err
		}
	}
	return nil
}

// ReadConfigRegister reads the config info for all the devices.
func (d Devices) ReadConfigRegister(offset, size int64) ([]uint64, error) {
	var vals []uint64
	for _, p := range d {
		val, err := p.ReadConfigRegister(offset, size)
		if err != nil {
			return nil, err
		}
		vals = append(vals, val)
	}
	return vals, nil
}

// WriteConfigRegister writes the config info for all the devices.
func (d Devices) WriteConfigRegister(offset, size int64, val uint64) error {
	for _, p := range d {
		if err := p.WriteConfigRegister(offset, size, val); err != nil {
			return err
		}
	}
	return nil
}
