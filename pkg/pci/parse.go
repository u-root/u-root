package pci

import (
	"bufio"
	"bytes"
)

func isHex(b byte) bool {
	return ('a' <= b && b <= 'f') || ('A' <= b && b <= 'F') || ('0' <= b && b <= '9')
}

// scan searches for Vendor and Device lines from the input *bufio.Scanner based
// on pci.ids format. Found Vendors and Devices are added to the input map.
// This implimentation expects an input pci.ids to have comments, blank lines,
// sub-devices, and classes already removed.
func scan(s *bufio.Scanner, ids map[string]Vendor) error {
	var currentVendor string
	var line string

	for s.Scan() {
		line = s.Text()

		switch {
		case isHex(line[0]) && isHex(line[1]):
			currentVendor = line[:4]
			ids[currentVendor] = Vendor{Name: line[6:], Devices: make(map[string]Device)}
		case currentVendor != "" && line[0] == '\t' && isHex(line[1]):
			ids[currentVendor].Devices[line[1:5]] = Device(line[7:])
		}
	}
	return nil
}

func parse(input []byte) (map[string]Vendor, error) {
	ids := make(map[string]Vendor)

	s := bufio.NewScanner(bytes.NewReader(input))
	err := scan(s, ids)

	return ids, err
}
