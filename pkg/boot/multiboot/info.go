// Copyright 2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package multiboot

import (
	"bytes"
	"encoding/binary"

	"github.com/u-root/u-root/pkg/ubinary"
)

type flag uint32

const (
	flagInfoMemory flag = 1 << iota
	flagInfoBootDev
	flagInfoCmdLine
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

	CmdLine uint32

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

// marshal writes out the exact bytes of multiboot info
// expected by the kernel being loaded.
func (i *info) marshal() ([]byte, error) {
	var buf bytes.Buffer
	if err := binary.Write(&buf, ubinary.NativeEndian, i); err != nil {
		return nil, err
	}

	size := (buf.Len() + 3) &^ 3
	_, err := buf.Write(bytes.Repeat([]byte{0}, size-buf.Len()))
	return buf.Bytes(), err
}
