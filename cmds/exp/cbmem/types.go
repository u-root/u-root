// Copyright 2016-2021 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build linux

package main

import "fmt"

type seg struct {
	off int64
	dat []byte
}

const (
	// TableSize is the size of the LB Table Header in bytes.
	TableSize = 48
)

// TS is an ID and timestamp
type TS struct {
	EntryID    uint32
	EntryStamp uint64
}

// TSHeader e is the header of a time stamp table.
type TSHeader struct {
	BaseTime    uint64
	MaxEntries  uint16
	TickFreqMHZ uint16
	NumEntries  uint32
}

// TimeStamps is a TSHeader and its time stamps.
type TimeStamps struct {
	TSHeader
	TS []TS
}

// Header is the common cbmem header.
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

// Record is the common CBMEM record header: a tag and a type.
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

// Pointer is a record containing just an address to another record.
type Pointer struct {
	Record
	Addr uint32
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

type timeStampTableEntry struct {
	Record
	TimeStampTable uint32
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

// MTRREntry is MTRR Entry record.
type MTRREntry struct {
	Record
	Index uint32
}

// BoardIDEntry is the ID for a board.
type BoardIDEntry struct {
	Record
	BoardID uint32
}

// MACEntry is a Mac Entry record.
type MACEntry struct {
	MACaddr []uint8
	pad     []uint8
}

// LBEntry is a set of MAC addresses.
type LBEntry struct {
	Record
	Count    uint32
	MACAddrs []MACEntry
}

// LBRAMEntry is a defintion of RAM (is this even used?)
type LBRAMEntry struct {
	Record
	RAMCode uint32
}

// SPIFlashEntry defines properties of a SPI part.
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

// CBMemEntry defines the Address and other properties of a CBMem area.
type CBMemEntry struct {
	Record
	Address   uint64
	EntrySize uint32
	ID        uint32
}

// CMOSTable is a CMOS table header, defining the length of the table.
// CMOS entries have a variable-length encoding.
type CMOSTable struct {
	Record
	HeaderLength uint32
}

// CMOSEntry defines a CMOS entry. This is little used any more.
type CMOSEntry struct {
	Record
	Bit      uint32
	Length   uint32
	Config   uint32
	ConfigID uint32
	Name     []uint8
}

// CMOSEnums defines a CMOS entry enumeration type.
type CMOSEnums struct {
	Record
	ConfigID uint32
	Value    uint32
	Text     []uint8
}

// CMOSDefaults define CMOS defaults.
type CMOSDefaults struct {
	Record
	NameLength uint32
	Name       []uint8
	DefaultSet []uint8
}

// CMOSChecksum defines a checksum over a range of CMOS
type CMOSChecksum struct {
	Record
	RangeStart uint32
	RangeEnd   uint32
	Location   uint32
	Type       uint32
}

// CBmem is the set of all CBmem records in a coreboot image.
type CBmem struct {
	Memory           *memoryEntry
	MemConsole       *memconsoleEntry
	Consoles         []string
	TimeStampsTable  Pointer
	TimeStamps       *TimeStamps
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
