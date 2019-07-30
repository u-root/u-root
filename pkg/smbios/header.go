// Copyright 2016-2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package smbios

import (
	"bytes"
	"encoding/binary"
	"errors"
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
	if len(data) < 4 {
		return errors.New("data too short for a header")
	}
	b := bytes.NewBuffer(data)
	binary.Read(b, binary.LittleEndian, &h.Type)
	binary.Read(b, binary.LittleEndian, &h.Length)
	binary.Read(b, binary.LittleEndian, &h.Handle)
	return nil
}

// String returns string representation os the header.
func (h *Header) String() string {
	return fmt.Sprintf(
		"Handle 0x%04x, DMI type %d, %d bytes\n%s",
		h.Handle, h.Type, h.Length, h.Type)
}
