// Copyright 2020 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package multiboot

import (
	"github.com/u-root/u-root/pkg/uio"
)

type mutibootInfo struct {
	cmdline uint64

	elems []elem
}

func (m *mutibootInfo) marshal() []byte {
	buf := uio.NewNativeEndianBuffer(nil)
	buf.Write64(m.cmdline)
	buf.Write64(uint64(len(m.elems)))

	// These elems are encoded as TLV.
	// Which is nice, because then we don't have to encode offsets and shit.
	for _, el := range m.elems {
		b := el.marshal()
		buf.Write32(uint32(el.typ()))
		// The element size is the total size including itself - typ +
		// size + data.
		buf.Write64(uint64(len(b)) + 8 + 4)
		buf.WriteData(b)
	}
	return buf.Data()
}

type elem interface {
	typ() mutibootType
	marshal() []byte
}

type mutibootType uint32

const (
	MUTIBOOT_INVALID_TYPE  mutibootType = 0
	MUTIBOOT_MEMRANGE_TYPE mutibootType = 1
	MUTIBOOT_MODULE_TYPE   mutibootType = 2
	MUTIBOOT_VBE_TYPE      mutibootType = 3
	MUTIBOOT_EFI_TYPE      mutibootType = 4
	MUTIBOOT_LOADESX_TYPE  mutibootType = 5
)

type mutibootMemRange struct {
	startAddr uint64
	length    uint64
	memType   uint32
}

func (m mutibootMemRange) typ() mutibootType {
	return MUTIBOOT_MEMRANGE_TYPE
}

func (m *mutibootMemRange) marshal() []byte {
	buf := uio.NewNativeEndianBuffer(nil)
	buf.Write64(m.startAddr)
	buf.Write64(m.length)
	buf.Write32(m.memType)
	return buf.Data()
}

type mutibootModuleRange struct {
	startPageNum uint64
	numPages     uint32
}

type mutibootModule struct {
	cmdline    uint64
	moduleSize uint64
	ranges     []mutibootModuleRange
}

func (m mutibootModule) typ() mutibootType {
	return MUTIBOOT_MODULE_TYPE
}

func (m *mutibootModule) marshal() []byte {
	buf := uio.NewNativeEndianBuffer(nil)
	buf.Write64(m.cmdline)
	buf.Write64(m.moduleSize)
	buf.Write32(uint32(len(m.ranges)))
	for _, r := range m.ranges {
		buf.Write64(r.startPageNum)
		buf.Write32(r.numPages)
		// Padding.
		buf.Write32(0)
	}
	return buf.Data()
}

type mutibootEfiFlags uint32

const (
	// 64-bit ARM EFI. (Why would we have to tell the next kernel that it's
	// an aarch64 EFI? Shouldn't it know?)
	MUTIBOOT_EFI_ARCH64 mutibootEfiFlags = 1 << 0

	// EFI Secure Boot in progress.
	MUTIBOOT_EFI_SECURE_BOOT mutibootEfiFlags = 1 << 1

	// UEFI memory map is valid rather than mutiboot memory map.
	MUTIBOOT_EFI_MMAP mutibootEfiFlags = 1 << 2
)

type mutibootEfi struct {
	flags  mutibootEfiFlags
	systab uint64

	// Only set if flags & MUTIBOOT_EFI_MMAP.
	memmap         uint64
	memmapNumDescs uint32
	memmapDescSize uint32
	memmapVersion  uint32
}

func (m mutibootEfi) typ() mutibootType {
	return MUTIBOOT_EFI_TYPE
}

func (m *mutibootEfi) marshal() []byte {
	buf := uio.NewNativeEndianBuffer(nil)
	buf.Write32(uint32(m.flags))
	buf.Write64(m.systab)
	buf.Write64(m.memmap)
	buf.Write32(m.memmapNumDescs)
	buf.Write32(m.memmapDescSize)
	buf.Write32(m.memmapVersion)
	return buf.Data()
}
