// Copyright 2016-2021 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package smbios

import (
	"errors"
	"fmt"
)

// SystemSlots is defined in DSP0134 7.10.
type SystemSlots struct {
	Table
	SlotDesignation      string // 04h
	SlotType             uint8  // 05h
	SlotDataBusWidth     uint8  // 06h
	CurrentUsage         uint8  // 07h
	SlotLength           uint8  // 08h
	SlotID               uint16 // 09h
	SlotCharacteristics1 uint8  // 0Bh
	SlotCharacteristics2 uint8  // 0Ch
	SegmentGroupNumber   uint16 // 0Dh
	BusNumber            uint8  // 0Fh
	DeviceFunctionNumber uint8  // 10h
	DataBusWidth         uint8  // 11h
	// TODO: Peer grouping count and Peer groups are omitted for now
}

// ParseSystemSlots parses a generic Table into SystemSlots.
func ParseSystemSlots(t *Table) (*SystemSlots, error) {
	return parseSystemSlots(parseStruct, t)
}

func parseSystemSlots(parseFn parseStructure, t *Table) (*SystemSlots, error) {
	if t.Type != TableTypeSystemSlots {
		return nil, fmt.Errorf("invalid table type %d", t.Type)
	}
	if t.Len() < 0x11 {
		return nil, errors.New("required fields missing")
	}

	ss := &SystemSlots{Table: *t}
	_, err := parseFn(t, 0 /* off */, false /* complete */, ss)
	if err != nil {
		return nil, err
	}
	return ss, nil
}
