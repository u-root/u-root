// Copyright 2021 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package op contains available SPI opcodes. The opcode is typically sent at
// the beginning of a SPI transaction.
package op

const (
	// PageProgram programs a page on the flash chip.
	PageProgram byte = 0x02
	// Read reads from the flash chip.
	Read byte = 0x03
	// WriteDisable disables writing.
	WriteDisable byte = 0x04
	// ReadStatus reads the status register.
	ReadStatus byte = 0x05
	// WriteEnable enables writing.
	WriteEnable byte = 0x06
	// SectorErase erases a sector to the value 0xff.
	SectorErase byte = 0x20
	// ReadSFDP reads from the SFDP.
	ReadSFDP byte = 0x5a
	// ReadID reads the JEDEC ID.
	ReadJEDECID byte = 0x9f
	// PRD/RES
	PRDRES = 0xab
	// Enter4BA enters 4-byte addressing mode.
	Enter4BA byte = 0xb7
	// BlockErase erases a block to the value 0xff.
	BlockErase byte = 0xd8
	// Exit4BA exits 4-byte addressing mode.
	Exit4BA byte = 0xe9
)
