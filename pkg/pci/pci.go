// Copyright 2012-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pci

import (
	"encoding/binary"
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
	Class      string `pci:"class"`
	VendorName string
	DeviceName string
	Latency    byte
	IRQPin     byte
	IRQLine    string `pci:"irq"`
	Bridge     bool
	FullPath   string
	ExtraInfo  []string
	Config     []byte
	// The rest only gets filled in config space is read.
	// Type 0
	Control  Control
	Status   Status
	Resource string `pci:"resource"`
	BARS     []BAR

	// Type 1

}

// String concatenates PCI address, Vendor, and Device and other information
// to make a useful display for the user.
func (p *PCI) String() string {
	return strings.Join(append([]string{fmt.Sprintf("%s: %v: %v %v", p.Addr, p.Class, p.VendorName, p.DeviceName)}, p.ExtraInfo...), "\n")
}

// SetVendorDeviceName changes VendorName and DeviceName from a name to a number,
// if possible.
func (p *PCI) SetVendorDeviceName() {
	ids = newIDs()
	p.VendorName, p.DeviceName = Lookup(ids, p.Vendor, p.Device)
}

// ReadConfig reads the config space.
func (p *PCI) ReadConfig() error {
	dev := filepath.Join(p.FullPath, "config")
	c, err := ioutil.ReadFile(dev)
	if err != nil {
		return err

	}
	p.Config = c
	p.Control = Control(binary.LittleEndian.Uint16(c[4:6]))
	p.Status = Status(binary.LittleEndian.Uint16(c[6:8]))
	return nil
}

type barreg struct {
	offset int64
	*os.File
}

func (r *barreg) Read(b []byte) (int, error) {
	return r.ReadAt(b, r.offset)
}

func (r *barreg) Write(b []byte) (int, error) {
	return r.WriteAt(b, r.offset)
}

// ReadConfigRegister reads a configuration register of size 8, 16, 32, or 64.
// It will only work on little-endian machines.
func (p *PCI) ReadConfigRegister(offset, size int64) (uint64, error) {
	dev := filepath.Join(p.FullPath, "config")
	f, err := os.Open(dev)
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

// Read implements the BusReader interface for type bus. Iterating over each
// PCI bus device, and applying optional Filters to it.
func (bus *bus) Read(filters ...Filter) (Devices, error) {
	devices := make(Devices, 0, len(bus.Devices))
iter:
	for _, d := range bus.Devices {
		p, err := onePCI(d)
		if err != nil {
			return nil, err
		}
		p.Addr = filepath.Base(d)
		p.FullPath = d
		for _, f := range filters {
			if !f(p) {
				continue iter
			}
		}
		// In the older versions of this package, reading was conditional.
		// There is no harm done, and little performance lost, in just reading it.
		// It's less than a millisecond.
		// In all cases, the first 64 bits are visible, so setting vendor
		// and device names is also no problem. If we can't read any bytes
		// at all, that indicates a problem and it's worth passing that problem
		// up to higher levels.
		if err := p.ReadConfig(); err != nil {
			return nil, err
		}
		p.SetVendorDeviceName()

		c := p.Config
		// Fill in whatever random stuff we can, from the base config.
		p.Latency = c[LatencyTimer]
		if c[HeaderType]&HeaderTypeMask == HeaderTypeBridge {
			p.Bridge = true
		}
		p.IRQPin = c[IRQPin]
		devices = append(devices, p)
	}
	return devices, nil
}
