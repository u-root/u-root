// Copyright 2016-2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package smbios

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

// Header is the header common to all table types.
type Header struct {
	Type   TableType
	Length uint8
	Handle uint16
}

// Parse parses the header from the binary data.
func (h *Header) Parse(data []byte) error {
	return binary.Read(bytes.NewReader(data), binary.LittleEndian, h)
}

// String returns string representation os the header.
func (h *Header) String() string {
	return fmt.Sprintf(
		"Handle 0x%04X, DMI type %d, %d bytes\n%s",
		h.Handle, h.Type, h.Length, h.Type)
}

// MarshalBinary encodes the Header content into a binary
func (h *Header) MarshalBinary() ([]byte, error) {
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.LittleEndian, h)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
