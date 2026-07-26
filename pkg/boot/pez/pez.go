// Copyright 2026 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package pez contains a Extractor for the PE compressed Linux Image (vmlinuz, ZBOOT).
package pez

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"

	"github.com/klauspost/compress/zstd"
)

// Linux PE zboot header as defined in Linux source
// drivers/firmware/efi/libstub/zboot-header.S
type Header struct {
	Magic       uint32
	Type        uint32
	Offset      uint32
	Size        uint32
	Reserved    [2]uint32
	Compression [4]byte
}

const (
	magic      = 0x00005a4d
	typeZImage = 0x676d697a // "zimg"

	// dosHeaderPEOffsetAddr is the offset in the MS-DOS header (0x3c) where
	// the 4-byte offset to the PE header is located.
	dosHeaderPEOffsetAddr = 0x3c

	// peHeaderOffsetSize is the size of the PE header offset pointer (4 bytes).
	peHeaderOffsetSize = 4

	// peSignature is the 4-byte signature "PE\x00\x00" that marks the start of the PE header.
	peSignature    = "PE\x00\x00"
	peSignatureLen = 4

	// coffNumSectionsOffset is the offset in bytes from the start of the PE header
	// to the NumberOfSections field in the COFF File Header.
	// PE Signature (4 bytes) + Machine (2 bytes) = 6 bytes.
	coffNumSectionsOffset = 6
	coffNumSectionsSize   = 2

	// coffOptHeaderSizeOffset is the offset in bytes from the start of the PE header
	// to the SizeOfOptionalHeader field in the COFF File Header.
	// PE Signature (4 bytes) + COFF File Header fields up to SizeOfOptionalHeader (16 bytes) = 20 bytes.
	coffOptHeaderSizeOffset = 20
	coffOptHeaderSizeSize   = 2

	// coffHeaderSize is the size of the combined PE signature (4 bytes) and COFF File Header (20 bytes).
	coffHeaderSize = 24

	// peSectionEntrySize is the size of each section header entry in the Section Table (40 bytes).
	peSectionEntrySize = 40

	// peSectionNameSize is the maximum size of a section name in the section header (8 bytes).
	peSectionNameSize = 8

	// peSectionRawSizeOffset is the offset of the SizeOfRawData field relative to the start of a section header entry.
	peSectionRawSizeOffset = 16

	// peSectionRawOffsetOffset is the offset of the PointerToRawData field relative to the start of a section header entry.
	peSectionRawOffsetOffset = 20

	// pezHeaderSize is the total size of the PEZ (ZBOOT) compressed image header (28 bytes).
	pezHeaderSize = 28
)

var ErrImageTooSmall = errors.New("image too small")
var ErrMagicMismatch = errors.New("magic number mismatch")
var ErrNotZImage = errors.New("not a zimg")
var ErrUnsupportedCompression = errors.New("unsupported compression type")

// readAt is a helper function that reads a specified number of bytes from the io.ReaderAt stream.
// It handles EOF, unexpected EOF, and partial reads cleanly.
func readAt(img io.ReaderAt, offset int64, size int) ([]byte, error) {
	buf := make([]byte, size)
	n, err := img.ReadAt(buf, offset)
	if err != nil && err != io.EOF {
		return nil, err
	}
	// If we read 0 bytes because we are exactly at the end of the file (offset == file size),
	// we want to return a clean io.EOF rather than io.ErrUnexpectedEOF.
	if n == 0 && err == io.EOF {
		return nil, io.EOF
	}
	if n < size {
		return nil, io.ErrUnexpectedEOF
	}
	return buf, nil
}

