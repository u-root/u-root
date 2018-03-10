// Copyright 2012-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pci

//go:generate go run gen.go

import (
	"io/ioutil"
	"path/filepath"
	"reflect"
)

const (
	pciPath = "/sys/bus/pci/devices"
)

type bus struct {
	Devices []string
}

func onePCI(dir string) (*PCI, error) {
	var pci PCI
	v := reflect.TypeOf(pci)
	for ix := 0; ix < v.NumField(); ix++ {
		f := v.Field(ix)
		n := f.Tag.Get("pci")
		if n == "" {
			continue
		}
		s, err := ioutil.ReadFile(filepath.Join(dir, n))
		if err != nil {
			return nil, err
		}
		// Linux never understood /proc.
		reflect.ValueOf(&pci).Elem().Field(ix).SetString(string(s[2 : len(s)-1]))
	}
	pci.VendorName, pci.DeviceName = pci.Vendor, pci.Device
	return &pci, nil
}

// Read implements the BusReader interface for type bus. Iterating over each
// PCI bus device.
func (bus *bus) Read() (Devices, error) {
	devices := make(Devices, len(bus.Devices))
	for i, d := range bus.Devices {
		p, err := onePCI(d)
		if err != nil {
			return nil, err
		}
		p.Addr = filepath.Base(d)
		p.FullPath = d
		devices[i] = p
	}
	return devices, nil
}

// NewBusReader returns a BusReader, given a glob to match PCI devices against.
// If it can't glob in pciPath/g then it returns an error.
// We don't provide an option to do type I or PCIe MMIO config stuff.
func NewBusReader(g string) (busReader, error) {
	globs, err := filepath.Glob(filepath.Join(pciPath, g))
	if err != nil {
		return nil, err
	}

	return &bus{Devices: globs}, nil
}
