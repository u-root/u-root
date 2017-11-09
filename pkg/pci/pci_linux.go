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
	return &pci, nil
}

// Read impliments the BusReader interface for type bus. Iterating over each
// PCI bus device.
func (bus *bus) Read() (Devices, error) {
	devices := make(Devices, len(bus.Devices))
	for i, d := range bus.Devices {
		p, err := onePCI(d)
		if err != nil {
			return nil, err
		}
		p.Addr = filepath.Base(d)
		devices[i] = p
	}
	return devices, nil
}

// NewBusReader returns a BusReader. If we can't at least glob in
// /sys/bus/pci/devices then we just give up. We don't provide an option
// (yet) to do type I or PCIe MMIO config stuff.
func NewBusReader() (busReader, error) {
	globs, err := filepath.Glob("/sys/bus/pci/devices/*")
	if err != nil {
		return nil, err
	}

	return &bus{Devices: globs}, nil
}