// findLinuxSection scans the PE/COFF image header and section table to locate the ".linux"
// section. This section houses the compressed ARM64 EFI zboot kernel payload.
// Returns the file offset and size of the section if found, or an error.
func findLinuxSection(img io.ReaderAt) (int64, int64, error) {
	// Read the offset pointing to the start of the PE header from the DOS stub (0x3c).
	buf, err := readAt(img, dosHeaderPEOffsetAddr, peHeaderOffsetSize)
	if err != nil {
		return 0, 0, err
	}
	peOffset := binary.LittleEndian.Uint32(buf)

	// Read and verify the PE signature ("PE\x00\x00") at the start of the PE header.
	buf, err = readAt(img, int64(peOffset), peSignatureLen)
	if err != nil {
		return 0, 0, err
	}
	if string(buf) != peSignature {
		return 0, 0, fmt.Errorf("invalid PE signature")
	}

	// Read the number of sections from the COFF File Header.
	buf, err = readAt(img, int64(peOffset)+coffNumSectionsOffset, coffNumSectionsSize)
	if err != nil {
		return 0, 0, err
	}
	numSections := binary.LittleEndian.Uint16(buf)

	// Read the optional header size from the COFF File Header to compute where the Section Table begins.
	buf, err = readAt(img, int64(peOffset)+coffOptHeaderSizeOffset, coffOptHeaderSizeSize)
	if err != nil {
		return 0, 0, err
	}
	optHeaderSize := binary.LittleEndian.Uint16(buf)

	// The Section Table starts immediately after the Optional Header (PE Signature + COFF Header + Optional Header).
	sectionTableOffset := int64(peOffset) + coffHeaderSize + int64(optHeaderSize)
	for i := 0; i < int(numSections); i++ {
		// Seek to each 40-byte section entry and read it.
		entryOffset := sectionTableOffset + int64(i*peSectionEntrySize)
		buf, err = readAt(img, entryOffset, peSectionEntrySize)
		if err != nil {
			return 0, 0, err
		}

		// Read the 8-character section name and trim trailing null bytes.
		name := buf[0:peSectionNameSize]
		nameLen := 0
		for nameLen < peSectionNameSize && name[nameLen] != 0 {
			nameLen++
		}
		nameStr := string(name[:nameLen])

		// If this is the ".linux" section, extract its offset and size.
		if nameStr == ".linux" {
			rawSize := binary.LittleEndian.Uint32(buf[peSectionRawSizeOffset : peSectionRawSizeOffset+4])
			rawOffset := binary.LittleEndian.Uint32(buf[peSectionRawOffsetOffset : peSectionRawOffsetOffset+4])
			return int64(rawOffset), int64(rawSize), nil
		}
	}

	return 0, 0, fmt.Errorf("no .linux section found")
}

// Extract extracts and decompresses the embedded bootable ARM64 kernel payload
// from either a raw ZBOOT EFI executable or an outer PE/COFF Unified Kernel Image (UKI).
// It returns an io.ReaderAt stream for the raw decompressed kernel image.
func Extract(img io.ReaderAt) (io.ReaderAt, error) {
	// First, check if the input image is wrapped inside a PE/COFF executable.
	// If so, locate and extract the ".linux" section.
	if offset, size, err := findLinuxSection(img); err == nil {
		img = io.NewSectionReader(img, offset, size)
	}

	// Read the 28-byte PEZ / ZBOOT header.
	buf, err := readAt(img, 0, pezHeaderSize)
	if err != nil {
		return nil, err
	}

	header := Header{
		Magic:  binary.LittleEndian.Uint32(buf[0:4]),
		Type:   binary.LittleEndian.Uint32(buf[4:8]),
		Offset: binary.LittleEndian.Uint32(buf[8:12]),
		Size:   binary.LittleEndian.Uint32(buf[12:16]),
	}
	header.Reserved[0] = binary.LittleEndian.Uint32(buf[16:20])
	header.Reserved[1] = binary.LittleEndian.Uint32(buf[20:24])
	copy(header.Compression[:], buf[24:28])

	// Validate the ZBOOT magic and image type.
	if header.Magic != magic {
		return nil, fmt.Errorf("found magic %#x but PEZ expects %#x: %w", header.Magic, magic, ErrMagicMismatch)
	}
	if header.Type != typeZImage {
		return nil, fmt.Errorf("found type %#x but PEZ expects %#x: %w", header.Type, typeZImage, ErrNotZImage)
	}

	// Slice the raw compressed payload based on header offset and size.
	payload := io.NewSectionReader(img, int64(header.Offset), int64(header.Size))

	// Trim trailing null bytes from the compression algorithm name without extra allocations.
	compression := string(bytes.TrimRight(header.Compression[:], "\x00"))

	// Decompress the payload according to its compression algorithm.
	switch compression {
	case "zstd":
		r, err := zstd.NewReader(payload)
		if err != nil {
			return nil, err
		}
		decompressed, err := io.ReadAll(r)
		if err != nil {
			return nil, err
		}
		return bytes.NewReader(decompressed), nil
	default:
		return nil, fmt.Errorf("unsupported compression type %q (supported compressions: zstd): %w", compression, ErrUnsupportedCompression)
	}
}
