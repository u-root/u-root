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

// String implements Stringer.
func (c *Control) String() string {
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

// String implements Stringer.
func (c *Status) String() string {
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

// String implements Stringer.
func (bar *BAR) String() string {
	// This little test lets us create empty strings, which
	// the JSON marshaler can then omit. That way, non-bridge
	// PCI devices won't even have this stuff show up.
	if bar.Base == 0 {
		return ""
	}
	var typ string
	base := bar.Base
	switch bar.Attr & 0xf {
	case 0:
		typ = "Memory at %08x (32-bit, non-prefetchable) [size=%#x]"
	case 1:
		typ = "I/O ports at %04x [size=%d]"
	case 2:
		typ = "Memory at %08x (32-bit, low 1Mbyte, non-prefetchable) [size=%#x]"
		if base < 0x100000 && base >= 0xc0000 {
			typ = "Expansion ROM at %08x (low 1Mbyte) [size=%#x]"
			if base&1 == 0 {
				typ = "(Disabled)" + typ
			} else {
				base--
			}
		}
	case 4:
		typ = "Memory at %08x (64-bit, non-prefetchable) [size=%#x]"
	case 8:
		typ = "Memory at %08x (32-bit, prefetchable) [size=%#x]"
	case 0xc:
		typ = "Memory at %08x (64-bit, prefetchable) [size=%#x]"
	default:
		return fmt.Sprintf("Can't get type from %#x", bar.Attr)
	}
	sz := bar.Lim - base + 1
	return fmt.Sprintf("Region %d: "+typ, bar.Index, base, sz)
}
