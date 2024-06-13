// Copyright 2012-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pci

//go:generate go run gen.go

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

const (
	pciPath = "/sys/bus/pci/devices"
)

type bus struct {
	Devices []string
}

func readString(dir, file string) (string, error) {
	s, err := os.ReadFile(filepath.Join(dir, file))
	if err != nil {
		return "", err
	}
	return strings.TrimSuffix(string(s), "\n"), nil
}

func readUint(dir, file string, base, bits int) (uint64, error) {
	s, err := readString(dir, file)
	if err != nil {
		return 0, err
	}
	s = strings.TrimPrefix(s, "0x")
	return strconv.ParseUint(s, base, bits)
}

// OnePCI takes the name of a directory containing linux-style
// PCI files and returns a filled-in *PCI.
func OnePCI(dir string) (*PCI, error) {
	pci := PCI{
		Addr:     filepath.Base(dir),
		FullPath: dir,
	}
	var err error
	var n uint64

	if n, err = readUint(dir, "vendor", 16, 16); err != nil {
		return nil, err
	}
	pci.Vendor = uint16(n)
	if n, err = readUint(dir, "device", 16, 16); err != nil {
		return nil, err
	}
	pci.Device = uint16(n)
	if n, err = readUint(dir, "class", 16, 24); err != nil {
		return nil, err
	}
	pci.Class = uint32(n)
	if n, err = readUint(dir, "irq", 0, 0); err != nil {
		return nil, err
	}
	pci.IRQLine = uint(n)
	if pci.Resource, err = readString(dir, "resource"); err != nil {
		return nil, err
	}
	pci.VendorName, pci.DeviceName = fmt.Sprintf("%04x", pci.Vendor), fmt.Sprintf("%04x", pci.Device)
	pci.ClassName = "ClassUnknown"
	if nm, ok := ClassNames[pci.Class]; ok {
		pci.ClassName = nm
	}

	for i, r := range strings.Split(pci.Resource, "\n") {
		b, l, a, err := BaseLimType(r)
		// It's not clear how this can happen, if ever; could someone
		// hotunplug a device while we are scanning?
		if err != nil {
			return nil, fmt.Errorf("scanning resource %d(%s): %w", i, dir, err)
		}
		if b == 0 {
			continue
		}
		nb := BAR{
			Index: i,
			Base:  b,
			Lim:   l,
			Attr:  a,
		}
		switch i {
		case 13:
			pci.IO = nb
		case 14:
			pci.Mem = nb
		case 15:
			pci.PrefMem = nb
		default:
			pci.BARS = append(pci.BARS, nb)
		}
	}
	return &pci, nil
}

// BaseLimType parses a Linux resource string into base, limit, and attributes.
// The string must have three hex fields.
// Gaul was divided into three parts.
// So are the BARs.
func BaseLimType(bar string) (uint64, uint64, uint64, error) {
	f := strings.Fields(bar)
	if len(f) != 3 {
		return 0, 0, 0, fmt.Errorf("bar %q should have 3 fields", bar)
	}
	// They must all be parseable hex numbers.
	var vals [3]uint64
	for i, ff := range f {
		var err error
		if vals[i], err = strconv.ParseUint(ff, 0, 0); err != nil {
			return 0, 0, 0, err
		}
	}
	return vals[0], vals[1], vals[2], nil
}

// NewBusReader returns a BusReader, given a ...glob to match PCI devices against.
// If it can't glob in pciPath/g then it returns an error.
// For convenience, we use * as the glob if none are supplied.
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
	u := map[string]struct{}{}
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
