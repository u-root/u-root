// Copyright 2018 the LinuxBoot Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package uefi

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"github.com/linuxboot/fiano/pkg/uuid"
)

// FVFileType represents the different types possible in an EFI file.
type FVFileType uint8

// UEFI FV File types.
const (
	FVFileTypeAll FVFileType = iota
	FVFileTypeRaw
	FVFileTypeFreeForm
	FVFileTypeSECCore
	FVFileTypePEICore
	FVFileTypeDXECore
	FVFileTypePEIM
	FVFileTypeDriver
	FVFileTypeCombinedPEIMDriver
	FVFileTypeApplication
	FVFileTypeSMM
	FVFileTypeVolumeImage
	FVFileTypeCombinedSMMDXE
	FVFileTypeSMMCore
	FVFileTypeSMMStandalone
	FVFileTypeSMMCoreStandalone
	FVFileTypeOEMMin   FVFileType = 0xC0
	FVFileTypeOEMMax   FVFileType = 0xDF
	FVFileTypeDebugMin FVFileType = 0xE0
	FVFileTypeDebugMax FVFileType = 0xEF
	FVFileTypePad      FVFileType = 0xF0
	FVFileTypeFFSMin   FVFileType = 0xF0
	FVFileTypeFFSMax   FVFileType = 0xFF
)

// SupportedFiles is a list of files types which will be parsed. File types not
// on this list are treated as opaque binary blobs.
var SupportedFiles = map[FVFileType]bool{
	// These are the file types that we'll actually try to parse sections for.
	FVFileTypeFreeForm: true,
	FVFileTypeSECCore:  true,
	FVFileTypePEICore:  true,
	FVFileTypeDXECore:  true,
	// TODO: Commenting out this line prevents PEI modules from being
	// decompressed. This solves the problem of PEI being too big when recompressed.
	//FVFileTypePEIM:               true,
	FVFileTypeDriver:             true,
	FVFileTypeCombinedPEIMDriver: true,
	FVFileTypeApplication:        true,
	FVFileTypeSMM:                true,
	FVFileTypeVolumeImage:        true,
	FVFileTypeCombinedSMMDXE:     true,
	FVFileTypeSMMCore:            true,
	FVFileTypeSMMStandalone:      true,
	FVFileTypeSMMCoreStandalone:  true,
}

var fileTypeNames = map[FVFileType]string{
	FVFileTypeRaw:                "EFI_FV_FILETYPE_RAW",
	FVFileTypeFreeForm:           "EFI_FV_FILETYPE_FREEFORM",
	FVFileTypeSECCore:            "EFI_FV_FILETYPE_SECURITY_CORE",
	FVFileTypePEICore:            "EFI_FV_FILETYPE_PEI_CORE",
	FVFileTypeDXECore:            "EFI_FV_FILETYPE_DXE_CORE",
	FVFileTypePEIM:               "EFI_FV_FILETYPE_PEIM",
	FVFileTypeDriver:             "EFI_FV_FILETYPE_DRIVER",
	FVFileTypeCombinedPEIMDriver: "EFI_FV_FILETYPE_COMBINED_PEIM_DRIVER",
	FVFileTypeApplication:        "EFI_FV_FILETYPE_APPLICATION",
	FVFileTypeSMM:                "EFI_FV_FILETYPE_MM",
	FVFileTypeVolumeImage:        "EFI_FV_FILETYPE_FIRMWARE_VOLUME_IMAGE",
	FVFileTypeCombinedSMMDXE:     "EFI_FV_FILETYPE_COMBINED_MM_DXE",
	FVFileTypeSMMCore:            "EFI_FV_FILETYPE_MM_CORE",
	FVFileTypeSMMStandalone:      "EFI_FV_FILETYPE_MM_STANDALONE",
	FVFileTypeSMMCoreStandalone:  "EFI_FV_FILETYPE_MM_CORE_STANDALONE",
}

// String creates a string representation for the file type.
func (f FVFileType) String() string {
	switch {
	case FVFileTypeOEMMin <= f && f <= FVFileTypeOEMMax:
		return fmt.Sprintf("EFI_FV_FILETYPE_OEM (%#x)", uint8(f))
	case FVFileTypeDebugMin <= f && f <= FVFileTypeDebugMax:
		return fmt.Sprintf("EFI_FV_FILETYPE_DEBUG (%#x)", uint8(f))
	// We use the non-inclusive '<' operator here because pad files belong
	// to the FFS filetype, but are also their own type.
	case FVFileTypeFFSMin < f && f <= FVFileTypeFFSMax:
		return fmt.Sprintf("EFI_FV_FILETYPE_FFS (%#x)", uint8(f))
	case f == FVFileTypePad:
		return "EFI_FV_FILETYPE_FFS_PAD"
	}
	if t, ok := fileTypeNames[f]; ok {
		return t
	}
	return "UNKNOWN"
}

