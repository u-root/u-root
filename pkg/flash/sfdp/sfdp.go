// Copyright 2021 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package sfdp reads the SFDP (Serial Flash Discoverable Parameters) from a
// flash chip where supported. The SFDP is in a separate address space of the
// flash chip and usually read-only. It is used to discover the features
// implemented by the flash chip.
//
// This supports v1.0 of SFDP. Support for v1.5 is incomplete.
//
// Useful references:
// * Your flash chip's datasheet
// * https://chromium.googlesource.com/chromiumos/platform/ec/+/master/include/sfdp.h
// * https://www.macronix.com/Lists/ApplicationNote/Attachments/1870/AN114v1-SFDP%20Introduction.pdf
// * JEDEC specifications
package sfdp

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"strings"
)

const (
	// magic appears as the first 4 bytes of the SFDP table.
	magic = "SFDP"
)

// Param contains information necessary to lookup a parameter in the SFDP
// tables.
type Param struct {
	Table uint16
	Dword int
	Shift int
	Bits  int
}

// Parameters from the Basic Table (id=0).
var (
	ParamBlockSectorEraseSize           = Param{0, 0, 0x00, 0x02}
	ParamWriteGranularity               = Param{0, 0, 0x02, 0x01}
	ParamWriteEnableInstructionRequired = Param{0, 0, 0x03, 0x01}
	ParamWriteEnableOpcodeSelect        = Param{0, 0, 0x04, 0x01}
	Param4KBEraseOpcode                 = Param{0, 0, 0x08, 0x08}
	Param112FastRead                    = Param{0, 0, 0x10, 0x01}
	ParamAddressBytesNumberUsed         = Param{0, 0, 0x11, 0x02}
	ParamDoubleTransferRateClocking     = Param{0, 0, 0x13, 0x01}
	Param122FastReadSupported           = Param{0, 0, 0x14, 0x01}
	Param144FastReadSupported           = Param{0, 0, 0x15, 0x01}
	Param114FastReadSupported           = Param{0, 0, 0x16, 0x01}
	ParamFlashMemoryDensity             = Param{0, 1, 0x00, 0x20}
	Param144FastReadNumberOfWaitStates  = Param{0, 2, 0x00, 0x05}
	Param144FastReadNumberOfModeBits    = Param{0, 2, 0x05, 0x03}
	Param144FastReadOpcode              = Param{0, 2, 0x08, 0x08}
	Param114FastReadNumberOfWaitStates  = Param{0, 2, 0x10, 0x05}
	Param114FastReadNumberOfModeBits    = Param{0, 2, 0x15, 0x03}
	Param114FastReadOpcode              = Param{0, 2, 0x18, 0x08}
	Param112FastReadNumberOfWaitStates  = Param{0, 3, 0x00, 0x05}
	Param112FastReadNumberOfModeBits    = Param{0, 3, 0x05, 0x03}
	Param112FastReadOpcode              = Param{0, 3, 0x08, 0x08}
	Param122FastReadNumberOfWaitStates  = Param{0, 3, 0x10, 0x05}
	Param122FastReadNumberOfModeBits    = Param{0, 3, 0x15, 0x03}
	Param122FastReadOpcode              = Param{0, 3, 0x18, 0x08}
	Param222FastReadSupported           = Param{0, 4, 0x00, 0x01}
	Param444FastReadSupported           = Param{0, 4, 0x04, 0x01}
)

// ParamLookupEntry is a single entry in the BasicTableLookup.
type ParamLookupEntry struct {
	Name  string
	Param Param
}

// BasicTableLookup maps a textual param name to the Param for debug utilities.
// These are in a separate data structure to facilitate dead code elimination
// and reduce size for programs which do not require this information.
var BasicTableLookup = []ParamLookupEntry{
	{"BlockSectorEraseSize", ParamBlockSectorEraseSize},
	{"WriteGranularity", ParamWriteGranularity},
	{"WriteEnableInstructionRequired", ParamWriteEnableInstructionRequired},
	{"WriteEnableOpcodeSelect", ParamWriteEnableOpcodeSelect},
	{"4KBEraseOpcode", Param4KBEraseOpcode},
	{"112FastRead", Param112FastRead},
	{"AddressBytesNumberUsed", ParamAddressBytesNumberUsed},
	{"DoubleTransferRateClocking", ParamDoubleTransferRateClocking},
	{"122FastReadSupported", Param122FastReadSupported},
	{"144FastReadSupported", Param144FastReadSupported},
	{"114FastReadSupported", Param114FastReadSupported},
	{"FlashMemoryDensity", ParamFlashMemoryDensity},
	{"144FastReadNumberOfWaitStates", Param144FastReadNumberOfWaitStates},
	{"144FastReadNumberOfModeBits", Param144FastReadNumberOfModeBits},
	{"144FastReadOpcode", Param144FastReadOpcode},
	{"114FastReadNumberOfWaitStates", Param114FastReadNumberOfWaitStates},
	{"114FastReadNumberOfModeBits", Param114FastReadNumberOfModeBits},
	{"114FastReadOpcode", Param114FastReadOpcode},
	{"112FastReadNumberOfWaitStates", Param112FastReadNumberOfWaitStates},
	{"112FastReadNumberOfModeBits", Param112FastReadNumberOfModeBits},
	{"112FastReadOpcode", Param112FastReadOpcode},
	{"122FastReadNumberOfWaitStates", Param122FastReadNumberOfWaitStates},
	{"122FastReadNumberOfModeBits", Param122FastReadNumberOfModeBits},
	{"122FastReadOpcode", Param122FastReadOpcode},
	{"222FastReadSupported", Param222FastReadSupported},
	{"444FastReadSupported", Param444FastReadSupported},
}

