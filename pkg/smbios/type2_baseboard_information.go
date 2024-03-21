// Copyright 2016-2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package smbios

import (
	"encoding/binary"
	"errors"
	"fmt"
	"strings"
)

// Much of this is auto-generated. If adding a new type, see README for instructions.

// BaseboardInfo is defined in DSP0134 7.3.
type BaseboardInfo struct {
	Header
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

// MarshalBinary encodes the BaseboardInfo content into a binary
func (bi *BaseboardInfo) MarshalBinary() ([]byte, error) {
	t, err := bi.toTable()
	if err != nil {
		return nil, err
	}
	return t.MarshalBinary()
}

func (bi *BaseboardInfo) toTable() (*Table, error) {
	h, err := bi.Header.MarshalBinary()
	if err != nil {
		return nil, err
	}

	var d []byte
	var tableStr []string
	id := byte(1)

	d = append(d, h...)
	if bi.Manufacturer != "" {
		d = append(d, id)
		id++
		tableStr = append(tableStr, bi.Manufacturer)
	} else {
		d = append(d, 0)
	}
	if bi.Product != "" {
		d = append(d, id)
		id++
		tableStr = append(tableStr, bi.Product)
	} else {
		d = append(d, 0)
	}
	if bi.Version != "" {
		d = append(d, id)
		id++
		tableStr = append(tableStr, bi.Version)
	} else {
		d = append(d, 0)
	}
	if bi.SerialNumber != "" {
		d = append(d, id)
		id++
		tableStr = append(tableStr, bi.SerialNumber)
	} else {
		d = append(d, 0)
	}
	if bi.AssetTag != "" {
		d = append(d, id)
		id++
		tableStr = append(tableStr, bi.AssetTag)
	} else {
		d = append(d, 0)
	}

	d = append(d, byte(bi.BoardFeatures))

	if bi.LocationInChassis != "" {
		d = append(d, id)
		tableStr = append(tableStr, bi.LocationInChassis)
	} else {
		d = append(d, 0)
	}

	chb := make([]byte, 2)
	binary.LittleEndian.PutUint16(chb, bi.ChassisHandle)
	d = append(d, chb...)

	d = append(d, byte(bi.BoardType))
	d = append(d, bi.NumberOfContainedObjectHandles)

	if bi.NumberOfContainedObjectHandles != uint8(len(bi.ContainedObjectHandles)) {
		return nil, fmt.Errorf("invalid number of contained object handles, NumberOfContainedObjectHandles: %d, len of ContainedObjectHandles: %d", bi.NumberOfContainedObjectHandles, len(bi.ContainedObjectHandles))
	}

	for _, coh := range bi.ContainedObjectHandles {
		cohb := make([]byte, 2)
		binary.LittleEndian.PutUint16(cohb, coh)
		d = append(d, cohb...)
	}

	t := &Table{
		Header:  bi.Header,
		data:    d,
		strings: tableStr,
	}
	return t, nil
}

// ParseBaseboardInfo parses a generic Table into BaseboardInfo.
func ParseBaseboardInfo(t *Table) (*BaseboardInfo, error) {
	return parseBaseboardInfo(parseStruct, t)
}

func parseBaseboardInfo(parseFn parseStructure, t *Table) (*BaseboardInfo, error) {
	if t.Type != TableTypeBaseboardInfo {
		return nil, fmt.Errorf("invalid table type %d", t.Type)
	}
	// Defined in DSP0134 7.3, length of the structure is at least 08h.
	if t.Len() < 0x8 {
		return nil, errors.New("required fields missing")
	}
	bi := &BaseboardInfo{Header: t.Header}
	off, err := parseFn(t, 0 /* off */, false /* complete */, bi)
	if err != nil {
		return nil, err
	}
	if bi.NumberOfContainedObjectHandles > 0 {
		if t.Len() < off+2*int(bi.NumberOfContainedObjectHandles) {
			return nil, errors.New("invalid data length")
		}
		for i := 0; i < int(bi.NumberOfContainedObjectHandles); i++ {
			h, err := t.GetWordAt(off)
			if err != nil {
				return nil, err
			}
			off += 2
			bi.ContainedObjectHandles = append(bi.ContainedObjectHandles, h)
		}
	}

	// Reassign length to reflect actual object handle amount.
	// Since the parsing process will only parse NumberOfContainedObjectHandles handles and ignore the all trailing "zero value" object handle.
	bi.Length = uint8(off)
	return bi, nil
}

func (bi *BaseboardInfo) String() string {
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
		lines = append(lines, fmt.Sprintf("\t0x%04X", h))
	}
	return strings.Join(lines, "\n\t")
}

// BoardFeatures is defined in DSP0134 7.3.1.
type BoardFeatures uint8

// BoardFeatures fields are defined in DSP0134 7.3.1
const (
	BoardFeaturesIsHotSwappable                  BoardFeatures = 1 << 4 // Set to 1 if the board is hot swappable
	BoardFeaturesIsReplaceable                   BoardFeatures = 1 << 3 // Set to 1 if the board is replaceable
	BoardFeaturesIsRemovable                     BoardFeatures = 1 << 2 // Set to 1 if the board is removable
	BoardFeaturesRequiresAtLeastOneDaughterBoard BoardFeatures = 1 << 1 // Set to 1 if the board requires at least one daughter board or auxiliary card to function
	BoardFeaturesIsAHostingBoard                 BoardFeatures = 1 << 0 // Set to 1 if the board is a hosting board (for example, a motherboard)
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
	BoardTypeOther                                   BoardType = 0x02 // Other
	BoardTypeServerBlade                             BoardType = 0x03 // Server Blade
	BoardTypeConnectivitySwitch                      BoardType = 0x04 // Connectivity Switch
	BoardTypeSystemManagementModule                  BoardType = 0x05 // System Management Module
	BoardTypeProcessorModule                         BoardType = 0x06 // Processor Module
	BoardTypeIOModule                                BoardType = 0x07 // I/O Module
	BoardTypeMemoryModule                            BoardType = 0x08 // Memory Module
	BoardTypeDaughterBoard                           BoardType = 0x09 // Daughter board
	BoardTypeMotherboardIncludesProcessorMemoryAndIO BoardType = 0x0a // Motherboard (includes processor, memory, and I/O)
	BoardTypeProcessorMemoryModule                   BoardType = 0x0b // Processor/Memory Module
	BoardTypeProcessorIOModule                       BoardType = 0x0c // Processor/IO Module
	BoardTypeInterconnectBoard                       BoardType = 0x0d // Interconnect board
)

func (v BoardType) String() string {
	names := map[BoardType]string{
		BoardTypeUnknown:                                 "Unknown",
		BoardTypeOther:                                   "Other",
		BoardTypeServerBlade:                             "Server Blade",
		BoardTypeConnectivitySwitch:                      "Connectivity Switch",
		BoardTypeSystemManagementModule:                  "System Management Module",
		BoardTypeProcessorModule:                         "Processor Module",
		BoardTypeIOModule:                                "I/O Module",
		BoardTypeMemoryModule:                            "Memory Module",
		BoardTypeDaughterBoard:                           "Daughter board",
		BoardTypeMotherboardIncludesProcessorMemoryAndIO: "Motherboard",
		BoardTypeProcessorMemoryModule:                   "Processor/Memory Module",
		BoardTypeProcessorIOModule:                       "Processor/IO Module",
		BoardTypeInterconnectBoard:                       "Interconnect board",
	}
	if name, ok := names[v]; ok {
		return name
	}
	return fmt.Sprintf("%#x", uint8(v))
}
