// Copyright 2012-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pci

// Lookup takes PCI hex (as strings) vendor and device ID values and returns human
// readable labels for both the vendor and device. Returns the input ID value if
// if label is not found in the ids map.
func Lookup(ids map[string]Vendor, vendor string, device string) (string, string) {

	if v, f1 := ids[vendor]; f1 {
		if d, f2 := v.Devices[device]; f2 {
			return v.Name, string(d)
		}
		// If entry for device doesn't exist return the hex ID
		return v.Name, device
	}
	// If entry for vendor doesn't exist both hex IDs
	return vendor, device
}
