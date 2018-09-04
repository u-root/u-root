// Copyright 2018 the LinuxBoot Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package uefi

import (
	"errors"
	"fmt"
)

// PDRegion represents the PD Region in the firmware.
type PDRegion struct {
	// holds the raw data
	buf []byte
	//Metadata for extraction and recovery
	ExtractPath string
	// This is a pointer to the Region struct laid out in the ifd
	Position *Region
}

// NewPDRegion parses a sequence of bytes and returns a PDRegion
// object, if a valid one is passed, or an error. It also points to the
// Region struct uncovered in the ifd.
func NewPDRegion(buf []byte, r *Region) (*PDRegion, error) {
	pdr := PDRegion{buf: buf, Position: r}
	return &pdr, nil
}

// Buf returns the buffer.
// Used mostly for things interacting with the Firmware interface.
func (pd *PDRegion) Buf() []byte {
	return pd.buf
}

// SetBuf sets the buffer.
// Used mostly for things interacting with the Firmware interface.
func (pd *PDRegion) SetBuf(buf []byte) {
	pd.buf = buf
}

// Apply calls the visitor on the PDRegion.
func (pd *PDRegion) Apply(v Visitor) error {
	return v.Visit(pd)
}

// ApplyChildren calls the visitor on each child node of PDRegion.
func (pd *PDRegion) ApplyChildren(v Visitor) error {
	return nil
}

// Validate Region
func (pd *PDRegion) Validate() []error {
	// TODO: Add more verification if needed.
	errs := make([]error, 0)
	if pd.Position == nil {
		errs = append(errs, errors.New("PDRegion position is nil"))
	}
	if !pd.Position.Valid() {
		errs = append(errs, fmt.Errorf("PDRegion is not valid, region was %v", *pd.Position))
	}
	return errs
}
