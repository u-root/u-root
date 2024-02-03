// Copyright 2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package chips contains chips known to work with the flash tool.
package chips

import (
	"fmt"
	"os"

	"github.com/u-root/u-root/pkg/flash/op"
)

const (
	k = 1024
	m = 1024 * 1024
)

// EraseBlock defines the size and the opcode used for an erase block.
// In Flash Chips, the erase region is aligned to the size.
type EraseBlock struct {
	Size int
	Op   op.OpCode
}

type ID int

// Chip defines a chip.
type Chip struct {
	Vendor      string
	Chip        string
	ID          ID
	ArraySize   int64
	PageSize    int64
	SectorSize  int64
	BlockSize   int64
	Is4BA       bool
	EraseBlocks []EraseBlock

	WriteEnableInstructionRequired bool
	WriteEnableOpcodeSelect        op.OpCode
	Write                          op.OpCode
	Read                           op.OpCode
}

// String implements string for Chip.
func (c *Chip) String() string {
	return fmt.Sprintf("Vendor:%s Chip:%s ID:%06x Size:%d 4BA:%v", c.Vendor, c.Chip, c.ID, c.ArraySize, c.Is4BA)
}

// Lookup finds a Chip by id, returning os.ErrNotExist if it is not found.
func Lookup(id ID) (*Chip, error) {
	for _, c := range Chips {
		if c.ID == id {
			return &c, nil
		}
	}
	return nil, os.ErrNotExist
}

// Chips are all the chips we know about.
// Note that the test assumes that SST25VF016B
// is first.
var Chips = []Chip{
	{
		Vendor:    "SST",
		Chip:      "SST25VF016B",
		ID:        0xbf2541,
		ArraySize: 2 * m,
		// This is the real page size.
		// The kernel gets an error on the ioctl.
		// PageSize:   256 * 1024,
		PageSize:   1, // make it 1 until we get AAI 1024,
		SectorSize: 4 * k,
		BlockSize:  64 * k,
		Is4BA:      false,
		EraseBlocks: []EraseBlock{
			{
				Size: 4 * k,
				Op:   0x20,
			},
			{
				Size: 32 * k,
				Op:   0x52,
			},
			{
				Size: 64 * k,
				Op:   0xD8,
			},
			{
				Size: 2 * m,
				Op:   0x60,
			},
			{
				Size: 2 * m,
				Op:   0xc7,
			},
		},

		WriteEnableInstructionRequired: true,
		WriteEnableOpcodeSelect:        op.WriteEnable,
		Write:                          op.AAI,
		Read:                           op.Read,
	},
}
