// Copyright 2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package scuzz

const (
	// golangci's requirements around unused constants
	// are inconsistent with how one writes user level drivers.
	// Generally, in a user level driver, one wants to lay out all
	// constants for all registers, unused or not; it makes
	// the code far easier to read and work with.
	// Anyway, I'll comment stuff out until used, since the directives
	// don't work and are ugly.
	// Sure, you can add a nolint tag, but as your
	//	bidi int8 = -4  nolint:golint,unused
	//from direction = -3
	to direction = -2
	//none direction = -1

	oldSchoolBlockLen = 512

	//	ataUsingLBA uint8 = (1 << 6)  nolint:golint,unused
	//	ataStatDRQ  uint8 = (1 << 3)  nolint:golint,unused
	//	ataStatErr  uint8 = (1 << 0)  nolint:golint,unused

	//	read  uint8 = 0  nolint:golint,unused
	//ataTo int32 = 1

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

	tdirTo = 0 << 3
	//	tdirFrom  = 1 << 3  nolint:golint,unused
	checkCond = 1 << 5
)

// Commands. We only export those we implement.
const (
	// securityUnlock unlocks the drive with a given 32-byte password
	securityUnlock = 0xf2
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

	// commandDataBlock defines a command and any associated data.
	commandDataBlock [ata16Len]byte

	// statusBlock is the status returned from a drive operation.
	statusBlock [maxStatusBlockLen]byte

	direction int32
)
