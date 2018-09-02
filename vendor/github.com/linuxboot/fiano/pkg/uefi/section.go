// Copyright 2018 the LinuxBoot Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package uefi

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"unsafe"

	"github.com/linuxboot/fiano/pkg/lzma"
	"github.com/linuxboot/fiano/pkg/unicode"
	"github.com/linuxboot/fiano/pkg/uuid"
)

const (
	// SectionMinLength is the minimum length of a file section header.
	SectionMinLength = 0x04
	// SectionExtMinLength is the minimum length of an extended file section header.
	SectionExtMinLength = 0x08
)

// SectionType holds a section type value
type SectionType uint8

// UEFI Section types
const (
	SectionTypeAll                 SectionType = 0x00
	SectionTypeCompression         SectionType = 0x01
	SectionTypeGUIDDefined         SectionType = 0x02
	SectionTypeDisposable          SectionType = 0x03
	SectionTypePE32                SectionType = 0x10
	SectionTypePIC                 SectionType = 0x11
	SectionTypeTE                  SectionType = 0x12
	SectionTypeDXEDepEx            SectionType = 0x13
	SectionTypeVersion             SectionType = 0x14
	SectionTypeUserInterface       SectionType = 0x15
	SectionTypeCompatibility16     SectionType = 0x16
	SectionTypeFirmwareVolumeImage SectionType = 0x17
	SectionTypeFreeformSubtypeGUID SectionType = 0x18
	SectionTypeRaw                 SectionType = 0x19
	SectionTypePEIDepEx            SectionType = 0x1b
	SectionMMDepEx                 SectionType = 0x1c
)

var sectionTypeNames = map[SectionType]string{
	SectionTypeCompression:         "EFI_SECTION_COMPRESSION",
	SectionTypeGUIDDefined:         "EFI_SECTION_GUID_DEFINED",
	SectionTypeDisposable:          "EFI_SECTION_DISPOSABLE",
	SectionTypePE32:                "EFI_SECTION_PE32",
	SectionTypePIC:                 "EFI_SECTION_PIC",
	SectionTypeTE:                  "EFI_SECTION_TE",
	SectionTypeDXEDepEx:            "EFI_SECTION_DXE_DEPEX",
	SectionTypeVersion:             "EFI_SECTION_VERSION",
	SectionTypeUserInterface:       "EFI_SECTION_USER_INTERFACE",
	SectionTypeCompatibility16:     "EFI_SECTION_COMPATIBILITY16",
	SectionTypeFirmwareVolumeImage: "EFI_SECTION_FIRMWARE_VOLUME_IMAGE",
	SectionTypeFreeformSubtypeGUID: "EFI_SECTION_FREEFORM_SUBTYPE_GUID",
	SectionTypeRaw:                 "EFI_SECTION_RAW",
	SectionTypePEIDepEx:            "EFI_SECTION_PEI_DEPEX",
	SectionMMDepEx:                 "EFI_SECTION_MM_DEPEX",
}

// String creates a string representation for the file type.
func (s SectionType) String() string {
	if t, ok := sectionTypeNames[s]; ok {
		return t
	}
	return "UNKNOWN"
}

// GUIDEDSectionAttribute holds a GUIDED section attribute bitfield
type GUIDEDSectionAttribute uint16

// UEFI GUIDED Section Attributes
const (
	GUIDEDSectionProcessingRequired GUIDEDSectionAttribute = 0x01
	GUIDEDSectionAuthStatusValid    GUIDEDSectionAttribute = 0x02
)

// Well-known GUIDs.
var (
	LZMAGUID    = *uuid.MustParse("EE4E5898-3914-4259-9D6E-DC7BD79403CF")
	LZMAX86GUID = *uuid.MustParse("D42AE6BD-1352-4BFB-909A-CA72A6EAE889")
)

// SectionHeader represents an EFI_COMMON_SECTION_HEADER as specified in
// UEFI PI Spec 3.2.4 Firmware File Section
type SectionHeader struct {
	Size [3]uint8 `json:"-"`
	Type SectionType
}

