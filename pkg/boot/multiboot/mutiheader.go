// Copyright 2020 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package multiboot

import (
	"bytes"
	"encoding/binary"
	"io"
	"log"

	"github.com/u-root/u-root/pkg/ubinary"
)

// mutibootMagic is both the magic value found in the mutiboot kernel header as
// well as the value the loaded OS expects in EAX at boot time.
const mutibootMagic = 0x1BADB005

type mutibootHeaderFlag uint32

const (
	// Kernel runs in EL1, not EL2.
	MUTIBOOT_ARCH_FLAG_ARM64_EL1 mutibootHeaderFlag = 1 << 0
	// Must pass video info to OS.
	MUTIBOOT_FLAG_VIDEO mutibootHeaderFlag = 1 << 2

	// rts_vaddr field is valid.
	MUTIBOOT_FLAG_EFI_RTS_OLD mutibootHeaderFlag = 1 << 17
	// rts vaddr and size fields valid.
	MUTIBOOT_FLAG_EFI_RTS_NEW mutibootHeaderFlag = 1 << 18
	// LoadESX version field valid.
	MUTIBOOT_FLAG_LOADESX_VERSION mutibootHeaderFlag = 1 << 19
	// Video min fields valid.
	MUTIBOOT_FLAG_VIDEO_MIN mutibootHeaderFlag = 1 << 20
)

type mutibootVideoMode uint32

const (
	MUTIBOOT_VIDEO_GRAPHIC = 0
	MUTIBOOT_VIDEO_TEXT    = 1
)

type mutibootHeader struct {
	Magic    uint32
	Flags    mutibootHeaderFlag
	Checksum uint32

	// unused.
	_ uint32
	_ uint32

	// video stuff
	MinWidth  uint32
	MinHeight uint32
	MinDepth  uint32
	ModeType  mutibootVideoMode
	Width     uint32
	Height    uint32
	Depth     uint32

	RuntimeServicesVAddr uint64
	RuntimeServicesSize  uint64
	LoadESXVersion       uint32
}

func (m *mutibootHeader) name() string {
	return "mutiboot"
}

func (m *mutibootHeader) bootMagic() uintptr {
	return mutibootMagic
}

// parseMutiHeader parses mutiboot header.
func parseMutiHeader(r io.Reader) (*mutibootHeader, error) {
	sizeofHeader := binary.Size(mutibootHeader{})

	var hdr mutibootHeader
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
		if err := binary.Read(br, ubinary.NativeEndian, &hdr); err != nil {
			return nil, err
		}
		if hdr.Magic == mutibootMagic && (hdr.Magic+uint32(hdr.Flags)+hdr.Checksum) == 0 {
			/*if hdr.Flags&flagHeaderUnsupported != 0 {
				return hdr, ErrFlagsNotSupported
			}*/
			if hdr.Flags&MUTIBOOT_FLAG_VIDEO != 0 {
				log.Print("VideoMode flag is not supproted yet, trying to load anyway")
			}
			return &hdr, nil
		}
		// The Multiboot header must be 64-bit aligned.
		buf = buf[8:]
	}
	return nil, ErrHeaderNotFound
}
