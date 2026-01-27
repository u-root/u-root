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
)

var ErrImageTooSmall = errors.New("image too small")
var ErrMagicMismatch = errors.New("magic number mismatch")
var ErrNotZImage = errors.New("not a zimg")
var ErrUnsupportedCompression = errors.New("unsupported compression type")

func Extract(img io.ReaderAt) (io.ReaderAt, error) {
	header := Header{}
	if err := binary.Read(io.NewSectionReader(img, 0, 28), binary.LittleEndian, &header); err != nil {
		return nil, err
	}

	if header.Magic != magic {
		return nil, fmt.Errorf("found magic %#x but PEZ expects %#x: %w", header.Magic, magic, ErrMagicMismatch)
	}
	if header.Type != typeZImage {
		return nil, fmt.Errorf("found type %#x but PEZ expects %#x: %w", header.Type, typeZImage, ErrNotZImage)
	}

	payload := io.NewSectionReader(img, int64(header.Offset), int64(header.Size))
	compression := string(header.Compression[:])

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
