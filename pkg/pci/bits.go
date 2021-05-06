// Copyright 2021 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pci

import (
	"fmt"
	"strconv"
	"strings"
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
	// Gaul was divided into three parts.
	// So are the BARs.
	b := strings.Fields(string(*bar))
	// If the bar is not something that matches the known format,
	// your kernel is broken. Just return something to be printed.
	if len(b) != 3 {
		return fmt.Sprintf("Could not parse %q", string(*bar))
	}
	// The type is the last byte of the 3rd field.
	// Kind of wish there were a substring operator that
	// Did All The Right Things.
	t := b[2][len(b[2])-1:]
	var typ string
	switch t {
	case "0":
		typ = "Memory at %08x (32-bit, non-prefetchable) [size=%#x]"
	case "1":
		typ = "I/O ports at %04x [size=%d]"
	case "4":
		typ = "Memory at %08x (32-bit, non-prefetchable) [size=%#x]"
	case "8":
		typ = "Memory at %08x (32-bit, prefetchable) [size=%#x]"
	case "c":
		typ = "Memory at %016x (64-bit, prefetchable) [size=%#x]"
	default:
		return fmt.Sprintf("Can't get type from %q", string(*bar))
	}
	base, err := strconv.ParseUint(b[0], 0, 0)
	if err != nil {
		return fmt.Sprintf("Could not parse %q", string(*bar))
		base = 0
	}
	end, err := strconv.ParseUint(b[1], 0, 0)
	if err != nil {
		return fmt.Sprintf("Could not parse %q", string(*bar))
		end = 0
	}
	sz := end - base + 1
	return fmt.Sprintf(typ, base, sz)
}
