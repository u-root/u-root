// Copyright 2018 the LinuxBoot Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package uefi

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"

	"github.com/linuxboot/fiano/pkg/uuid"
)

// FirmwareVolume constants
const (
	FirmwareVolumeFixedHeaderSize  = 56
	FirmwareVolumeMinSize          = FirmwareVolumeFixedHeaderSize + 8 // +8 for the null block that terminates the block list
	FirmwareVolumeExtHeaderMinSize = 20
)

// Valid FV GUIDs
var (
	FFS1      = uuid.MustParse("7a9354d9-0468-444a-81ce-0bf617d890df")
	FFS2      = uuid.MustParse("8c8ce578-8a3d-4f1c-9935-896185c32dd3")
	FFS3      = uuid.MustParse("5473c07a-3dcb-4dca-bd6f-1e9689e7349a")
	EVSA      = uuid.MustParse("fff12b8d-7696-4c8b-a985-2747075b4f50")
	NVAR      = uuid.MustParse("cef5b9a3-476d-497f-9fdc-e98143e0422c")
	EVSA2     = uuid.MustParse("00504624-8a59-4eeb-bd0f-6b36e96128e0")
	AppleBoot = uuid.MustParse("04adeead-61ff-4d31-b6ba-64f8bf901f5a")
	PFH1      = uuid.MustParse("16b45da2-7d70-4aea-a58d-760e9ecb841d")
	PFH2      = uuid.MustParse("e360bdba-c3ce-46be-8f37-b231e5cb9f35")
)

// FVGUIDs holds common FV type names
var FVGUIDs = map[uuid.UUID]string{
	*FFS1:      "FFS1",
	*FFS2:      "FFS2",
	*FFS3:      "FFS3",
	*EVSA:      "NVRAM_EVSA",
	*NVAR:      "NVRAM_NVAR",
	*EVSA2:     "NVRAM_EVSA2",
	*AppleBoot: "APPLE_BOOT",
	*PFH1:      "PFH1",
	*PFH2:      "PFH2",
}

// These are the FVs we actually try to parse beyond the header
// We don't parse anything except FFS2 and FFS3
var supportedFVs = map[uuid.UUID]bool{
	*FFS2: true,
	*FFS3: true,
}

// Block describes number and size of the firmware volume blocks
type Block struct {
	Count uint32
	Size  uint32
}

// FirmwareVolumeFixedHeader contains the fixed fields of a firmware volume
// header
type FirmwareVolumeFixedHeader struct {
	_               [16]uint8
	FileSystemGUID  uuid.UUID
	Length          uint64
	Signature       uint32
	Attributes      uint32 // UEFI PI spec volume 3.2.1 EFI_FIRMWARE_VOLUME_HEADER
	HeaderLen       uint16
	Checksum        uint16
	ExtHeaderOffset uint16
	Reserved        uint8 `json:"-"`
	Revision        uint8
	// _               [3]uint8
}

// FirmwareVolumeExtHeader contains the fields of an extended firmware volume
// header
type FirmwareVolumeExtHeader struct {
	FVName        uuid.UUID
	ExtHeaderSize uint32
}

// FirmwareVolume represents a firmware volume. It combines the fixed header and
// a variable list of blocks
type FirmwareVolume struct {
	FirmwareVolumeFixedHeader
	// there must be at least one that is zeroed and indicates the end of the
	// block list
	// We don't really have to care about blocks because we just read everything in.
	Blocks []Block
	FirmwareVolumeExtHeader
	Files []*File `json:",omitempty"`

	// Variables not in the binary for us to keep track of stuff/print
	DataOffset  uint64
	fvType      string
	buf         []byte
	FVOffset    uint64 // Byte offset from start of BIOS region.
	ExtractPath string
	Resizable   bool // Determines if this FV is resizable.
}

// Buf returns the buffer.
// Used mostly for things interacting with the Firmware interface.
func (fv *FirmwareVolume) Buf() []byte {
	return fv.buf
}

// SetBuf sets the buffer.
// Used mostly for things interacting with the Firmware interface.
func (fv *FirmwareVolume) SetBuf(buf []byte) {
	fv.buf = buf
}

