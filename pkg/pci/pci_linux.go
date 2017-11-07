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

// Read impliments the BusReader interface for type bus.
func (bus *bus) Read() (Devices, error) {
	pci := make([]*PCI, len(bus.Devices))
	for i, d := range bus.Devices {
		p, err := onePCI(d)
		if err != nil {
			return Devices{}, err
		}
		p.Addr = filepath.Base(d)
		pci[i] = p
	}
	return Devices{PCIs: pci}, nil
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