// Stock GUIDS
var (
	ZeroGUID = uuid.MustParse("00000000-0000-0000-0000-000000000000")
	FFGUID   = uuid.MustParse("FFFFFFFF-FFFF-FFFF-FFFF-FFFFFFFFFFFF")
)

// FileAlignments specifies the correct alignments based on the field in the file header.
var fileAlignments = []uint64{
	// These alignments not computable, we have to look them up.
	1,
	16,
	128,
	512,
	1024,
	4 * 1024,
	32 * 1024,
	64 * 1024,
	128 * 1024,
	256 * 1024,
	512 * 1024,
	1024 * 1024,
	2 * 1024 * 1024,
	4 * 1024 * 1024,
	8 * 1024 * 1024,
	16 * 1024 * 1024,
}

const (
	// FileHeaderMinLength is the minimum length of a firmware file header.
	FileHeaderMinLength = 0x18
	// FileHeaderExtMinLength is the minimum length of an extended firmware file header.
	FileHeaderExtMinLength = 0x20
	// EmptyBodyChecksum is the value placed in the File IntegrityCheck field if the body checksum bit isn't set.
	EmptyBodyChecksum uint8 = 0xAA
)

// IntegrityCheck holds the two 8 bit checksums for the file header and body separately.
type IntegrityCheck struct {
	Header uint8
	File   uint8
}

type fileAttr uint8

// FileHeader represents an EFI File header.
type FileHeader struct {
	UUID       uuid.UUID      // This is the GUID of the file.
	Checksum   IntegrityCheck `json:"-"`
	Type       FVFileType
	Attributes fileAttr
	Size       [3]uint8 `json:"-"`
	State      uint8    `json:"-"`
}

// Checks if the large file attribute is set
func (a fileAttr) isLarge() bool {
	return a&0x01 != 0
}

// GetAlignment returns the byte alignment specified by the file header.
func (a fileAttr) GetAlignment() uint64 {
	alignVal := (a & 0x38) >> 3
	alignVal |= (a & 0x02) << 2
	return fileAlignments[alignVal]
}

// Sets the large file attribute.
func (a *fileAttr) setLarge(large bool) {
	if large {
		*a |= 0x01
	} else {
		*a &= 0xFE
	}
}

// Checks if we need to checksum the file body
func (a fileAttr) hasChecksum() bool {
	return a&0x40 != 0
}

// HeaderLen is a helper function to return the length of the file header
// depending on the file size
func (f *File) HeaderLen() uint64 {
	if f.Header.Attributes.isLarge() {
		return FileHeaderExtMinLength
	}
	return FileHeaderMinLength
}

func (f *File) checksumHeader() uint8 {
	fh := f.Header
	headerSize := FileHeaderMinLength
	if fh.Attributes.isLarge() {
		headerSize = FileHeaderExtMinLength
	}
	// Sum over header without State and IntegrityCheck.File.
	// To do that we just sum over the whole header and subtract.
	// UEFI PI Spec 3.2.3 EFI_FFS_FILE_HEADER
	sum := Checksum8(f.buf[:headerSize])
	sum -= fh.Checksum.File
	sum -= fh.State
	return sum
}

// FileHeaderExtended represents an EFI File header with the
// large file attribute set.
// We also use this as the generic header for all EFI files, regardless of whether
// they are actually large. This makes it easier for us to just return one type
// All sizes are also copied into the ExtendedSize field so we only have to check once
type FileHeaderExtended struct {
	FileHeader
	ExtendedSize uint64 `json:"-"`
}

// File represents an EFI File.
type File struct {
	Header   FileHeaderExtended
	Type     string
	Sections []*Section `json:",omitempty"`

	//Metadata for extraction and recovery
	buf         []byte
	ExtractPath string
	DataOffset  uint64
}

// Buf returns the buffer.
// Used mostly for things interacting with the Firmware interface.
func (f *File) Buf() []byte {
	return f.buf
}

// SetBuf sets the buffer.
// Used mostly for things interacting with the Firmware interface.
func (f *File) SetBuf(buf []byte) {
	f.buf = buf
}

// Apply calls the visitor on the File.
func (f *File) Apply(v Visitor) error {
	return v.Visit(f)
}

