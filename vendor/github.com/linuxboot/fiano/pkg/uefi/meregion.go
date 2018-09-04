// Copyright 2018 the LinuxBoot Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package uefi

import (
	"errors"
	"fmt"
)

// MERegion represents the ME Region in the firmware.
type MERegion struct {
	// holds the raw data
	buf []byte
	//Metadata for extraction and recovery
	ExtractPath string
	// This is a pointer to the Region struct laid out in the ifd
	Position *Region
}

// NewMERegion parses a sequence of bytes and returns a MERegion
// object, if a valid one is passed, or an error. It also points to the
// Region struct uncovered in the ifd.
func NewMERegion(buf []byte, r *Region) (*MERegion, error) {
	me := MERegion{buf: buf, Position: r}
	return &me, nil
}

// Buf returns the buffer.
// Used mostly for things interacting with the Firmware interface.
func (me *MERegion) Buf() []byte {
	return me.buf
}

// SetBuf sets the buffer.
// Used mostly for things interacting with the Firmware interface.
func (me *MERegion) SetBuf(buf []byte) {
	me.buf = buf
}

// Apply calls the visitor on the MERegion.
func (me *MERegion) Apply(v Visitor) error {
	return v.Visit(me)
}

// ApplyChildren calls the visitor on each child node of MERegion.
func (me *MERegion) ApplyChildren(v Visitor) error {
	return nil
}

// Validate Region
func (me *MERegion) Validate() []error {
	// TODO: Add more verification if needed.
	errs := make([]error, 0)
	if me.Position == nil {
		errs = append(errs, errors.New("MERegion position is nil"))
	}
	if !me.Position.Valid() {
		errs = append(errs, fmt.Errorf("MERegion is not valid, region was %v", *me.Position))
	}
	return errs
}