// SFDP (Serial Flash Discoverable Parameters) holds a copy of the tables of the SFDP.
//
// The structure is:
//
//	SFDP
//	 |--> Header
//	 \--> []Parameter
//	         |--> ParameterHeader
//	         \--> Table: A copy of the table's contents.
type SFDP struct {
	Header
	Parameters []Parameter
}

// Header is the header of the SFDP.
type Header struct {
	// Signature is 0x50444653 ("SFDP") if the chip supports SFDP.
	Signature                uint32
	MinorRev                 uint8
	MajorRev                 uint8
	NumberOfParameterHeaders uint8
	_                        uint8
}

// Parameter holds a single table.
type Parameter struct {
	ParameterHeader
	Table []byte
}

// ParameterHeader is the header of Parameter.
type ParameterHeader struct {
	IDLSB    uint8
	MinorRev uint8
	MajorRev uint8
	// Length is the number of dwords.
	Length uint8
	// The top byte of Pointer is the MSB of the Id for revision 1.5.
	Pointer uint32
}

// ID returns the ID for the table. Ths size of the ID depends on the SFDP
// version.
func (p ParameterHeader) ID(sfdpMajorRev, sfdpMinorRev uint8) uint16 {
	id := uint16(p.IDLSB)
	if sfdpMajorRev == 1 && sfdpMinorRev == 5 {
		id |= uint16((p.Pointer >> 16) & 0xff00)
	}
	return id
}

// UnsupportedError is returned if no SFDP can be found.
type UnsupportedError struct{}

// Error implements the error interface.
func (*UnsupportedError) Error() string {
	return "could not find an SFDP"
}

// TableError is returned if the given table id cannot be found.
type TableError struct {
	WantTableID uint16
}

// Error implements the error interface.
func (e *TableError) Error() string {
	return fmt.Sprintf("could not find table %#x", e.WantTableID)
}

// DwordError is returned if the given dword index cannot be found in the given table.
type DwordError struct {
	WantTableID uint16
	WantDword   int
}

// Error implements the error interface.
func (e *DwordError) Error() string {
	return fmt.Sprintf("could not find dword %#x in table %#x", e.WantDword, e.WantTableID)
}

// Read reads the SFDP and all tables from the flash chip.
func Read(r io.ReaderAt) (*SFDP, error) {
	headerBuf := make([]byte, binary.Size(Header{}))
	if _, err := r.ReadAt(headerBuf, 0); err != nil {
		return nil, err
	}
	if string(headerBuf[:4]) != magic {
		return nil, &UnsupportedError{}
	}
	var header Header
	if err := binary.Read(bytes.NewBuffer(headerBuf), binary.LittleEndian, &header); err != nil {
		return nil, err
	}

	parametersBuf := make([]byte, binary.Size(ParameterHeader{})*int(header.NumberOfParameterHeaders+1))
	if _, err := r.ReadAt(parametersBuf, int64(binary.Size(Header{}))); err != nil {
		return nil, err
	}
	sfdp := &SFDP{
		Header:     header,
		Parameters: make([]Parameter, header.NumberOfParameterHeaders+1),
	}
	for i := range sfdp.Parameters {
		p := &sfdp.Parameters[i]
		if err := binary.Read(bytes.NewBuffer(parametersBuf), binary.LittleEndian, &p.ParameterHeader); err != nil {
			return nil, err
		}
		p.Table = make([]byte, int(p.Length)*4)
		if _, err := r.ReadAt(p.Table, int64(p.Pointer)&0x00ffffff); err != nil {
			return nil, err
		}
	}
	return sfdp, nil
}

// Dword reads a dword from the SFDP table with the given table id and dword
// index.
func (s *SFDP) Dword(table uint16, dword int) (uint32, error) {
	for _, p := range s.Parameters {
		if p.ID(s.Header.MajorRev, s.Header.MinorRev) == table {
			byteIdx := dword * 4
			if dword < 0 || byteIdx >= len(p.Table) {
				return 0, &DwordError{WantTableID: table, WantDword: dword}
			}
			return uint32(p.Table[byteIdx]) | (uint32(p.Table[byteIdx+1]) << 8) | (uint32(p.Table[byteIdx+2]) << 16) | (uint32(p.Table[byteIdx+3]) << 24), nil
		}
	}
	return 0, &TableError{WantTableID: table}
}

// Param reads a Param from the SFDP table.
func (s *SFDP) Param(p Param) (int64, error) {
	dword, err := s.Dword(p.Table, p.Dword)
	if err != nil {
		return 0, err
	}
	return (int64(dword) >> p.Shift) & ((1 << int64(p.Bits)) - 1), nil
}

// PrettyPrint prints each parameter from the lookup in a human-readable format.
func (s *SFDP) PrettyPrint(w io.Writer, l []ParamLookupEntry) error {
	// Get the max width of the param name.
	width := 0
	for _, p := range l {
		if len(p.Name) > width {
			width = len(p.Name)
		}
	}

	// Print the parameters.
	padding := strings.Repeat(" ", width)
	for _, p := range l {
		val, err := s.Param(p.Param)
		if err == nil {
			_, err := fmt.Fprintf(w, "%s%s %#x\n", p.Name, padding[len(p.Name):], val)
			if err != nil {
				return err
			}
		} else {
			_, err := fmt.Fprintf(w, "%s%s Error: %v\n", p.Name, padding[len(p.Name):], err)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
