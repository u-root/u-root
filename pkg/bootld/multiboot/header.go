// Copyright 2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Multiboot header as defined in
// https://www.gnu.org/software/grub/manual/multiboot/multiboot.html#Header-layout
package multiboot

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
	"log"

	"github.com/u-root/u-root/pkg/ubinary"
)

var ErrHeaderNotFound = errors.New("multiboot header not found")
var ErrFlagsNotSupported = errors.New("multiboot header flags not supported yet")

const headerMagic = 0x1BADB002

type headerFlag uint32

const (
	flagHeaderPageAlign  headerFlag = 0x00000001
	flagHeaderMemoryInfo            = 0x00000002

	// ignore this flag for now
	flagHeaderMultibootVideoMode = 0x00000004

	flagHeaderUnsupported = 0x0000FFF8
)

// mandatory is a mandatory part of Multiboot v1 header.
type mandatory struct {
	Magic    uint32
	Flags    headerFlag
	Checksum uint32
}

// optional is an optional part of Multiboot v1 header.
type optional struct {
	HeaderAddr  uint32
	LoadAddr    uint32
	LoadEndAddr uint32
	BSSEndAddr  uint32
	EntryAddr   uint32

	ModeType uint32
	Width    uint32
	Height   uint32
	Depth    uint32
}

// header represents a Multiboot v1 header loaded from the file.
type header struct {
	mandatory
	optional
}

// parseHeader parses multiboot header as defined in
// https://www.gnu.org/software/grub/manual/multiboot/multiboot.html#OS-image-format
func parseHeader(r io.Reader) (header, error) {
	mandatorySize := binary.Size(mandatory{})
	optionalSize := binary.Size(optional{})
	sizeofHeader := mandatorySize + optionalSize
	var hdr header
	// The multiboot header must be contained completely within the
	// first 8192 bytes of the OS image.
	buf := make([]byte, 8192)
	n, err := io.ReadAtLeast(r, buf, mandatorySize)
	if err != nil {
		return hdr, err
	}
	buf = buf[:n]

	// Append zero bytes to the end of buffer to be able to read hdr
	// in a single binary.Read() when the mandatory
	// part of the header starts near the 8192 boundary.
	buf = append(buf, make([]byte, optionalSize)...)
	br := new(bytes.Reader)
	for len(buf) >= sizeofHeader {
		br.Reset(buf)
		if err := binary.Read(br, ubinary.NativeEndian, &hdr); err != nil {
			return hdr, err
		}
		if hdr.Magic == headerMagic && (hdr.Magic+uint32(hdr.Flags)+hdr.Checksum) == 0 {
			if hdr.Flags&flagHeaderUnsupported != 0 {
				return hdr, ErrFlagsNotSupported
			}
			if hdr.Flags&flagHeaderMultibootVideoMode != 0 {
				log.Print("VideoMode flag is not supproted yet, trying to load anyway")
			}
			return hdr, nil
		}
		// The Multiboot header must be 32-bit aligned.
		buf = buf[4:]
	}
	return hdr, ErrHeaderNotFound
}
