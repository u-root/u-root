// Copyright 2016-2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package smbios

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"strings"
)

// Table is a generic type of table that does not parsed fields,
// it only allows access to its contents by offset.
type Table struct {
	Header
	data    []byte   `smbios:"-"` // Structured part of the table, this includes the raw data in Header.
	strings []string `smbios:"-"` // Strings section.
}

var (
	// ErrTableNotFound is retuned if table with specified type is not found.
	ErrTableNotFound = errors.New("table not found")

	// ErrUnsupportedTableType is returned by ParseTypedTable if this table type is not supported and cannot be parsed.
	ErrUnsupportedTableType = errors.New("unsupported table type")

	errEndOfTable = errors.New("end of table")

	tableSep = []byte{0, 0}
)

const (
	outOfSpec = "<OUT OF SPEC>"
)

// Len returns length of the structured part of the table.
func (t *Table) Len() int {
	return len(t.data)
}

// MarshalBinary encodes the table content into a binary
func (t *Table) MarshalBinary() ([]byte, error) {
	var result []byte
	result = append(result, t.data...)
	for _, s := range t.strings {
		result = append(result, []byte(s)...)
		result = append(result, 0x0) // string terminator: 0x0
	}
	if len(t.strings) == 0 { // If there's no string, table still needs a string terminator.
		result = append(result, 0x0) // string terminator: 0x0
	}
	result = append(result, 0x0) // table terminator: 0x0
	return result, nil
}

// GetByteAt returns a byte from the structured part at the specified offset.
func (t *Table) GetByteAt(offset int) (uint8, error) {
	if offset > len(t.data)-1 {
		return 0, fmt.Errorf("invalid offset %d", offset)
	}
	return t.data[offset], nil
}

// GetBytesAt returns a number of bytes from the structured part at the specified offset.
func (t *Table) GetBytesAt(offset, length int) ([]byte, error) {
	if offset > len(t.data)-length {
		return nil, fmt.Errorf("invalid offset %d", offset)
	}
	return t.data[offset : offset+length], nil
}

// GetWordAt returns a 16-bit word from the structured part at the specified offset.
func (t *Table) GetWordAt(offset int) (res uint16, err error) {
	if offset > len(t.data)-2 {
		return 0, fmt.Errorf("invalid offset %d", offset)
	}
	err = binary.Read(bytes.NewReader(t.data[offset:offset+2]), binary.LittleEndian, &res)
	return res, err
}

// GetDWordAt returns a 32-bit word from the structured part at the specified offset.
func (t *Table) GetDWordAt(offset int) (res uint32, err error) {
	if offset > len(t.data)-4 {
		return 0, fmt.Errorf("invalid offset %d", offset)
	}
	err = binary.Read(bytes.NewReader(t.data[offset:offset+4]), binary.LittleEndian, &res)
	return res, err
}

// GetQWordAt returns a 64-bit word from the structured part at the specified offset.
func (t *Table) GetQWordAt(offset int) (res uint64, err error) {
	if offset > len(t.data)-8 {
		return 0, fmt.Errorf("invalid offset %d", offset)
	}
	err = binary.Read(bytes.NewReader(t.data[offset:offset+8]), binary.LittleEndian, &res)
	return res, err
}

// GetStringAt returns a string pointed to by the byte at the specified offset in the structured part.
// NB: offset is not the string index.
func (t *Table) GetStringAt(offset int) (string, error) {
	if offset >= len(t.data) {
		return "", fmt.Errorf("invalid offset %d", offset)
	}
	stringIndex := t.data[offset]
	switch {
	case stringIndex == 0:
		return "Not Specified", nil
	case int(stringIndex) <= len(t.strings):
		return t.strings[stringIndex-1], nil
	default:
		return "<BAD INDEX>", fmt.Errorf("invalid string index %d", stringIndex)
	}
}

func (t *Table) String() string {
	lines := []string{
		t.Header.String(),
		"\tHeader and Data:",
	}
	data := t.data
	for len(data) > 0 {
		ld := data
		if len(ld) > 16 {
			ld = ld[:16]
		}
		ls := make([]string, len(ld))
		for i, d := range ld {
			ls[i] = fmt.Sprintf("%02X", d)
		}
		lines = append(lines, "\t\t"+strings.Join(ls, " "))
		data = data[len(ld):]
	}
	if len(t.strings) > 0 {
		lines = append(lines, "\tStrings:")
		for _, s := range t.strings {
			lines = append(lines, "\t\t"+s)
		}
	}
	return strings.Join(lines, "\n")
}

// ParseTable parses a table from byte stream.
// Returns the parsed table and remaining data.
func ParseTable(data []byte) (*Table, []byte, error) {
	var err error
	var h Header
	if err = h.Parse(data); err != nil {
		return nil, data, err
	}
	if len(data) < int(h.Length)+2 /* string terminator length */ {
		return nil, data, errors.New("data too short")
	}
	structData := data[:h.Length]
	data = data[h.Length:]
	stringData := data
	var strings []string
	for len(data) > 0 && err == nil {
		end := bytes.IndexByte(stringData, 0)
		if end < 0 {
			return nil, data, errors.New("unterminated string")
		}
		s := string(stringData[:end])
		stringData = stringData[end+1:]
		if len(s) > 0 {
			strings = append(strings, s)
		}
		if end == 0 { // End of strings
			break
		}
	}
	// One would think that next table always follows previous table's strings.
	// One would be wrong.
	endOfTableIndex := bytes.Index(data, tableSep)
	if endOfTableIndex < 0 {
		return nil, nil, errors.New("end of table not found")
	}
	data = data[endOfTableIndex+2:]
	if h.Type == TableTypeEndOfTable {
		err = errEndOfTable
	}
	return &Table{Header: h, data: structData, strings: strings}, data, err
}

func kmgt(v uint64) string {
	switch {
	case v >= 1024*1024*1024*1024 && v%(1024*1024*1024*1024) == 0:
		return fmt.Sprintf("%d TB", v/(1024*1024*1024*1024))
	case v >= 1024*1024*1024 && v%(1024*1024*1024) == 0:
		return fmt.Sprintf("%d GB", v/(1024*1024*1024))
	case v >= 1024*1024 && v%(1024*1024) == 0:
		return fmt.Sprintf("%d MB", v/(1024*1024))
	case v >= 1024 && v%1024 == 0:
		return fmt.Sprintf("%d kB", v/1024)
	default:
		return fmt.Sprintf("%d bytes", v)
	}
}
