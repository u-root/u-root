// Copyright 2021 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pci

import (
	"fmt"
)

// bit puller-aparters. There's a case to be made for the usual tables
// with widths and values and stuff but this way ends up being easier
// to read, surprisingly.

// Control register bits are packed, bit 0 to bit 11. Thank you, PCI people.
var crBits = []string{
	"I/O",
	"Memory",
	"DMA",
	"Special",
	"MemWINV",
	"VGASnoop",
	"ParErr",
	"Stepping",
	"SERR",
	"FastB2B",
	"DisInt",
}

func (c *PCIControl) String() string {
	var s string
	for i, n := range crBits {
		if len(s) > 0 {
			s = s + " "
		}
		s += n
		ix := (1<<i)&uint16(*c) != 0
		if ix {
			s += "+"
		} else {
			s += "-"
		}
	}
	return s
}

var stBits = []string{
	"Reserved",
	"Reserved",
	"INTx",
	"Cap",
	"66MHz",
	"UDF",
	"FastB2b",
	"ParErr",
	"DEVSEL",
	">Tabort",
	"<Tabort",
	"<MABORT",
	">SERR",
	"<PERR",
}

func (c *PCIStatus) String() string {
	var s string

	for i, n := range stBits {
		switch i {
		case 9: // the only multi-bit field
			spd := (uint16(*c) & 0x600) >> 9
			if len(s) > 0 {
				s = s + " "
			}
			s += fmt.Sprintf("DEVSEL=%s", []string{"fast", "medium", "slow", "reserved"}[spd])
		case 0:
		case 1:
		case 10:
			continue
		default:
			ix := (1<<i)&uint16(*c) != 0
			if len(s) > 0 {
				s = s + " "
			}
			s += n
			if ix {
				s += "+"
			} else {
				s += "-"
			}
		}
	}
	return s
}
