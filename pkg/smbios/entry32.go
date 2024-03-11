// Copyright 2016-2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package smbios

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

// Entry32 is the SMBIOS 32-Bit entry point structure, described in DSP0134 5.2.1.
type Entry32 struct {
	Anchor             [4]uint8
	Checksum           uint8
	Length             uint8
	SMBIOSMajorVersion uint8
	SMBIOSMinorVersion uint8
	StructMaxSize      uint16 // Max size of a single table among all the tables, this definition is different from the one in Entry64.
	Revision           uint8
	Reserved           [5]uint8
	IntAnchor          [5]uint8
	IntChecksum        uint8
	StructTableLength  uint16
	StructTableAddr    uint32
	NumberOfStructs    uint16
	BCDRevision        uint8
}

// UnmarshalBinary unmarshals the SMBIOS 32-Bit entry point structure from binary data.
func (e *Entry32) UnmarshalBinary(data []byte) error {
	if len(data) < 0x1f {
		return fmt.Errorf("invalid entry point stucture length %d", len(data))
	}
	if err := binary.Read(bytes.NewReader(data), binary.LittleEndian, e); err != nil {
		return err
	}
	if !bytes.Equal(e.Anchor[:], []byte("_SM_")) {
		return fmt.Errorf("invalid anchor string %q", string(e.Anchor[:]))
	}
	if int(e.Length) != 0x1f {
		return fmt.Errorf("length mismatch: %d vs %d", e.Length, len(data))
	}
	cs := calcChecksum(data[:e.Length], 4)
	if e.Checksum != cs {
		return fmt.Errorf("checksum mismatch: 0x%02x vs 0x%02x", e.Checksum, cs)
	}
	if !bytes.Equal(e.IntAnchor[:], []byte("_DMI_")) {
		return fmt.Errorf("invalid intermediate anchor string %q", string(e.Anchor[:]))
	}
	intCs := calcChecksum(data[0x10:0x1f], 5)
	if e.IntChecksum != intCs {
		return fmt.Errorf("intermediate checksum mismatch: 0x%02x vs 0x%02x", e.IntChecksum, intCs)
	}
	return nil
}

// MarshalBinary marshals the SMBIOS 32-Bit entry point structure to binary data.
func (e *Entry32) MarshalBinary() ([]byte, error) {
	buf := bytes.NewBuffer(nil)
	if err := binary.Write(buf, binary.LittleEndian, e); err != nil {
		return nil, err
	}
	// Adjust checksums.
	data := buf.Bytes()
	data[0x15] = calcChecksum(data[0x10:0x1f], 5)
	data[4] = calcChecksum(data, 4)
	return data, nil
}
