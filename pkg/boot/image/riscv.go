// Copyright 2026 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package image

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
)

const (
	// RISCVMagic2 is the magic2 field of a RISC-V Image header, "RSC\x05".
	//
	// The header also has a magic field, but it is deprecated as of header
	// version 0.2 and may be removed, so magic2 is what we check.
	// (Documentation/arch/riscv/boot-image-header.rst)
	RISCVMagic2 = 0x05435352
)

// ErrNoImageSize is returned when a RISC-V Image header carries no image
// size. Unlike arm64 there is no fallback to guess at: the RISC-V
// documentation makes the field mandatory.
var ErrNoImageSize = errors.New("image size is zero")

// RISCVHeader is the header of a RISC-V Image.
//
// It is laid out like the arm64 header but is not identical to it: magic is 64
// bits wide here, and version occupies what arm64 uses as reserved space.
// (Documentation/arch/riscv/boot-image-header.rst)
type RISCVHeader struct {
	Code0      uint32 `offset:"0x00"`
	Code1      uint32 `offset:"0x04"`
	TextOffset uint64 `offset:"0x08"`
	ImageSize  uint64 `offset:"0x10"`
	Flags      uint64 `offset:"0x18"`
	Version    uint32 `offset:"0x20"`
	Res1       uint32 `offset:"0x24"`
	Res2       uint64 `offset:"0x28"`
	Magic      uint64 `offset:"0x30"`
	Magic2     uint32 `offset:"0x38"`
	Res3       uint32 `offset:"0x3c"`
}

// MajorVersion returns the major version of the header format.
func (h RISCVHeader) MajorVersion() uint16 {
	return uint16(h.Version >> 16)
}

// MinorVersion returns the minor version of the header format.
func (h RISCVHeader) MinorVersion() uint16 {
	return uint16(h.Version & 0xffff)
}

// RISCVImage abstracts a RISC-V Image.
type RISCVImage struct {
	Header RISCVHeader
	Data   []byte
}

// ParseRISCVFromBytes parses a RISC-V Image from a byte slice.
func ParseRISCVFromBytes(data []byte) (*RISCVImage, error) {
	img := &RISCVImage{}

	if err := binary.Read(bytes.NewBuffer(data), binary.LittleEndian, &img.Header); err != nil {
		return img, fmt.Errorf("unmarshaling riscv header: %w", err)
	}

	if img.Header.Magic2 != RISCVMagic2 {
		return img, ErrBadMagic
	}

	// "Image size is mandatory for boot loader to load kernel image.
	// Booting will fail otherwise."
	// (Documentation/arch/riscv/boot-image-header.rst)
	//
	// Unlike arm64 there is no older-kernel fallback to guess at here.
	if img.Header.ImageSize == 0 {
		return img, ErrNoImageSize
	}

	// Bit 0 of flags is the kernel endianness, 1 if big endian.
	if img.Header.Flags&0x1 != 0 {
		return img, ErrBadEndianness
	}

	img.Data = data

	return img, nil
}
