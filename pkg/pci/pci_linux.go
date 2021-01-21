// Copyright 2012-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pci

//go:generate go run gen.go

import (
	"io/ioutil"
	"path/filepath"
	"reflect"
	"sort"
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

// NewBusReader returns a BusReader, given a ...glob to match PCI devices against.
// If it can't glob in pciPath/g then it returns an error.
// For convenience, we use * as the glob if none are supplied.
// We don't provide an option to do type I or PCIe MMIO config stuff.
func NewBusReader(globs ...string) (BusReader, error) {
	if len(globs) == 0 {
		globs = []string{"*"}
	}
	var exp []string
	for _, g := range globs {
		gg, err := filepath.Glob(filepath.Join(pciPath, g))
		if err != nil {
			return nil, err
		}
		exp = append(exp, gg...)
	}
	// uniq
	var u = map[string]struct{}{}
	for _, e := range exp {
		u[e] = struct{}{}
	}
	exp = []string{}
	for v := range u {
		exp = append(exp, v)
	}
	// sort. This might even sort like a shell would do it.
	sort.Strings(exp)
	return &bus{Devices: exp}, nil
}
