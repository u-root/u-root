// Copyright 2021 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package uefi

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

// These headers are coming from EDK2
// MdePkg/Include/Pi/PiFirmwareVolume.h
// MdePkg/Include/Pi/PiFirmwareFile.h

type EFIFirmwareVolumeHeader struct {
	ZeroVector      [16]uint8
	FileSystemGUID  [16]uint8
	FvLength        uint64
	Signature       [4]uint8
	Attributes      uint32
	HeaderLength    uint16
	Checksum        uint16
	ExtHeaderOffset uint16
	Reserved        uint8
	Revision        uint8
}

type EFIFFSFileHeader struct {
	Name           [16]uint8
	HeaderChecksum uint8
	FileChecksum   uint8
	Type           uint8
	Attributes     uint8
	Size           [3]uint8
	State          uint8
}

const (
	EFIFFSAttribLargeFile       uint8 = 0x01
	EFICommonSectionHeaderSize  int   = 4
	EFICommonSectionHeader2Size int   = 8
	EFIFFSFileHeaderSize        int   = 24
	EFIFFSFileHeader2Size       int   = 32
)

const (
	EFISectionTypePE32    uint8 = 0x10
	EFISectionTypeFVImage uint8 = 0x17
)

const (
	EFIFVFileTypeSecurityCore        = 0x03
	EFIFVFileTypeFirmwareVolumeImage = 0x0b
)

// UnmarshalBinary unmarshals the FiwmreareVolumeHeader from binary data.
func (e *EFIFirmwareVolumeHeader) UnmarshalBinary(data []byte) error {
	if len(data) < 0x38 {
		return fmt.Errorf("invalid entry point stucture length %d", len(data))
	}
	if err := binary.Read(bytes.NewReader(data), binary.LittleEndian, e); err != nil {
		return err
	}
	if !bytes.Equal(e.Signature[:], []byte("_FVH")) {
		return fmt.Errorf("invalid Signature string %q", string(e.Signature[:]))
	}
	return nil
}

// UnmarshalBinary unmarshals the EFIFFSFileHeader from binary data.
func (e *EFIFFSFileHeader) UnmarshalBinary(data []byte) error {
	if len(data) < EFIFFSFileHeaderSize {
		return fmt.Errorf("invalid entry point stucture length %d", len(data))
	}
	if err := binary.Read(bytes.NewReader(data), binary.LittleEndian, e); err != nil {
		return err
	}
	return nil
}

// findSecurityCorePEEntry finds SEC PE entry in Firmware Volume
func findSecurityCorePEEntry(data []byte) (offset int, err error) {
	var fvh EFIFirmwareVolumeHeader
	var ffs EFIFFSFileHeader
	if err = fvh.UnmarshalBinary(data); err != nil {
		return 0, err
	}
	offset += int(fvh.HeaderLength)
	for offset < int(fvh.FvLength) {
		if err = ffs.UnmarshalBinary(data[offset:]); err != nil {
			break
		}
		fs := int(ffs.Size[0]) + int(ffs.Size[1])<<8 + int(ffs.Size[2])<<16
		large := ffs.Attributes&EFIFFSAttribLargeFile != 0
		// file size should not be 0
		if fs == 0 {
			return 0, fmt.Errorf("file is corrupt")
		}
		switch ffs.Type {
		case EFIFVFileTypeSecurityCore:
			peo, err := findSectionInFFS(EFISectionTypePE32, data[offset:offset+fs], large)
			if err == nil {
				return offset + peo, nil
			}
		case EFIFVFileTypeFirmwareVolumeImage:
			fvo, err := findSectionInFFS(EFISectionTypeFVImage, data[offset:offset+fs], large)
			if err == nil {
				offset2, err := findSecurityCorePEEntry(data[offset+fvo:])
				if err == nil {
					return offset + fvo + offset2, nil
				}
			}
		}

		// next FFS needs to be aligned with 8 bytes.
		if fs&7 != 0 {
			fs &= ^7
			fs += 8
		}
		offset += fs
	}
	return 0, fmt.Errorf("unable to find SEC ffs in this file")
}

// findSectionInFFS finds given first Given SectionType's entry in FFS
func findSectionInFFS(sectionType uint8, data []byte, isLargeFile bool) (cursor int, err error) {
	fhs := EFIFFSFileHeaderSize
	shs := EFICommonSectionHeaderSize
	if isLargeFile {
		fhs = EFIFFSFileHeader2Size
		shs = EFICommonSectionHeader2Size
	}
	cursor += fhs
	for cursor < len(data) {
		if data[cursor+3] == sectionType {
			return cursor + shs, nil
		}
		ss := int(data[cursor])
		ss += int(data[cursor+1]) << 8
		ss += int(data[cursor+2]) << 16
		if ss == 0 {
			return 0, fmt.Errorf("unable to parse FFS")
		}
		cursor += ss
	}
	return cursor, fmt.Errorf("cannot find PE32 entry in FFS")
}
