// Copyright 2012-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pci

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"os"
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

type barreg struct {
	offset int64
	*os.File
}

func (r *barreg) Read(b []byte) (int, error) {
	return r.ReadAt(b, r.offset)
}

func (w *barreg) Write(b []byte) (int, error) {
	return w.WriteAt(b, w.offset)
}

// ReadConfigRegister reads a configuration register of size 8, 16, 32, or 64.
// It will only work on little-endian machines.
func (p *PCI) ReadConfigRegister(offset, size int64) (uint64, error) {
	f, err := os.Open(filepath.Join(p.FullPath, "config"))
	if err != nil {
		return 0, err
	}
	defer f.Close()
	var reg uint64
	r := &barreg{offset: offset, File: f}
	switch size {
	default:
		return 0, fmt.Errorf("%d is not valid: only options are 8, 16, 32, 64", size)
	case 64:
		err = binary.Read(r, binary.LittleEndian, &reg)
	case 32:
		var val uint32
		err = binary.Read(r, binary.LittleEndian, &val)
		reg = uint64(val)
	case 16:
		var val uint16
		err = binary.Read(r, binary.LittleEndian, &val)
		reg = uint64(val)
	case 8:
		var val uint8
		err = binary.Read(r, binary.LittleEndian, &val)
		reg = uint64(val)
	}
	return reg, err
}

// WriteConfigRegister writes a configuration register of size 8, 16, 32, or 64.
// It will only work on little-endian machines.
func (p *PCI) WriteConfigRegister(offset, size int64, val uint64) error {
	f, err := os.OpenFile(filepath.Join(p.FullPath, "config"), os.O_WRONLY, 0)
	if err != nil {
		return err
	}
	defer f.Close()
	w := &barreg{offset: offset, File: f}
	switch size {
	default:
		return fmt.Errorf("%d is not valid: only options are 8, 16, 32, 64", size)
	case 64:
		err = binary.Write(w, binary.LittleEndian, &val)
	case 32:
		var v = uint32(val)
		err = binary.Write(w, binary.LittleEndian, &v)
	case 16:
		var v = uint16(val)
		err = binary.Write(w, binary.LittleEndian, &v)
	case 8:
		var v = uint8(val)
		err = binary.Write(w, binary.LittleEndian, &v)
	}
	return err
}
