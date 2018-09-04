// Copyright 2018 the LinuxBoot Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package uefi

import (
	"errors"
	"fmt"
)

// GBERegion represents the GBE Region in the firmware.
type GBERegion struct {
	// holds the raw data
	buf []byte
	//Metadata for extraction and recovery
	ExtractPath string
	// This is a pointer to the Region struct laid out in the ifd
	Position *Region
}

// NewGBERegion parses a sequence of bytes and returns a GBERegion
// object, if a valid one is passed, or an error. It also points to the
// Region struct uncovered in the ifd.
func NewGBERegion(buf []byte, r *Region) (*GBERegion, error) {
	gbe := GBERegion{buf: buf, Position: r}
	return &gbe, nil
}

// Buf returns the buffer.
// Used mostly for things interacting with the Firmware interface.
func (gbe *GBERegion) Buf() []byte {
	return gbe.buf
}

// SetBuf sets the buffer.
// Used mostly for things interacting with the Firmware interface.
func (gbe *GBERegion) SetBuf(buf []byte) {
	gbe.buf = buf
}

// Apply calls the visitor on the GBERegion.
func (gbe *GBERegion) Apply(v Visitor) error {
	return v.Visit(gbe)
}

// ApplyChildren calls the visitor on each child node of GBERegion.
func (gbe *GBERegion) ApplyChildren(v Visitor) error {
	return nil
}

// Validate Region
func (gbe *GBERegion) Validate() []error {
	// TODO: Add more verification if needed.
	errs := make([]error, 0)
	if gbe.Position == nil {
		errs = append(errs, errors.New("GBERegion position is nil"))
	}
	if !gbe.Position.Valid() {
		errs = append(errs, fmt.Errorf("GBERegion is not valid, region was %v", *gbe.Position))
	}
	return errs
}