// SectionExtHeader represents an EFI_COMMON_SECTION_HEADER2 as specified in
// UEFI PI Spec 3.2.4 Firmware File Section
type SectionExtHeader struct {
	SectionHeader
	ExtendedSize uint32 `json:"-"`
}

// SectionGUIDDefinedHeader contains the fields for a EFI_SECTION_GUID_DEFINED
// encapsulated section header.
type SectionGUIDDefinedHeader struct {
	GUID       uuid.UUID
	DataOffset uint16
	Attributes uint16
}

// SectionGUIDDefined contains the type specific fields for a
// EFI_SECTION_GUID_DEFINED section.
type SectionGUIDDefined struct {
	SectionGUIDDefinedHeader

	// Metadata
	Compression string
}

// GetBinHeaderLen returns the length of the binary typ specific header
func (s *SectionGUIDDefined) GetBinHeaderLen() uint32 {
	return uint32(unsafe.Sizeof(s.SectionGUIDDefinedHeader))
}

// TypeHeader interface forces type specific headers to report their length
type TypeHeader interface {
	GetBinHeaderLen() uint32
}

// TypeSpecificHeader is used for marshalling and unmarshalling from JSON
type TypeSpecificHeader struct {
	Type   SectionType
	Header TypeHeader
}

var headerTypes = map[SectionType]func() TypeHeader{
	SectionTypeGUIDDefined: func() TypeHeader { return &SectionGUIDDefined{} },
}

// UnmarshalJSON unmarshals a TypeSpecificHeader struct and correctly deduces the
// type of the interface.
func (t *TypeSpecificHeader) UnmarshalJSON(b []byte) error {
	var getType struct {
		Type   SectionType
		Header json.RawMessage
	}
	if err := json.Unmarshal(b, &getType); err != nil {
		return err
	}
	factory, ok := headerTypes[getType.Type]
	if !ok {
		return fmt.Errorf("unknown TypeSpecificHeader type '%v', unable to unmarshal", getType.Type)
	}
	t.Type = SectionType(getType.Type)
	t.Header = factory()
	return json.Unmarshal(getType.Header, &t.Header)
}

// DepExOpCode is one opcode for the dependency expression section.
type DepExOpCode string

// DepExOpCodes maps the numeric code to the string.
var DepExOpCodes = map[byte]DepExOpCode{
	0x0: "BEFORE",
	0x1: "AFTER",
	0x2: "PUSH",
	0x3: "AND",
	0x4: "OR",
	0x5: "NOT",
	0x6: "TRUE",
	0x7: "FALSE",
	0x8: "END",
	0x9: "SOR",
}

// DepExOp contains one operation for the dependency expression.
type DepExOp struct {
	OpCode DepExOpCode
	GUID   *uuid.UUID `json:",omitempty"`
}

// Section represents a Firmware File Section
type Section struct {
	Header SectionExtHeader
	Type   string
	buf    []byte

	// Metadata for extraction and recovery
	ExtractPath string
	FileOrder   int `json:"-"`

	// Type specific fields
	// TODO: It will be simpler if this was not an interface
	TypeSpecific *TypeSpecificHeader `json:",omitempty"`

	// For EFI_SECTION_USER_INTERFACE
	Name string `json:",omitempty"`

	// For EFI_SECTION_DXE_DEPEX, EFI_SECTION_PEI_DEPEX, and EFI_SECTION_MM_DEPEX
	DepEx []DepExOp `json:",omitempty"`

	// Encapsulated firmware
	Encapsulated []*TypedFirmware `json:",omitempty"`
}

// Buf returns the buffer.
// Used mostly for things interacting with the Firmware interface.
func (s *Section) Buf() []byte {
	return s.buf
}

// SetBuf sets the buffer.
// Used mostly for things interacting with the Firmware interface.
func (s *Section) SetBuf(buf []byte) {
	s.buf = buf
}

// Apply calls the visitor on the Section.
func (s *Section) Apply(v Visitor) error {
	return v.Visit(s)
}

