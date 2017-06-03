package main

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
		reflect.ValueOf(&pci).Elem().Field(ix).SetString(string(s[2:len(s)-1]))
	}
	pci.VendorName, pci.DeviceName = lookup(pci.Vendor, pci.Device)
	return &pci, nil
}

func (bus *bus) Read() ([]*PCI, error) {
	var pci []*PCI
	for _, d := range bus.Devices {
		p, err := onePCI(d)
		if err != nil {
			return nil, err
		}
		p.Addr = filepath.Base(d)
		pci = append(pci, p)
	}
	return pci, nil
}

// NewBusReader returns a BusReader. If we can't at least glob in
// /sys/bus/pci/devices then we just give up. We don't provide an option
// (yet) to do type I or PCIe MMIO config stuff.
func NewBusReader() (BusReader, error) {
	globs, err := filepath.Glob("/sys/bus/pci/devices/*")
	if err != nil {
		return nil, err
	}

	return &bus{Devices: globs}, nil
}
