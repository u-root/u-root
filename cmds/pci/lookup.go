// Copyright 2012-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bufio"
	"bytes"

	"github.com/u-root/u-root/pkg/log"
)

const (
	FindingVendor = iota
	FindingDevice
)

func isHex(b byte) bool {
	return ('a' <= b && b <= 'f') || ('A' <= b && b <= 'F') || ('0' <= b && b <= '9')
}

func lookup(vendor, device string) (string, string) {
	var state = FindingVendor
	var line string
	var lineno int
	vendorName := vendor
	deviceName := device

	s := bufio.NewScanner(bytes.NewReader(pciids))
	for s.Scan() {
		line = s.Text()
		lineno++

		log.Printf("Line %d: %v", lineno, line)

		switch {
		case len(line) < 7, line[0] == '#':
			log.Printf("discard")
			continue

		case state == FindingVendor && isHex(line[0]) && isHex(line[1]):
			log.Printf("vendor check %s against %s", vendor[:4], line[:4])
			if vendor[:4] == line[:4] {
				vendorName = line[6:]
				state = FindingDevice
			}

		// There are subdevices, ignore them.
		case state == FindingDevice && line[0:2] == "\t\t":
			log.Printf("Subdevice")

		case state == FindingDevice && (line[0] != '\t' || !isHex(line[1])):
			log.Printf("Finding device: no more devices for this vendor")
			return vendorName, deviceName

		case state == FindingDevice && line[0] == '\t' && isHex(line[1]):
			log.Printf("device check %s against %s", device[:4], line[1:5])
			if device[:4] == line[1:5] {
				deviceName = line[7:]
				return vendorName, deviceName
			}
		}
	}
	return vendorName, deviceName
}
