// Copyright 2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package multiboot

import (
	"bytes"
	"encoding/binary"
)

var sizeofInfo = uint32(binary.Size(info{}))

type flag uint32

const (
	flagInfoMemory flag = 1 << iota
	flagInfoBootDev
	flagInfoCmdline
	flagInfoMods
	flagInfoAoutSyms
	flagInfoElfSHDR
	flagInfoMemMap
	flagInfoDriveInfo
	flagInfoConfigTable
	flagInfoBootLoaderName
	flagInfoAPMTable
	flagInfoVideoInfo
	flagInfoFrameBuffer
)

// info represents the Multiboot v1 info passed to the loaded kernel.
type info struct {
	Flags    flag
	MemLower uint32
	MemUpper uint32

	// BootDevice is not supported, always zero.
	BootDevice uint32

	Cmdline uint32

	ModsCount uint32
	ModsAddr  uint32

	// Syms is not supported, always zero array.
	Syms [4]uint32

	MmapLength uint32
	MmapAddr   uint32

	// Following fields except BootLoaderName are not suppoted yet,
	// the values are always set to zeros.

	DriversLength uint32
	DrivesrAddr   uint32

	ConfigTable uint32

	BootLoaderName uint32

	APMTable uint32

	VBEControlInfo  uint32
	VBEModeInfo     uint32
	VBEMode         uint16
	VBEInterfaceSeg uint16
	VBEInterfaceOff uint16
	VBEInterfaceLen uint16

	FramebufferAddr   uint16
	FramebufferPitch  uint16
	FramebufferWidth  uint32
	FramebufferHeight uint32
	FramebufferBPP    byte
	FramebufferType   byte
	ColorInfo         [6]byte
}

type infoWrapper struct {
	info

	Cmdline        string
	BootLoaderName string
}

// marshal writes out the exact bytes of multiboot info
// expected by the kernel being loaded.
func (iw *infoWrapper) marshal(base uintptr) ([]byte, error) {
	offset := sizeofInfo + uint32(base)
	iw.info.Cmdline = offset
	offset += uint32(len(iw.Cmdline)) + 1
	iw.info.BootLoaderName = offset
	iw.info.Flags |= flagInfoCmdline | flagInfoBootLoaderName

	buf := bytes.Buffer{}
	if err := binary.Write(&buf, binary.NativeEndian, iw.info); err != nil {
		return nil, err
	}

	for _, s := range []string{iw.Cmdline, iw.BootLoaderName} {
		if _, err := buf.WriteString(s); err != nil {
			return nil, err
		}
		if err := buf.WriteByte(0); err != nil {
			return nil, err
		}
	}

	size := (buf.Len() + 3) &^ 3
	_, err := buf.Write(bytes.Repeat([]byte{0}, size-buf.Len()))
	return buf.Bytes(), err
}

func (iw infoWrapper) size() (uint, error) {
	b, err := iw.marshal(0)
	return uint(len(b)), err
}
