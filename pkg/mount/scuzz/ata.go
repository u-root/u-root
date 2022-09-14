// Copyright 2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package scuzz

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"strings"
)

// direction is the transfer direction.
type direction int32

// These are the acceptable constants for the sg_io_hdr dxfer_direction field
// defined by Linux.
const (
	_SG_DXFER_NONE     direction = -1
	_SG_DXFER_TO_DEV   direction = -2
	_SG_DXFER_FROM_DEV direction = -3
)

const (
	oldSchoolBlockLen = 512

	//	ataUsingLBA uint8 = (1 << 6)  nolint:golint,unused
	//	ataStatDRQ  uint8 = (1 << 3)  nolint:golint,unused
	//	ataStatErr  uint8 = (1 << 0)  nolint:golint,unused

	//	read  uint8 = 0  nolint:golint,unused
	// ataTo int32 = 1

	// ioPIO indicates we are doing programmed IO
	//	ioPIO = 0  nolint:golint,unused
	// ioDMA indicates we are doing DMA
	//	ioDMA = 1  nolint:golint,unused

	//	opCheckCondition = 0x02  nolint:golint,unused
	//	opDriverSense    = 0x08  nolint:golint,unused

	ata16    = 0x85
	ata16Len = 16
	cdbSize  = ata16Len

	maxStatusBlockLen = 32

	lba48    = 1
	nonData  = (3 << 1)
	pioIn    = (4 << 1)
	pioOut   = (5 << 1)
	protoDMA = (6 << 1)

	//	tlenNoData = 0 << 0  nolint:golint,unused
	//	tlenFeat   = 1 << 0  nolint:golint,unused
	tlenNsect = 2 << 0

	//	tlenBytes   = 0 << 2  nolint:golint,unused
	tlenSectors = 1 << 2

	tdirTo    = 0 << 3
	tdirFrom  = 1 << 3
	checkCond = 1 << 5
)

type (
	// Cmd is an ATA command. See the ATA standard starting in the 80s.
	Cmd uint8

	// dataBlock is a classic 512-byte ATA block.
	// It is not completely clear we need to do SG operations
	// in units of a block, but libraries and commands
	// seem to think so, and we're not sure we can verify
	// for all kernels what the rules are.
	dataBlock [oldSchoolBlockLen]byte

	// Some ATA blocks are best dealt with as blocks of "words".
	// They are in big-endian order.
	wordBlock [oldSchoolBlockLen / 2]uint16

	// commandDataBlock defines a command and any associated data.
	commandDataBlock [ata16Len]byte

	// statusBlock is the status returned from a drive operation.
	statusBlock [maxStatusBlockLen]byte
)

func (b dataBlock) toWordBlock() (wordBlock, error) {
	var w wordBlock
	err := binary.Read(bytes.NewBuffer(b[:]), binary.BigEndian, &w)
	return w, err
}

// mustLBA confirms that we are dealing with an LBA device.
// This means a post-2003 device. The standard hdparm command deals
// with all kinds of obsolete stuff; we don't care.
// The spec does not have a name for the bits or the offsets
// other than "Obsolete", "Retired", "Must be zero" or "Must be one".
// We follow the practice of existing code of not naming them either.
func (w wordBlock) mustLBA() error {
	check := []struct {
		off  int
		mask uint16
		bit  uint16
	}{
		{0, 0x8000, 0x8000},
		{49, 0x200, 0x200},
		{83, 0xc000, 0x4000},
		{86, 0x400, 0x400},
	}
	for _, c := range check {
		v, m, b := w[c.off], c.mask, c.bit
		if (v & m) != b {
			return fmt.Errorf("unsupported and probably non-ATA48 ddevice: word at offset %d: %#x and should be %#x", c.off, c.mask&v, b)
		}
	}
	return nil
}

// ataString writes out an ATAString in decoded human-readable form.
//
// ATA Command Set 4, Section 3.4.9: "Each pair of bytes in an ATA string is
// swapped ...". That the space character is used as padding is not mandated in
// the spec, but is used in the ATA string example, and on every device we've
// tried.
func ataString(a []byte) string {
	var s strings.Builder
	for i := 0; i < len(a); i += 2 {
		s.WriteByte(a[i+1])
		s.WriteByte(a[i])
	}
	return strings.TrimSpace(string(s.String()))
}

// unpackIdentify unpacks a wordBlock into an Info.
func unpackIdentify(s statusBlock, d dataBlock, w wordBlock) *Info {
	var info Info
	info.NumberSectors = binary.LittleEndian.Uint64(d[200:208])

	info.ECCBytes = uint(w[22])

	info.OrigSerial = string(d[20:40])
	info.OrigFirmwareRevision = string(d[46:54])
	info.OrigModel = string(d[54:94])

	info.Serial = ataString(d[20:40])
	info.FirmwareRevision = ataString(d[46:54])
	info.Model = ataString(d[54:94])

	info.MasterRevision = binary.LittleEndian.Uint16(d[184:186])
	info.SecurityStatus = DiskSecurityStatus(binary.LittleEndian.Uint16(d[256:258]))

	info.TrustedComputingSupport = w[48]
	return &info
}