// Apply calls the visitor on the FirmwareVolume.
func (fv *FirmwareVolume) Apply(v Visitor) error {
	return v.Visit(fv)
}

// ApplyChildren calls the visitor on each child node of FirmwareVolume.
func (fv *FirmwareVolume) ApplyChildren(v Visitor) error {
	for _, f := range fv.Files {
		if err := f.Apply(v); err != nil {
			return err
		}
	}
	return nil
}

// GetErasePolarity gets the erase polarity
func (fv *FirmwareVolume) GetErasePolarity() uint8 {
	if fv.Attributes&0x800 != 0 {
		return 0xFF
	}
	return 0
}

// Validate Firmware Volume
func (fv *FirmwareVolume) Validate() []error {
	// TODO: Add more verification if needed.
	errs := make([]error, 0)
	// Check for min length
	fvlen := uint64(len(fv.buf))
	// We need this check in case HeaderLen doesn't exist, and bail out early
	if fvlen < FirmwareVolumeMinSize {
		errs = append(errs, fmt.Errorf("length too small!, buffer is only %#x bytes long", fvlen))
		return errs
	}
	// Check header length
	if fv.HeaderLen < FirmwareVolumeMinSize {
		errs = append(errs, fmt.Errorf("header length too small, got: %#x", fv.HeaderLen))
		return errs
	}
	// Check for full header and bail out if its not fully formed.
	if fvlen < uint64(fv.HeaderLen) {
		errs = append(errs, fmt.Errorf("buffer smaller than header!, header is %#x bytes, buffer is %#x bytes",
			fv.HeaderLen, fvlen))
		return errs
	}
	// Do we want to fail in this case? maybe not.
	if FVGUIDs[fv.FileSystemGUID] == "" {
		errs = append(errs, fmt.Errorf("unknown FV type! Guid was %v", fv.FileSystemGUID))
	}
	// UEFI PI spec says version should always be 2
	if fv.Revision != 2 {
		errs = append(errs, fmt.Errorf("revision should be 2, was %v", fv.Revision))
	}
	// Check Signature
	fvSigInt := binary.LittleEndian.Uint32([]byte("_FVH"))
	if fv.Signature != fvSigInt {
		errs = append(errs, fmt.Errorf("signature was not _FVH, got: %#08x", fv.Signature))
	}
	// Check length
	if fv.Length != fvlen {
		errs = append(errs, fmt.Errorf("length mismatch!, header has %#x, buffer is %#x bytes long", fv.Length, fvlen))
	}
	// Check checksum
	sum, err := Checksum16(fv.buf[:fv.HeaderLen])
	if err != nil {
		errs = append(errs, fmt.Errorf("unable to checksum FV header: %v", err))
	} else if sum != 0 {
		errs = append(errs, fmt.Errorf("header did not sum to 0, got: %#x", sum))
	}

	for _, f := range fv.Files {
		errs = append(errs, f.Validate()...)
	}
	return errs
}

func fillFFs(b []byte) {
	for i := range b {
		b[i] = 0xFF
	}
}

// InsertFile appends the file to the end of the buffer according to alignment requirements.
func (fv *FirmwareVolume) InsertFile(alignedOffset uint64, fBuf []byte) error {
	// fv.Length should contain the minimum fv size.
	// If Resizable is not set, this is the exact FV size.
	fvLen := fv.Length
	if !fv.Resizable && alignedOffset > fv.Length {
		return fmt.Errorf("insufficient space in %#x bytes FV, files too big, offset was %#x",
			fvLen, alignedOffset)
	}
	bufLen := uint64(len(fv.buf))
	if bufLen > alignedOffset {
		return fmt.Errorf("aligned offset is in the middle of the FV, offset was %#x, fv buffer was %#x",
			alignedOffset, bufLen)
	}

	// add padding for alignment
	for i, num := uint64(0), alignedOffset-bufLen; i < num; i++ {
		fv.buf = append(fv.buf, Attributes.ErasePolarity)
	}

	// Check size
	fLen := uint64(len(fBuf))
	if fLen == 0 {
		log.Fatal("FUCK")
	}
	if fLen+alignedOffset > fvLen && !fv.Resizable {
		// TODO: Actually loop through and calculate the full size so we know how much to reduce by.
		// For now we just return early
		return fmt.Errorf("insufficient space in %#x bytes FV, files too big, offset was %#x, length was %#x",
			fvLen, alignedOffset, len(fBuf))
	}
	// Overwrite old data in the firmware volume.
	fv.buf = append(fv.buf, fBuf...)
	return nil
}

