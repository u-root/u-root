// Copyright 2016-2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package smbios

import (
	"errors"
	"fmt"
	"strings"
)

// Much of this is auto-generated. If adding a new type, see README for instructions.

// BaseboardInformation is defined in DSP0134 7.3.
type BaseboardInformation struct {
	Table
	Manufacturer                   string        // 04h
	Product                        string        // 05h
	Version                        string        // 06h
	SerialNumber                   string        // 07h
	AssetTag                       string        // 08h
	BoardFeatures                  BoardFeatures // 09h
	LocationInChassis              string        // 0Ah
	ChassisHandle                  uint16        // 0Bh
	BoardType                      BoardType     // 0Dh
	NumberOfContainedObjectHandles uint8         // 0Eh
	ContainedObjectHandles         []uint16      `smbios:"-"` // 0Fh
}

// NewBaseboardInformation parses a generic Table into BaseboardInformation.
func NewBaseboardInformation(t *Table) (*BaseboardInformation, error) {
	if t.Type != TableTypeBaseboardInformation {
		return nil, fmt.Errorf("invalid table type %d", t.Type)
	}
	if t.Len() < 0xf {
		return nil, errors.New("required fields missing")
	}
	bi := &BaseboardInformation{Table: *t}
	off, err := parseStruct(t, 0 /* off */, false /* complete */, bi)
	if err != nil {
		return nil, err
	}
	if bi.NumberOfContainedObjectHandles > 0 {
		if t.Len() != off+2*int(bi.NumberOfContainedObjectHandles) {
			return nil, errors.New("invalid data length")
		}
		for i := 0; i < int(bi.NumberOfContainedObjectHandles); i++ {
			h, err := t.GetWordAt(off)
			if err != nil {
				return nil, err
			}
			bi.ContainedObjectHandles = append(bi.ContainedObjectHandles, h)
		}
	}
	return bi, nil
}

func (bi *BaseboardInformation) String() string {
	lines := []string{
		bi.Header.String(),
		fmt.Sprintf("Manufacturer: %s", bi.Manufacturer),
		fmt.Sprintf("Product Name: %s", bi.Product),
		fmt.Sprintf("Version: %s", bi.Version),
		fmt.Sprintf("Serial Number: %s", bi.SerialNumber),
		fmt.Sprintf("Asset Tag: %s", bi.AssetTag),
		fmt.Sprintf("Features:\n%s", bi.BoardFeatures),
		fmt.Sprintf("Location In Chassis: %s", bi.LocationInChassis),
		fmt.Sprintf("Chassis Handle: 0x%04X", bi.ChassisHandle),
		fmt.Sprintf("Type: %s", bi.BoardType),
		fmt.Sprintf("Contained Object Handles: %d", bi.NumberOfContainedObjectHandles),
	}
	for _, h := range bi.ContainedObjectHandles {
		lines = append(lines, fmt.Sprintf("0x%04X", h))
	}
	return strings.Join(lines, "\n\t")
}

// BoardFeatures is defined in DSP0134 7.3.1.
type BoardFeatures uint8

// BoardFeatures fields are defined in DSP0134 7.3.1
const (
	BoardFeaturesIsHotSwappable                  BoardFeatures = 1 << 4 // Set to 1 if the board is hot swappable
	BoardFeaturesIsReplaceable                                 = 1 << 3 // Set to 1 if the board is replaceable
	BoardFeaturesIsRemovable                                   = 1 << 2 // Set to 1 if the board is removable
	BoardFeaturesRequiresAtLeastOneDaughterBoard               = 1 << 1 // Set to 1 if the board requires at least one daughter board or auxiliary card to function
	BoardFeaturesIsAHostingBoard                               = 1 << 0 // Set to 1 if the board is a hosting board (for example, a motherboard)
)

func (v BoardFeatures) String() string {
	var lines []string
	if v&BoardFeaturesIsAHostingBoard != 0 {
		lines = append(lines, "Board is a hosting board")
	}
	if v&BoardFeaturesRequiresAtLeastOneDaughterBoard != 0 {
		lines = append(lines, "Board requires at least one daughter board")
	}
	if v&BoardFeaturesIsRemovable != 0 {
		lines = append(lines, "Board is removable")
	}
	if v&BoardFeaturesIsReplaceable != 0 {
		lines = append(lines, "Board is replaceable")
	}
	if v&BoardFeaturesIsHotSwappable != 0 {
		lines = append(lines, "Board is hot swappable")
	}
	return "\t\t" + strings.Join(lines, "\n\t\t")
}

// BoardType is defined in DSP0134 7.3.2.
type BoardType uint8

// BoardType values are defined in DSP0134 7.3.2
const (
	BoardTypeUnknown                                 BoardType = 0x01 // Unknown
	BoardTypeOther                                             = 0x02 // Other
	BoardTypeServerBlade                                       = 0x03 // Server Blade
	BoardTypeConnectivitySwitch                                = 0x04 // Connectivity Switch
	BoardTypeSystemManagementModule                            = 0x05 // System Management Module
	BoardTypeProcessorModule                                   = 0x06 // Processor Module
	BoardTypeIOModule                                          = 0x07 // I/O Module
	BoardTypeMemoryModule                                      = 0x08 // Memory Module
	BoardTypeDaughterBoard                                     = 0x09 // Daughter board
	BoardTypeMotherboardIncludesProcessorMemoryAndIO           = 0x0a // Motherboard (includes processor, memory, and I/O)
	BoardTypeProcessorMemoryModule                             = 0x0b // Processor/Memory Module
	BoardTypeProcessorIOModule                                 = 0x0c // Processor/IO Module
	BoardTypeInterconnectBoard                                 = 0x0d // Interconnect board
)

func (v BoardType) String() string {
	switch v {
	case BoardTypeUnknown:
		return "Unknown"
	case BoardTypeOther:
		return "Other"
	case BoardTypeServerBlade:
		return "Server Blade"
	case BoardTypeConnectivitySwitch:
		return "Connectivity Switch"
	case BoardTypeSystemManagementModule:
		return "System Management Module"
	case BoardTypeProcessorModule:
		return "Processor Module"
	case BoardTypeIOModule:
		return "I/O Module"
	case BoardTypeMemoryModule:
		return "Memory Module"
	case BoardTypeDaughterBoard:
		return "Daughter board"
	case BoardTypeMotherboardIncludesProcessorMemoryAndIO:
		return "Motherboard"
	case BoardTypeProcessorMemoryModule:
		return "Processor/Memory Module"
	case BoardTypeProcessorIOModule:
		return "Processor/IO Module"
	case BoardTypeInterconnectBoard:
		return "Interconnect board"
	}
	return fmt.Sprintf("%#x", uint8(v))
}
