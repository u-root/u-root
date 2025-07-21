// Copyright 2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mtd

import "slices"

import "fmt"

// VendorFromID returns a Vendor or error given a VendorID.
func VendorFromID(v VendorID) (Vendor, error) {
	for _, vendor := range vendors {
		if vendor.id == v {
			return &vendor, nil
		}
	}
	return nil, fmt.Errorf("%v: not a known vendor", v)
}

// VendorFromName returns a Vendor or error given a VendorName.
func VendorFromName(v VendorName) (Vendor, error) {
	for _, vendor := range vendors {
		if slices.Contains(vendor.names, v) {
			return &vendor, nil
		}
	}
	return nil, fmt.Errorf("%v: not a known vendor", v)
}

// Chip returns a Chip or error given a ChipID.
func (v *vendor) Chip(id ChipID) (Chip, error) {
	for _, d := range devices {
		if d.vendor == v.names[0] && d.id == id {
			return &d, nil
		}
	}
	return nil, fmt.Errorf("no chip with id %#x for vendor %q", id, v.Synonyms())
}

// ID returns a VendorID.
func (v *vendor) ID() VendorID {
	return v.id
}

// Name returns the preferred VendorName.
func (v *vendor) Name() VendorName {
	return v.names[0]
}

// Synonyms returns all the synonyms.
func (v *vendor) Synonyms() []VendorName {
	return v.names[1:]
}

// ChipFromVIDDID will return a Chip struct, given a Vendor and Device ID.
func ChipFromVIDDID(vid VendorID, did ChipID) (Chip, error) {
	v, err := VendorFromID(vid)
	if err != nil {
		return nil, err
	}
	return v.Chip(did)
}

// ID returns the ChipID.
func (c *ChipDevice) ID() ChipID {
	return c.id
}

// Name returns the canonical chip name.
func (c *ChipDevice) Name() ChipName {
	return c.devices[0]
}

// Synonyms returns all synonyms for a chip.
func (c *ChipDevice) Synonyms() []ChipName {
	return c.devices[1:]
}

// Size returns a ChipSize in bytes.
func (c *ChipDevice) Size() ChipSize {
	return ChipSize(c.pageSize * c.numPages)
}

// String is a stringer for a ChipDevice.
func (c *ChipDevice) String() string {
	s := fmt.Sprintf("%s/%s: %d pages, %d pagesize, %#x bytes", c.vendor, c.devices[0], c.numPages, c.pageSize, c.Size())
	if len(c.devices) > 1 {
		s = s + fmt.Sprintf(", synonyms %v", c.devices[1:])
	}
	if len(c.remarks) > 0 {
		s = s + fmt.Sprintf(", remarks: %s", c.remarks)
	}
	return s
}

// Supported returns true if a chip is supported by this package.
func Supported(c Chip) bool {
	return c.Size() != 0
}
