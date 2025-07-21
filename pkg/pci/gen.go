// Copyright 2012-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build ignore

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
var IDs =[]Vendor {
`
)

func isHex(b byte) bool {
	return ('a' <= b && b <= 'f') || ('A' <= b && b <= 'F') || ('0' <= b && b <= '9')
}

// scan searches for Vendor and Device lines from the input *bufio.Scanner based
// on pci.ids format. Found Vendors and Devices are added to the input ids map.
func scan(s *bufio.Scanner) []pci.Vendor {
	var ids []pci.Vendor
	var currentVendor uint16
	var line string
	i := -1

	for s.Scan() {
		line = s.Text()

		switch {
		case len(line) > 2 && isHex(line[0]) && isHex(line[1]):
			v, err := strconv.ParseUint(line[:4], 16, 16)
			if err != nil {
				log.Printf("Bad hex constant for vendor: %v", line[:4])
				continue
			}
			currentVendor = uint16(v)
			ids = append(ids, pci.Vendor{ID: uint16(v), Name: string(line[6:]), Devices: []pci.Device{}})
			i++

		case len(line) > 8 && currentVendor != 0 && line[0] == '\t' && isHex(line[1]) && isHex(line[3]):
			v, err := strconv.ParseUint(line[1:5], 16, 16)
			if err != nil {
				log.Printf("Bad hex constant for device: %v", line[1:5])
				continue
			}
			ids[i].Devices = append(ids[i].Devices, pci.Device{ID: uint16(v), Name: string(line[7:])})
		}
	}
	return ids
}

func parse(input []byte) []pci.Vendor {
	s := bufio.NewScanner(bytes.NewReader(input))
	return scan(s)
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
	for _, vendor := range ids {
		code += fmt.Sprintf("Vendor{ID: %#04x, ", vendor.ID)
		code += fmt.Sprintf("Name: %q, Devices: []Device{\n", vendor.Name)
		for _, dev := range vendor.Devices {
			code += fmt.Sprintf("Device{ID:%#04x, Name:%q,},\n", dev.ID, dev.Name)
		}
		code += fmt.Sprintf("},\n},\n")
	}
	code += fmt.Sprintf("}\n")
	err = ioutil.WriteFile("pciids.go", []byte(code), 0o666)
}