// FindFirmwareVolumeOffset searches for a firmware volume signature, "_FVH"
// using 8-byte alignment. If found, returns the offset from the start of the
// bios region, otherwise returns -1.
func FindFirmwareVolumeOffset(data []byte) int64 {
	if len(data) < 32 {
		return -1
	}
	var (
		offset int64
		fvSig  = []byte("_FVH")
	)
	for offset = 32; offset < int64(len(data)); offset += 8 {
		if bytes.Equal(data[offset:offset+4], fvSig) {
			return offset - 40 // the actual volume starts 40 bytes before the signature
		}
	}
	return -1
}

// NewFirmwareVolume parses a sequence of bytes and returns a FirmwareVolume
// object, if a valid one is passed, or an error
func NewFirmwareVolume(data []byte, fvOffset uint64, resizable bool) (*FirmwareVolume, error) {
	fv := FirmwareVolume{Resizable: resizable}

	if len(data) < FirmwareVolumeMinSize {
		return nil, fmt.Errorf("Firmware Volume size too small: expected %v bytes, got %v",
			FirmwareVolumeMinSize,
			len(data),
		)
	}
	reader := bytes.NewReader(data)
	if err := binary.Read(reader, binary.LittleEndian, &fv.FirmwareVolumeFixedHeader); err != nil {
		return nil, err
	}
	// read the block map
	blocks := make([]Block, 0)
	for {
		var block Block
		if err := binary.Read(reader, binary.LittleEndian, &block); err != nil {
			return nil, err
		}
		if block.Count == 0 && block.Size == 0 {
			// found the terminating block
			break
		}
		blocks = append(blocks, block)
	}
	fv.Blocks = blocks

	// Parse the extended header and figure out the start of data
	fv.DataOffset = uint64(fv.HeaderLen)
	if fv.ExtHeaderOffset != 0 && uint64(fv.ExtHeaderOffset) < fv.Length-FirmwareVolumeExtHeaderMinSize {
		// jump to ext header offset.
		r := bytes.NewReader(data[fv.ExtHeaderOffset:])
		if err := binary.Read(r, binary.LittleEndian, &fv.FirmwareVolumeExtHeader); err != nil {
			return nil, fmt.Errorf("unable to parse FV extended header, got: %v", err)
		}
		// TODO: will the ext header ever end before the regular header? I don't believe so. Add a check?
		fv.DataOffset = uint64(fv.ExtHeaderOffset) + uint64(fv.ExtHeaderSize)
	}
	// Make sure DataOffset is 8 byte aligned at least.
	// TODO: handle alignment field in header.
	fv.DataOffset = Align8(fv.DataOffset)

	fv.fvType = FVGUIDs[fv.FileSystemGUID]
	fv.FVOffset = fvOffset

	// slice the buffer
	fv.buf = data[:fv.Length]

	// Parse the files.
	// TODO: handle fv data alignment.
	// Start from the end of the fv header.
	// Test if the fv type is supported.
	if _, ok := supportedFVs[fv.FileSystemGUID]; !ok {
		return &fv, nil
	}
	lh := fv.Length - FileHeaderMinLength
	var prevLen uint64
	for offset := fv.DataOffset; offset < lh; offset += prevLen {
		offset = Align8(offset)
		file, err := NewFile(data[offset:])
		if err != nil {
			return nil, fmt.Errorf("unable to construct firmware file at offset %#x into FV: %v", offset, err)
		}
		if file == nil {
			// We've reached free space. Terminate
			break
		}
		fv.Files = append(fv.Files, file)
		prevLen = file.Header.ExtendedSize
	}
	return &fv, nil
}
