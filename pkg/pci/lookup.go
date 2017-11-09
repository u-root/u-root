// Copyright 2012-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pci

import (
	"bufio"
	"bytes"
)

const (
	findingVendor = iota
	findingDevice
)

var debug = func(s string, arg ...interface{}) {} //{log.Printf(s, arg...)}

func isHex(b byte) bool {
	return ('a' <= b && b <= 'f') || ('A' <= b && b <= 'F') || ('0' <= b && b <= '9')
}

func lookup(vendor, device string) (string, string) {
	var state = findingVendor
	var line string
	var lineno int
	vendorName := vendor
	deviceName := device
	s := bufio.NewScanner(bytes.NewReader(pciids))

	for s.Scan() {
		line = s.Text()
		lineno++
		debug("Line %d: %v\n", lineno, line)
		switch {
		case len(line) < 7, line[0] == '#':
			debug("discard")
			continue
		case state == findingVendor && isHex(line[0]) && isHex(line[1]):
			debug("vendor check %s against %s", vendor[:4], line[:4])
			if vendor[:4] == line[:4] {
				vendorName = line[6:]
				state = findingDevice
			}
		// There are subdevices, ignore them.
		case state == findingDevice && line[0:2] == "\t\t":
			debug("Subdevice")
		case state == findingDevice && (line[0] != '\t' || !isHex(line[1])):
			debug("Finding device: no more devices for this vendor")
			return vendorName, deviceName
		case state == findingDevice && line[0] == '\t' && isHex(line[1]):
			debug("device check %s against %s", device[:4], line[1:5])
			if device[:4] == line[1:5] {
				deviceName = line[7:]
				return vendorName, deviceName
			}
		}
	}
	return vendorName, deviceName
}
