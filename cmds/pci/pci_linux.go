package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"reflect"
	"strconv"
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
		s = s[:len(s)-1]
		i, err := strconv.ParseUint(string(s), 0, 0)
		if err != nil {
			return nil, fmt.Errorf("%v: expected number, got %v: %v", n, string(s), err)
		}
		log.Printf("n is %v, s %v, i %v", n, s, i)
		reflect.ValueOf(&pci).Elem().Field(ix).SetUint(i)
	}
	ve, d := lookup(fmt.Sprintf("%04x", pci.Vendor), fmt.Sprintf("%04x", pci.Device))
	log.Printf("Lookup (%v, %v)", ve, d)

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
