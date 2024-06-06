// Copyright 2012-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build ignore
// +build ignore

package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"strconv"

	"github.com/u-root/u-root/pkg/pci"
)

var (
	pciidspath = [...]string{"/usr/share/misc/pci.ids"}
	code       = `// Copyright 2012-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pci
var pciids =Vendors {
`
)

func isHex(b byte) bool {
	return ('a' <= b && b <= 'f') || ('A' <= b && b <= 'F') || ('0' <= b && b <= '9')
}

// scan searches for Vendor and Device lines from the input *bufio.Scanner based
// on pci.ids format. Found Vendors and Devices are added to the input ids map.
func scan(s *bufio.Scanner, ids map[uint16]pci.Vendor) {
	var currentVendor uint16
	var line string

	for s.Scan() {
		line = s.Text()

		switch {
		case len(line) > 1 && isHex(line[0]) && isHex(line[1]):
			v, err := strconv.ParseUint(line[:4], 16, 16)
			if err != nil {
				log.Printf("Bad hex constant for vendor: %v", line[:4])
				continue
			}
			currentVendor = uint16(v)
			ids[currentVendor] = pci.Vendor{Name: line[6:], Devices: make(map[uint16]pci.DeviceName)}

		case len(line) > 8 && currentVendor != 0 && line[0] == '\t' && isHex(line[1]) && isHex(line[3]):
			v, err := strconv.ParseUint(line[1:5], 16, 16)
			if err != nil {
				log.Printf("Bad hex constant for device: %v", line[1:5])
				continue
			}
			ids[currentVendor].Devices[uint16(v)] = pci.DeviceName(line[7:])
		}
	}
}

func parse(input []byte) pci.Vendors {
	ids := make(map[uint16]pci.Vendor)
	s := bufio.NewScanner(bytes.NewReader(input))
	scan(s, ids)
	return ids
}

func main() {
	var (
		b   []byte
		err error
	)
	for _, p := range pciidspath {
		b, err = ioutil.ReadFile(p)
		if err == nil {
			break
		}
	}
	if err != nil {
		log.Fatal("can not find a file in %q", pciidspath)
	}
	ids := parse(b)
	for vendor, devices := range ids {
		code += fmt.Sprintf("%#04x: ", vendor)
		code += fmt.Sprintf("Vendor{Name: %q, Devices: map[uint16]DeviceName{\n", devices.Name)
		for j, name := range devices.Devices {
			code += fmt.Sprintf("%#04x:%q,\n", j, name)
		}
		code += fmt.Sprintf("},\n},\n")
	}
	code += fmt.Sprintf("}\n")
	err = ioutil.WriteFile("pciids.go", []byte(code), 0666)
}
