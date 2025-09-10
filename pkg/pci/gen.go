// Copyright 2012-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build ignore

package main

import (
	"bufio"
	"bytes"
	"fmt"
	"go/format"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/u-root/u-root/pkg/pci"
)

var (
	pciidsurl = "https://raw.githubusercontent.com/pciutils/pciids/refs/heads/master/pci.ids"
	header    = `// Copyright 2012-2017 the u-root Authors. All rights reserved
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
		buf bytes.Buffer
	)

	resp, err := http.Get(pciidsurl)
	if err != nil {
		log.Fatalf("unable to download pciids from github:%v", err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	b, err = io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("unable to consume pciids body from github:%v", err)
	}

	ids := parse(b)

	fmt.Fprint(&buf, header)

	for _, vendor := range ids {
		fmt.Fprintf(&buf, "{ID: %#04x, Name: %q, Devices: []Device{\n", vendor.ID, vendor.Name)
		for _, dev := range vendor.Devices {
			fmt.Fprintf(&buf, "{ID:%#04x, Name:%q,},\n", dev.ID, dev.Name)
		}
		fmt.Fprintf(&buf, "},\n},\n")
	}
	fmt.Fprintf(&buf, "}\n")

	p, err := format.Source(buf.Bytes())
	if err != nil {
		log.Fatalf("unable to format source: %v", err)
	}

	err = os.WriteFile("pciids.go", p, 0o666)
}