// ApplyChildren calls the visitor on each child node of Section.
func (s *Section) ApplyChildren(v Visitor) error {
	for _, f := range s.Encapsulated {
		if err := f.Value.Apply(v); err != nil {
			return err
		}
	}
	return nil
}

// GenSecHeader generates a full binary header for the section data.
// It assumes that the passed in section struct already contains section data in the buffer,
// the section type in the Type field, and the type specific header in the TypeSpecific field.
// It modifies the calling Section.
func (s *Section) GenSecHeader() error {
	var err error
	// Calculate size
	headerLen := uint32(SectionMinLength)
	if s.TypeSpecific != nil && s.TypeSpecific.Header != nil {
		headerLen += s.TypeSpecific.Header.GetBinHeaderLen()
	}
	s.Header.ExtendedSize = uint32(len(s.buf)) + headerLen // TS header lengths are part of headerLen at this point
	if s.Header.ExtendedSize >= 0xFFFFFF {
		headerLen += 4 // Add space for the extended header.
		s.Header.ExtendedSize += 4
	}

	// Set the correct data offset for GUID Defined headers.
	// This is terrible
	if s.Header.Type == SectionTypeGUIDDefined {
		gd := s.TypeSpecific.Header.(*SectionGUIDDefined)
		gd.DataOffset = uint16(headerLen)
		// append type specific header in front of data
		tsh := new(bytes.Buffer)
		if err = binary.Write(tsh, binary.LittleEndian, &gd.SectionGUIDDefinedHeader); err != nil {
			return err
		}
		s.buf = append(tsh.Bytes(), s.buf...)
	}

	// Append common header
	s.Header.Size = Write3Size(uint64(s.Header.ExtendedSize))
	h := new(bytes.Buffer)
	if s.Header.ExtendedSize >= 0xFFFFFF {
		err = binary.Write(h, binary.LittleEndian, &s.Header)
	} else {
		err = binary.Write(h, binary.LittleEndian, &s.Header.SectionHeader)
	}
	if err != nil {
		return err
	}
	s.buf = append(h.Bytes(), s.buf...)
	return nil
}

// Validate File Section
func (s *Section) Validate() []error {
	errs := make([]error, 0)
	buflen := uint32(len(s.buf))
	blankSize := [3]uint8{0xFF, 0xFF, 0xFF}

	// Size Checks
	sh := &s.Header
	if sh.Size == blankSize {
		if buflen < SectionExtMinLength {
			errs = append(errs, fmt.Errorf("section length too small!, buffer is only %#x bytes long for extended header",
				buflen))
			return errs
		}
	} else if uint32(Read3Size(s.Header.Size)) != sh.ExtendedSize {
		errs = append(errs, errors.New("section size not copied into extendedsize"))
		return errs
	}
	if buflen != sh.ExtendedSize {
		errs = append(errs, fmt.Errorf("section size mismatch! Size is %#x, buf length is %#x",
			sh.ExtendedSize, buflen))
		return errs
	}

	return errs
}

