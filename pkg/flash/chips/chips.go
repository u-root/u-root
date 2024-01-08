// Copyright 2021 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package chip contains chips known to work with the flash tool.
package chips

import "github.com/u-root/u-root/pkg/flash/op"

const (
	k = 1024
	m = 1024 * 1024
)

type EraseBlock struct {
	size int
	op   uint8
}

type Chip struct {
	Vendor      string
	Chip        string
	ID          int
	Size        int64
	PageSize    int
	EraseBlocks []EraseBlock

	Unlock uint8
	Write  uint8
	Read   uint8
}

var Chips = []Chip{
	{
		Vendor:   "SST",
		Chip:     "SST25VF016B",
		ID:       0xbf2541,
		Size:     2048 * 1048576,
		PageSize: 256 * 1024,
		EraseBlocks: []EraseBlock{
			{
				size: 4 * k,
				op:   0x20,
			},
			{
				size: 32 * k,
				op:   0x52,
			},
			{
				size: 64 * k,
				op:   0xD8,
			},
			{
				size: 2 * m,
				op:   0x60,
			},
			{
				size: 2 * m,
				op:   0xc7,
			},
		},

		Unlock: op.WriteEnable,
		Write:  op.AAI,
		Read:   op.Read,
	},
}
