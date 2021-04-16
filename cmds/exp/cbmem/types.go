// Copyright 2016-2021 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build linux

package main

import "fmt"

type seg struct {
	off int64
	dat []byte
}

const (
	TableSize = 48
)

// TODO: we don't have a platform yet that creates these.
type TS struct {
	entry_id    uint32
	entry_stamp uint64
}

type TSTable struct {
	base_time     uint64
	max_entries   uint16
	tick_freq_mhz uint16
	num_entries   uint32
	entries       []TS
}

type Header struct {
	Signature    [4]uint8
	HeaderSz     uint32
	HeaderCSUM   uint32
	TableSz      uint32
	TableCSUM    uint32
	TableEntries uint32
}

func (h *Header) String() string {
	return fmt.Sprintf("Signature %s Header Size %#x HeaderCSUM %#08x Table Size %#x TableCSUM %#08x TableEntries %d",
		string(h.Signature[:4]), h.HeaderSz, h.HeaderCSUM, h.TableSz, h.TableCSUM, h.TableEntries)
}

type Record struct {
	Tag  uint32
	Size uint32
}

type memoryRange struct {
	Start uint64
	Size  uint64
	Mtype uint32
}
type memoryEntry struct {
	Record
	Maps []memoryRange
}
type hwrpbEntry struct {
	Record
	HwrPB uint64
}

type consoleEntry struct {
	Record
	Console uint32
}

type mainboardEntry struct {
	Record
	Vendor     string
	PartNumber string
}

type stringEntry struct {
	Record
	String []uint8
}

type timestampEntry struct {
	Record
	TimeStamp uint32
}
type serialEntry struct {
	Record
	Type     uint32
	BaseAddr uint32
	Baud     uint32
	RegWidth uint32
}

type memconsoleEntry struct {
	Record
	Address uint64
	CSize   uint32
	Cursor  uint32
	Data    string
}
type forwardEntry struct {
	Record
	Forward uint64
}
type framebufferEntry struct {
	Record
	PhysicalAddress  uint64
	XResolution      uint32
	YRresolution     uint32
	BytesPerLine     uint32
	BitsPerPixel     uint8
	RedMaskPos       uint8
	RedMaskSize      uint8
	GreenMaskPos     uint8
	GreenMaskSize    uint8
	BlueMaskPos      uint8
	BlueMaskSize     uint8
	ReservedMaskPos  uint8
	ReservedMaskSize uint8
}

// GPIO is General Purpose IO
type GPIO struct {
	Port     uint32
	Polarity uint32
	Value    uint32
	Name     []uint8
}

type gpioEntry struct {
	Record
	Count uint32
	GPIOs []GPIO
}

type rangeEntry struct {
	Record
	RangeStart uint64
	RangeSize  uint32
}

type cbmemEntry struct {
	Record
	CbmemAddr uint64
}

type MTRREntry struct {
	Record
	Index uint32
}

type BoardIDEntry struct {
	Record
	BoardID uint32
}
type MACEntry struct {
	MACaddr []uint8
	pad     []uint8
}

type LBEntry struct {
	Record
	Count    uint32
	MACAddrs []MACEntry
}

type LBRAMEntry struct {
	Record
	RAMCode uint32
}

type SPIFlashEntry struct {
	Record
	FlashSize  uint32
	SectorSize uint32
	EraseCmd   uint32
}
type bootMediaParamsEntry struct {
	Record
	FMAPOffset    uint64
	CBFSOffset    uint64
	CBFSSize      uint64
	BootMediaSize uint64
}
type CBMemEntry struct {
	Record
	Address   uint64
	EntrySize uint32
	ID        uint32
}
type CMOSTable struct {
	Record
	HeaderLength uint32
}
type CMOSEntry struct {
	Record
	Bit      uint32
	Length   uint32
	Config   uint32
	ConfigID uint32
	Name     []uint8
}
type CMOSEnums struct {
	Record
	ConfigID uint32
	Value    uint32
	Text     []uint8
}
type CMOSDefaults struct {
	Record
	NameLength uint32
	Name       []uint8
	DefaultSet []uint8
}
type CMOSChecksum struct {
	Record
	RangeStart uint32
	RangeEnd   uint32
	Location   uint32
	Type       uint32
}

type CBmem struct {
	Memory           *memoryEntry
	MemConsole       *memconsoleEntry
	Consoles         []string
	TimeStamps       []timestampEntry
	UART             []serialEntry
	MainBoard        mainboardEntry
	Hwrpb            hwrpbEntry
	CBMemory         []cbmemEntry
	BoardID          BoardIDEntry
	StringVars       map[string]string
	BootMediaParams  bootMediaParamsEntry
	VersionTimeStamp uint32
	Unknown          []uint32
	Ignored          []string
}