// NewSection parses a sequence of bytes and returns a Section
// object, if a valid one is passed, or an error.
func NewSection(buf []byte, fileOrder int) (*Section, error) {
	s := Section{FileOrder: fileOrder}
	// Read in standard header.
	r := bytes.NewReader(buf)
	if err := binary.Read(r, binary.LittleEndian, &s.Header.SectionHeader); err != nil {
		return nil, err
	}

	// Map type to string.
	s.Type = s.Header.Type.String()

	headerSize := unsafe.Sizeof(SectionHeader{})
	if s.Header.Size == [3]uint8{0xFF, 0xFF, 0xFF} {
		// Extended Header
		if err := binary.Read(r, binary.LittleEndian, &s.Header.ExtendedSize); err != nil {
			return nil, err
		}
		if s.Header.ExtendedSize == 0xFFFFFFFF {
			return nil, errors.New("section size and extended size are all FFs! there should not be free space inside a file")
		}
		headerSize = unsafe.Sizeof(SectionExtHeader{})
	} else {
		// Copy small size into big for easier handling.
		// Section's extended size is 32 bits unlike file's
		s.Header.ExtendedSize = uint32(Read3Size(s.Header.Size))
	}

	if buflen := len(buf); int(s.Header.ExtendedSize) > buflen {
		return nil, fmt.Errorf("section size mismatch! Section has size %v, but buffer is %v bytes big",
			s.Header.ExtendedSize, buflen)
	}
	// Slice buffer to the correct size.
	s.buf = buf[:s.Header.ExtendedSize]

	// Section type specific data
	switch s.Header.Type {
	case SectionTypeGUIDDefined:
		typeSpec := &SectionGUIDDefined{}
		if err := binary.Read(r, binary.LittleEndian, &typeSpec.SectionGUIDDefinedHeader); err != nil {
			return nil, err
		}
		s.TypeSpecific = &TypeSpecificHeader{Type: SectionTypeGUIDDefined, Header: typeSpec}

		// Determine how to interpret the section based on the GUID.
		var encapBuf []byte
		if typeSpec.Attributes&uint16(GUIDEDSectionProcessingRequired) != 0 {
			var err error
			switch typeSpec.GUID {
			case LZMAGUID:
				typeSpec.Compression = "LZMA"
				encapBuf, err = lzma.Decode(buf[typeSpec.DataOffset:])
			case LZMAX86GUID:
				typeSpec.Compression = "LZMAX86"
				encapBuf, err = lzma.DecodeX86(buf[typeSpec.DataOffset:])
			default:
				typeSpec.Compression = "UNKNOWN"
			}
			if err != nil {
				log.Print(err)
				typeSpec.Compression = "UNKNOWN"
				encapBuf = []byte{}
			}
		}

		for i, offset := 0, uint64(0); offset < uint64(len(encapBuf)); i++ {
			encapS, err := NewSection(encapBuf[offset:], i)
			if err != nil {
				return nil, fmt.Errorf("error parsing encapsulated section #%d at offset %d: %v",
					i, offset, err)
			}
			// Align to 4 bytes for now. The PI Spec doesn't say what alignment it should be
			// but UEFITool aligns to 4 bytes, and this seems to work on everything I have.
			offset = Align4(offset + uint64(encapS.Header.ExtendedSize))
			s.Encapsulated = append(s.Encapsulated, MakeTyped(encapS))
		}

	case SectionTypeUserInterface:
		s.Name = unicode.UCS2ToUTF8(s.buf[headerSize:])

	case SectionTypeFirmwareVolumeImage:
		fv, err := NewFirmwareVolume(s.buf[headerSize:], 0, true)
		if err != nil {
			return nil, err
		}
		s.Encapsulated = []*TypedFirmware{MakeTyped(fv)}

	case SectionTypeDXEDepEx, SectionTypePEIDepEx, SectionMMDepEx:
		var err error
		if s.DepEx, err = parseDepEx(s.buf[headerSize:]); err != nil {
			log.Println("warning:", err)
		}
	}

	return &s, nil
}

func parseDepEx(b []byte) ([]DepExOp, error) {
	depEx := []DepExOp{}
	r := bytes.NewBuffer(b)
	for {
		opCodeByte, err := r.ReadByte()
		if err != nil {
			return nil, errors.New("invalid DEPEX, no END")
		}
		if opCodeStr, ok := DepExOpCodes[opCodeByte]; ok {
			op := DepExOp{OpCode: opCodeStr}
			if opCodeStr == "BEFORE" || opCodeStr == "AFTER" || opCodeStr == "PUSH" {
				op.GUID = &uuid.UUID{}
				if err := binary.Read(r, binary.LittleEndian, op.GUID); err != nil {
					return nil, fmt.Errorf("invalid DEPEX, could not read GUID: %v", err)
				}
			}
			depEx = append(depEx, op)
			if opCodeStr == "END" {
				break
			}
		} else {
			return nil, fmt.Errorf("invalid DEPEX opcode, %#v", opCodeByte)
		}
	}
	return depEx, nil
}
