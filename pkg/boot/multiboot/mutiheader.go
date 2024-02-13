// Copyright 2020 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package multiboot

import (
	"bytes"
	"encoding/binary"
	"io"
	"log"
)

// esxBootInfoMagic is both the magic value found in the esxBootInfo kernel header as
// well as the value the loaded OS expects in EAX at boot time.
const esxBootInfoMagic = 0x1BADB005

type esxBootInfoHeaderFlag uint32

const (
	// Kernel runs in EL1, not EL2.
	ESXBOOTINFO_ARCH_FLAG_ARM64_EL1 esxBootInfoHeaderFlag = 1 << 0
	// Must pass video info to OS.
	ESXBOOTINFO_FLAG_VIDEO esxBootInfoHeaderFlag = 1 << 2

	// rts_vaddr field is valid.
	ESXBOOTINFO_FLAG_EFI_RTS_OLD esxBootInfoHeaderFlag = 1 << 17
	// rts vaddr and size fields valid.
	ESXBOOTINFO_FLAG_EFI_RTS_NEW esxBootInfoHeaderFlag = 1 << 18
	// LoadESX version field valid.
	ESXBOOTINFO_FLAG_LOADESX_VERSION esxBootInfoHeaderFlag = 1 << 19
	// Video min fields valid.
	ESXBOOTINFO_FLAG_VIDEO_MIN esxBootInfoHeaderFlag = 1 << 20
)

type esxBootInfoVideoMode uint32

const (
	ESXBOOTINFO_VIDEO_GRAPHIC = 0
	ESXBOOTINFO_VIDEO_TEXT    = 1
)

type esxBootInfoHeader struct {
	Magic    uint32
	Flags    esxBootInfoHeaderFlag
	Checksum uint32

	// unused.
	_ uint32
	_ uint32

	// video stuff
	MinWidth  uint32
	MinHeight uint32
	MinDepth  uint32
	ModeType  esxBootInfoVideoMode
	Width     uint32
	Height    uint32
	Depth     uint32

	RuntimeServicesVAddr uint64
	RuntimeServicesSize  uint64
	LoadESXVersion       uint32
}

func (m *esxBootInfoHeader) name() string {
	return "ESXBootInfo"
}

func (m *esxBootInfoHeader) bootMagic() uintptr {
	return esxBootInfoMagic
}

// parseMutiHeader parses esxBootInfo header.
func parseMutiHeader(r io.Reader) (*esxBootInfoHeader, error) {
	sizeofHeader := binary.Size(esxBootInfoHeader{})

	var hdr esxBootInfoHeader
	// The multiboot header must be contained completely within the
	// first 8192 bytes of the OS image.
	buf := make([]byte, 8192)
	n, err := io.ReadAtLeast(r, buf, sizeofHeader)
	if err != nil {
		return nil, err
	}
	buf = buf[:n]

	br := new(bytes.Reader)
	for len(buf) >= sizeofHeader {
		br.Reset(buf)
		if err := binary.Read(br, binary.NativeEndian, &hdr); err != nil {
			return nil, err
		}
		if hdr.Magic == esxBootInfoMagic && (hdr.Magic+uint32(hdr.Flags)+hdr.Checksum) == 0 {
			/*if hdr.Flags&flagHeaderUnsupported != 0 {
				return hdr, ErrFlagsNotSupported
			}*/
			if hdr.Flags&ESXBOOTINFO_FLAG_VIDEO != 0 {
				log.Print("VideoMode flag is not supproted yet, trying to load anyway")
			}
			return &hdr, nil
		}
		// The Multiboot header must be 64-bit aligned.
		buf = buf[8:]
	}
	return nil, ErrHeaderNotFound
}
