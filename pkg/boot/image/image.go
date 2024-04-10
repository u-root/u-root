// Copyright 2022 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package image contains a parser for Arm64 Linux Image format. It
// assumes little endian arm.
package image

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"math"
)

const (
	// Magic values used in Image header.
	Magic = 0x644d5241
)

var (
	kernelImageSize = uint64(math.Pow(2, 24)) // 16MB, a guess value similar to that used in kexec-tools.

	errBadMagic      = errors.New("bad header magic")
	errBadEndianness = errors.New("invalid Image endianness, expected little")
)

// Arm64Header is header for Arm64 Image.
type Arm64Header struct {
	Code0      uint32 `offset:"0x00"`
	Code1      uint32 `offset:"0x04"`
	TextOffset uint64 `offset:"0x08"`
	ImageSize  uint64 `offset:"0x10"`
	Flags      uint64 `offset:"0x18"`
	Res2       uint64 `offset:"0x20"`
	Res3       uint64 `offset:"0x28"`
	Res4       uint64 `offset:"0x30"`
	Magic      uint32 `offset:"0x38"`
	Res5       uint32 `offset:"0x3c"`
}

// Image abstracts Arm64 Image.
type Image struct {
	Header Arm64Header
	Data   []byte
}

// ParseFromBytes parse an Image from bytes slice.
func ParseFromBytes(data []byte) (*Image, error) {
	img := &Image{}

	if err := binary.Read(bytes.NewBuffer(data), binary.LittleEndian, &img.Header); err != nil {
		return img, fmt.Errorf("unmarshaling arm64 header: %w", err)
	}

	if img.Header.Magic != Magic {
		return img, errBadMagic
	}

	if img.Header.ImageSize == 0 {
		/* For 3.16 and older kernels. */
		img.Header.TextOffset = 0x80000
		img.Header.ImageSize = kernelImageSize
	}

	// v3.17: https://www.kernel.org/doc/Documentation/arm64/booting.txt
	//
	// NOTE(10000TB): For now assumes and support little endian arm.
	// Error out if Image is not little endian.
	if int(img.Header.Flags&0x1) != 0 {
		return img, errBadEndianness
	}

	img.Data = data

	return img, nil
}
