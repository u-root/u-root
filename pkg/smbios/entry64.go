// Copyright 2016-2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package smbios

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

// Entry64 is the SMBIOS 64-Bit entry point structure, described in DSP0134 5.2.2.
type Entry64 struct {
	Anchor             [5]uint8
	Checksum           uint8
	Length             uint8
	SMBIOSMajorVersion uint8
	SMBIOSMinorVersion uint8
	SMBIOSDocRev       uint8
	Revision           uint8
	Reserved           uint8
	StructMaxSize      uint32 // Max possible size of all the tables combined, this definition is different from the one in Entry32.
	StructTableAddr    uint64
}

// UnmarshalBinary unmarshals the SMBIOS 64-Bit entry point structure from binary data.
func (e *Entry64) UnmarshalBinary(data []byte) error {
	if len(data) < 0x18 {
		return fmt.Errorf("invalid entry point stucture length %d", len(data))
	}
	if err := binary.Read(bytes.NewReader(data), binary.LittleEndian, e); err != nil {
		return err
	}
	if !bytes.Equal(e.Anchor[:], []byte("_SM3_")) {
		return fmt.Errorf("invalid anchor string %q", string(e.Anchor[:]))
	}
	if int(e.Length) != 0x18 {
		return fmt.Errorf("length mismatch: %d vs %d", e.Length, len(data))
	}
	cs := calcChecksum(data[:e.Length], 5)
	if e.Checksum != cs {
		return fmt.Errorf("checksum mismatch: 0x%02x vs 0x%02x", e.Checksum, cs)
	}
	return nil
}

// MarshalBinary marshals the SMBIOS 64-Bit entry point structure to binary data.
func (e *Entry64) MarshalBinary() ([]byte, error) {
	buf := bytes.NewBuffer(nil)
	if err := binary.Write(buf, binary.LittleEndian, e); err != nil {
		return nil, err
	}
	// Adjust checksum.
	data := buf.Bytes()
	data[5] = calcChecksum(data, 5)
	return data, nil
}
