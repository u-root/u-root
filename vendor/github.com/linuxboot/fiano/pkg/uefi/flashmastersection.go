// Copyright 2018 the LinuxBoot Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package uefi

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

// FlashMasterSectionSize is the size in bytes of the FlashMaster section
const FlashMasterSectionSize = 12

// RegionPermissions holds the read/write permissions for other regions.
type RegionPermissions struct {
	ID    uint16
	Read  uint8
	Write uint8
}

func (r *RegionPermissions) String() string {
	return fmt.Sprintf("RegionPermissions{ID=%v, Read=0x%x, Write=0x%x}",
		r.ID, r.Read, r.Write)
}

// FlashMasterSection holds all the IDs and read/write permissions for other regions
// This controls whether the bios region can read/write to the ME for example.
type FlashMasterSection struct {
	BIOS RegionPermissions
	ME   RegionPermissions
	GBE  RegionPermissions
}

func (m *FlashMasterSection) String() string {
	return fmt.Sprintf("FlashMasterSection{Bios %v, Me %v, Gbe %v}",
		m.BIOS, m.ME, m.GBE)
}

// NewFlashMasterSection parses a sequence of bytes and returns a FlashMasterSection
// object, if a valid one is passed, or an error
func NewFlashMasterSection(buf []byte) (*FlashMasterSection, error) {
	if len(buf) < FlashMasterSectionSize {
		return nil, fmt.Errorf("Flash Master Section size too small: expected %v bytes, got %v",
			FlashMasterSectionSize,
			len(buf),
		)
	}
	var master FlashMasterSection
	reader := bytes.NewReader(buf)
	if err := binary.Read(reader, binary.LittleEndian, &master); err != nil {
		return nil, err
	}
	return &master, nil
}
