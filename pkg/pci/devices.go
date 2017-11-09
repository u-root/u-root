// Copyright 2012-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pci

import "bytes"

//Devices contains a slice of one or more PCI devices
type Devices []*PCI

// ToString concatenates multiple Devices' PCI address, Vendor, and Device
// to make a useful display for the user. Boolean argument toggles displaying
// numeric IDs or human readable labels.
func (d Devices) ToString(n bool) string {
	var buffer bytes.Buffer
	for _, pci := range d {
		buffer.WriteString(pci.ToString(n))
		buffer.WriteString("\n")
	}
	return buffer.String()
}

// String is a Stringer for fmt and others' convenience.
func (d Devices) String() string {
	return d.ToString(false)
}
