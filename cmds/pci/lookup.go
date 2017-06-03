// Copyright 2012-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bufio"
	"bytes"
	"log"
)

const (
	FindingVendor = iota
	FindingDevice
)

var 	debug                                = func(s string, arg ...interface{}) {} //{log.Printf(s, arg...)}

func isHex(b byte) bool {
	return ('a' <= b && b <= 'f') || ('A' <= b && b <= 'F') || ('0' <= b && b <= '9')
}

func lookup(vendor, device string) (string, string) {
	var state = FindingVendor
	var line string
	var lineno int
	vendorName := vendor
	deviceName := device
	defer log.Printf("(%v, %v) -> (%v, %v)", vendor, device, vendorName, deviceName)
	s := bufio.NewScanner(bytes.NewReader(pciids))

	for s.Scan() {
		line = s.Text()
		lineno++
		debug("Line %d: %v\n", lineno, line)
		switch {
		case len(line) == 0, line[0] == '#':
			continue
		case state == FindingVendor && isHex(line[0]) && isHex(line[1]):
			if vendor[:4] == line[:4] {
				vendorName = line[6:]
				state = FindingDevice
			}
		case state == FindingDevice && (line[0] != '\t' || !isHex(line[1])):
			return vendorName, deviceName
		case state == FindingDevice && line[0] == '\t' && isHex(line[1]):
			if device[:4] == line[1:5] {
				deviceName = line[7:]
				return vendorName, deviceName
			}
		}
	}
	return vendorName, deviceName
}
