// Copyright 2017-2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pci

import (
	"bufio"
	"bytes"
	"log"
	"strconv"
)

func isHex(b byte) bool {
	return ('a' <= b && b <= 'f') || ('A' <= b && b <= 'F') || ('0' <= b && b <= '9')
}

// scan searches for Vendor and Device lines from the input *bufio.Scanner based
// on pci.ids format. Found Vendors and Devices are added to the input ids map.
func scan(s *bufio.Scanner, ids map[uint16]Vendor) {
	var currentVendor uint16
	var line string

	for s.Scan() {
		line = s.Text()

		switch {
		case isHex(line[0]) && isHex(line[1]):
			v, err := strconv.ParseUint(line[:4], 16, 16)
			if err != nil {
				log.Printf("Bad hex constant for vendor: %v", line[:4])
				continue
			}
			currentVendor = uint16(v)
			ids[currentVendor] = Vendor{Name: line[6:], Devices: make(map[uint16]DeviceName)}

		case currentVendor != 0 && line[0] == '\t' && isHex(line[1]) && isHex(line[3]):
			v, err := strconv.ParseUint(line[1:5], 16, 16)
			if err != nil {
				log.Printf("Bad hex constant for device: %v", line[1:5])
				continue
			}
			ids[currentVendor].Devices[uint16(v)] = DeviceName(line[7:])
		}
	}
}

func parse(input []byte) Vendors {
	ids := make(map[uint16]Vendor)
	s := bufio.NewScanner(bytes.NewReader(input))
	scan(s, ids)
	return ids
}
