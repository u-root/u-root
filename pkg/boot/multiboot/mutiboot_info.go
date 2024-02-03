// Copyright 2020 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package multiboot

import (
	"github.com/u-root/uio/uio"
)

type esxBootInfoInfo struct {
	cmdline uint64

	elems []elem
}

func (m *esxBootInfoInfo) marshal() []byte {
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
	typ() esxBootInfoType
	marshal() []byte
}

type esxBootInfoType uint32

const (
	ESXBOOTINFO_INVALID_TYPE  esxBootInfoType = 0
	ESXBOOTINFO_MEMRANGE_TYPE esxBootInfoType = 1
	ESXBOOTINFO_MODULE_TYPE   esxBootInfoType = 2
	ESXBOOTINFO_VBE_TYPE      esxBootInfoType = 3
	ESXBOOTINFO_EFI_TYPE      esxBootInfoType = 4
	ESXBOOTINFO_LOADESX_TYPE  esxBootInfoType = 5
)

type esxBootInfoMemRange struct {
	startAddr uint64
	length    uint64
	memType   uint32
}

func (m esxBootInfoMemRange) typ() esxBootInfoType {
	return ESXBOOTINFO_MEMRANGE_TYPE
}

func (m *esxBootInfoMemRange) marshal() []byte {
	buf := uio.NewNativeEndianBuffer(nil)
	buf.Write64(m.startAddr)
	buf.Write64(m.length)
	buf.Write32(m.memType)
	return buf.Data()
}

type esxBootInfoModuleRange struct {
	startPageNum uint64
	numPages     uint32
}

type esxBootInfoModule struct {
	cmdline    uint64
	moduleSize uint64
	ranges     []esxBootInfoModuleRange
}

func (m esxBootInfoModule) typ() esxBootInfoType {
	return ESXBOOTINFO_MODULE_TYPE
}

func (m *esxBootInfoModule) marshal() []byte {
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

type esxBootInfoEfiFlags uint32

const (
	// 64-bit ARM EFI. (Why would we have to tell the next kernel that it's
	// an aarch64 EFI? Shouldn't it know?)
	ESXBOOTINFO_EFI_ARCH64 esxBootInfoEfiFlags = 1 << 0

	// EFI Secure Boot in progress.
	ESXBOOTINFO_EFI_SECURE_BOOT esxBootInfoEfiFlags = 1 << 1

	// UEFI memory map is valid rather than esxBootInfo memory map.
	ESXBOOTINFO_EFI_MMAP esxBootInfoEfiFlags = 1 << 2
)

type esxBootInfoEfi struct {
	flags  esxBootInfoEfiFlags
	systab uint64

	// Only set if flags & ESXBOOTINFO_EFI_MMAP.
	memmap         uint64
	memmapNumDescs uint32
	memmapDescSize uint32
	memmapVersion  uint32
}

func (m esxBootInfoEfi) typ() esxBootInfoType {
	return ESXBOOTINFO_EFI_TYPE
}

func (m *esxBootInfoEfi) marshal() []byte {
	buf := uio.NewNativeEndianBuffer(nil)
	buf.Write32(uint32(m.flags))
	buf.Write64(m.systab)
	buf.Write64(m.memmap)
	buf.Write32(m.memmapNumDescs)
	buf.Write32(m.memmapDescSize)
	buf.Write32(m.memmapVersion)
	return buf.Data()
}
