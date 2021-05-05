// Copyright 2012-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pci

import (
	"fmt"
	"io"
)

//Devices contains a slice of one or more PCI devices
type Devices []*PCI

// Print prints information to an io.Writer
func (d Devices) Print(o io.Writer, verbose int) error {
	for _, pci := range d {
		if _, err := fmt.Fprintf(o, "%s\n", pci.String()); err != nil {
			return err
		}
		if verbose >= 1 {
			if _, err := fmt.Fprintf(o, "\tControl: %s\n\tStatus: %s\n", pci.Control.String(), pci.Status.String()); err != nil {
				return err
			}
		}

		if verbose > 0 {
			fmt.Fprintf(o, "\n")
		}
	}
	return nil
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