// ApplyChildren calls the visitor on each child node of File.
func (f *File) ApplyChildren(v Visitor) error {
	for _, s := range f.Sections {
		if err := s.Apply(v); err != nil {
			return err
		}
	}
	return nil
}

// SetSize sets the size into the File struct.
// If resizeFile is true, if the file is too large the file will be enlarged to make space
// for the ExtendedHeader
func (f *File) SetSize(size uint64, resizeFile bool) {
	fh := &f.Header
	// See if we need the extended size
	// Check if size > 3 bytes size field
	fh.ExtendedSize = size
	fh.Attributes.setLarge(false)
	if fh.ExtendedSize > 0xFFFFFF {
		// Can't fit, need extended header
		if resizeFile {
			// Increase the file size by the additional space needed
			// for the extended header.
			fh.ExtendedSize += FileHeaderExtMinLength - FileHeaderMinLength
		}
		fh.Attributes.setLarge(true)
	}
	// This will set size to 0xFFFFFF if too big.
	fh.Size = Write3Size(fh.ExtendedSize)
}

// ChecksumAndAssemble takes in the fileData and assembles the file binary
func (f *File) ChecksumAndAssemble(fileData []byte) error {
	// Checksum the header and body, then write out the header.
	// To checksum the header we write the temporary header to the file buffer first.
	fh := &f.Header

	header := new(bytes.Buffer)
	err := binary.Write(header, binary.LittleEndian, fh)
	if err != nil {
		return fmt.Errorf("unable to construct binary header of file %v, got %v",
			fh.UUID, err)
	}
	f.buf = header.Bytes()
	// We need to get rid of whatever it sums to so that the overall sum is zero
	// Sorry about the name :(
	fh.Checksum.Header -= f.checksumHeader()

	// Checksum the body
	fh.Checksum.File = EmptyBodyChecksum
	if fh.Attributes.hasChecksum() {
		// if the empty checksum had been set to 0 instead of 0xAA
		// this could have been a bit nicer. BUT NOOOOOOO.
		fh.Checksum.File = 0 - Checksum8(fileData)
	}

	// Write out the updated header to the buffer with the new checksums.
	// Write the extended header only if the large attribute flag is set.
	header = new(bytes.Buffer)
	if fh.Attributes.isLarge() {
		err = binary.Write(header, binary.LittleEndian, fh)
	} else {
		err = binary.Write(header, binary.LittleEndian, fh.FileHeader)
	}
	if err != nil {
		return err
	}
	f.buf = header.Bytes()

	f.buf = append(f.buf, fileData...)
	return nil
}

// Validate Firmware File
func (f *File) Validate() []error {
	errs := make([]error, 0)
	buflen := uint64(len(f.buf))
	blankSize := [3]uint8{0xFF, 0xFF, 0xFF}
	if buflen < FileHeaderMinLength {
		errs = append(errs, fmt.Errorf("file length too small!, buffer is only %#x bytes long", buflen))
		return errs
	}

	// Size Checks
	fh := &f.Header
	if fh.Size == blankSize {
		if buflen < FileHeaderExtMinLength {
			errs = append(errs, fmt.Errorf("file %v length too small!, buffer is only %#x bytes long for extended header",
				fh.UUID, buflen))
			return errs
		}
		if !fh.Attributes.isLarge() {
			errs = append(errs, fmt.Errorf("file %v using extended header, but large attribute is not set",
				fh.UUID))
			return errs
		}
	} else if Read3Size(f.Header.Size) != fh.ExtendedSize {
		errs = append(errs, fmt.Errorf("file %v size not copied into extendedsize",
			fh.UUID))
		return errs
	}
	if buflen != fh.ExtendedSize {
		errs = append(errs, fmt.Errorf("file %v size mismatch! Size is %#x, buf length is %#x",
			fh.UUID, fh.ExtendedSize, buflen))
		return errs
	}

	// Header Checksums
	if sum := f.checksumHeader(); sum != 0 {
		errs = append(errs, fmt.Errorf("file %v header checksum failure! sum was %v",
			fh.UUID, sum))
	}

	// Body Checksum
	if !fh.Attributes.hasChecksum() && fh.Checksum.File != EmptyBodyChecksum {
		errs = append(errs, fmt.Errorf("file %v body checksum failure! Attribute was not set, but sum was %v instead of %v",
			fh.UUID, fh.Checksum.File, EmptyBodyChecksum))
	} else if fh.Attributes.hasChecksum() {
		headerSize := FileHeaderMinLength
		if fh.Attributes.isLarge() {
			headerSize = FileHeaderExtMinLength
		}
		if sum := Checksum8(f.buf[headerSize:]); sum != 0 {
			errs = append(errs, fmt.Errorf("file %v body checksum failure! sum was %v",
				fh.UUID, sum))
		}
	}

	for _, s := range f.Sections {
		errs = append(errs, s.Validate()...)
	}
	return errs
}

