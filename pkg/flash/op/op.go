// Copyright 2021 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package op contains available SPI opcodes. The opcode is typically sent at
// the beginning of a SPI transaction.
package op

import (
	"fmt"
	"strings"
)

type OpCode byte

const (
	// PageProgram programs a page on the flash chip.
	PageProgram OpCode = 0x02
	// Read reads from the flash chip.
	Read OpCode = 0x03
	// WriteDisable disables writing.
	WriteDisable OpCode = 0x04
	// ReadStatus reads the status register.
	ReadStatus OpCode = 0x05
	// WriteEnable enables writing.
	WriteEnable OpCode = 0x06
	// SectorErase erases a sector to the value 0xff.
	SectorErase OpCode = 0x20
	// ReadSFDP reads from the SFDP.
	ReadSFDP OpCode = 0x5a
	// ReadID reads the JEDEC ID.
	ReadJEDECID OpCode = 0x9f
	// PRD/RES
	PRDRES OpCode = 0xab
	// AAI is auto address increment
	AAI OpCode = 0xad
	// Enter4BA enters 4-OpCode addressing mode.
	Enter4BA OpCode = 0xb7
	// BlockErase erases a block to the value 0xff.
	BlockErase OpCode = 0xd8
	// Exit4BA exits 4-OpCode addressing mode.
	Exit4BA OpCode = 0xe9
)

func (o OpCode) String() string {
	switch o {
	case PageProgram:
		return "PageProgram"
	case Read:
		return "Read"
	case WriteDisable:
		return "WriteDisable"
	case ReadStatus:
		return "ReadStatus"
	case WriteEnable:
		return "WriteEnable"
	case SectorErase:
		return "SectorErase"
	case ReadSFDP:
		return "ReadSFDP"
	case ReadJEDECID:
		return "ReadJEDECID"
	case PRDRES:
		return "PRDRES"
	case AAI:
		return "AAI"
	case Enter4BA:
		return "Enter4BA"
	case BlockErase:
		return "BlockErase"
	case Exit4BA:
		return "Exit4BA"
	default:
		return fmt.Sprintf("Unknown(%02x)", byte(o))
	}
}

func (o OpCode) Bytes() []byte {
	return []byte{byte(o)}
}

type Status byte

// Status is not universally defined, but a few bits are common.
const (
	WriteBusy Status = 1 << iota
	WriteEnabled
	ByteProtect0
	ByteProtect1
	ByteProtect2
	ByteProtectP3
	AutoAddressIncrement
	ByteProtectLocked
)

var names = []string{
	"WriteBusy",
	"WriteEnabled",
	"ByteProtect0",
	"ByteProtect1",
	"ByteProtect2",
	"ByteProtectP3",
	"AutoAddressIncrement",
	"ByteProtectLocked",
}

func (status Status) String() string {
	var s string
	for i := range 8 {
		if byte(status)&(1<<i) != 0 {
			s = s + names[i] + "|"
		}
	}
	s = strings.TrimRight(s, "|")
	return s
}

func (status Status) Busy() bool {
	return (status & WriteBusy) != 0
}
