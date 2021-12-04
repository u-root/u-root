// Copyright 2012-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pci

import "fmt"

const venDevFmt = "%04x"

// Lookup takes PCI and device ID values and returns human
// readable labels for both the vendor and device. It returns the input ID value if
// if label is not found in the ids map.
func Lookup(ids map[uint16]Vendor, vendor uint16, device uint16) (string, string) {
	if v, ok := ids[vendor]; ok {
		if d, ok := v.Devices[device]; ok {
			return v.Name, string(d)
		}
		// If entry for device doesn't exist return the hex ID
		return v.Name, fmt.Sprintf(venDevFmt, device)
	}
	// If entry for vendor doesn't exist both hex IDs
	return fmt.Sprintf(venDevFmt, vendor), fmt.Sprintf(venDevFmt, device)
}
