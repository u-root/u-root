// Copyright 2012-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pci

import (
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"
)

// PCI is a PCI device. We will fill this in as we add options.
// For now it just holds two uint16 per the PCI spec.
type PCI struct {
	Addr       string
	Vendor     string `pci:"vendor"`
	Device     string `pci:"device"`
	VendorName string
	DeviceName string
	FullPath   string
	ExtraInfo  []string
}

// String concatenates PCI address, Vendor, and Device and other information
// to make a useful display for the user.
func (p *PCI) String() string {
	return strings.Join(append([]string{fmt.Sprintf("%s: %v %v", p.Addr, p.VendorName, p.DeviceName)}, p.ExtraInfo...), "\n")
}

// SetVendorDeviceName changes VendorName and DeviceName from a name to a number,
// if possible.
func (p *PCI) SetVendorDeviceName() {
	ids = newIDs()
	p.VendorName, p.DeviceName = Lookup(ids, p.Vendor, p.Device)
}

// ReadConfig reads the config space and adds it to ExtraInfo as a hexdump.
func (p *PCI) ReadConfig() error {
	c, err := ioutil.ReadFile(filepath.Join(p.FullPath, "config"))
	if err != nil {
		return err
	}
	p.ExtraInfo = append(p.ExtraInfo, hex.Dump(c))
	return nil
}