// CreatePadFile creates an empty pad file in order to align the next file.
func CreatePadFile(size uint64) (*File, error) {
	if size < FileHeaderMinLength {
		return nil, fmt.Errorf("size too small! min size required is %#x bytes, requested %#x",
			FileHeaderMinLength, size)
	}

	f := File{}
	fh := &f.Header

	// Create empty guid
	if Attributes.ErasePolarity == 0xFF {
		fh.UUID = *FFGUID
	} else if Attributes.ErasePolarity == 0 {
		fh.UUID = *ZeroGUID
	} else {
		return nil, fmt.Errorf("erase polarity not 0x00 or 0xFF, got %#x", Attributes.ErasePolarity)
	}

	// TODO: I see examples of this where the attributes are just 0 and not dependent on the
	// erase polarity. Is that right? Check and handle.
	fh.Attributes = 0

	// Set the size. If the file is too big, we take up more of the padding for the header.
	// This also sets the large file attribute if file is big.
	f.SetSize(size, false)
	fh.Type = FVFileTypePad

	// Create empty pad filedata based on size
	var fileData []byte
	fileData = make([]byte, size-FileHeaderMinLength)
	if fh.Attributes.isLarge() {
		fileData = make([]byte, size-FileHeaderExtMinLength)
	}
	// Fill with empty bytes
	for i, dataLen := 0, len(fileData); i < dataLen; i++ {
		fileData[i] = Attributes.ErasePolarity
	}

	fh.State = 0x07 ^ Attributes.ErasePolarity

	// Everything has been setup. Checksum and create.
	if err := f.ChecksumAndAssemble(fileData); err != nil {
		return nil, err
	}
	return &f, nil
}

// NewFile parses a sequence of bytes and returns a File
// object, if a valid one is passed, or an error. If no error is returned and the File
// pointer is nil, it means we've reached the volume free space at the end of the FV.
func NewFile(buf []byte) (*File, error) {
	f := File{}
	f.DataOffset = FileHeaderMinLength
	// Read in standard header.
	r := bytes.NewReader(buf)
	if err := binary.Read(r, binary.LittleEndian, &f.Header.FileHeader); err != nil {
		return nil, err
	}

	// Map type to string.
	f.Type = f.Header.Type.String()

	// TODO: Check Attribute flag as well. How important is the attribute flag? we already
	// have FFFFFF in the size
	if f.Header.Size == [3]uint8{0xFF, 0xFF, 0xFF} {
		// Extended Header
		if err := binary.Read(r, binary.LittleEndian, &f.Header.ExtendedSize); err != nil {
			return nil, err
		}
		if f.Header.ExtendedSize == 0xFFFFFFFFFFFFFFFF {
			// Start of free space
			// Note: this is not a pad file. Pad files also have valid headers.
			return nil, nil
		}
		f.DataOffset = FileHeaderExtMinLength
	} else {
		// Copy small size into big for easier handling.
		// Damn the 3 byte sizes.
		f.Header.ExtendedSize = Read3Size(f.Header.Size)
	}

	if buflen := len(buf); f.Header.ExtendedSize > uint64(buflen) {
		return nil, fmt.Errorf("File size too big! File with GUID: %v has length %v, but is only %v bytes big",
			f.Header.UUID, f.Header.ExtendedSize, buflen)
	}
	// Slice buffer to the correct size.
	f.buf = buf[:f.Header.ExtendedSize]

	// Parse sections
	if _, ok := SupportedFiles[f.Header.Type]; !ok {
		return &f, nil
	}
	for i, offset := 0, f.DataOffset; offset < f.Header.ExtendedSize; i++ {
		s, err := NewSection(f.buf[offset:], i)
		if err != nil {
			return nil, fmt.Errorf("error parsing sections of file %v: %v", f.Header.UUID, err)
		}
		offset += uint64(s.Header.ExtendedSize)
		// Align to 4 bytes for now. The PI Spec doesn't say what alignment it should be
		// but UEFITool aligns to 4 bytes, and this seems to work on everything I have.
		offset = Align4(offset)
		f.Sections = append(f.Sections, s)
	}

	return &f, nil
}
