// Copyright 2012-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// io lets you do IO operations.
//
// Synopsis:
//     io [inb|inw|inl] address
//     io [outb|outw|outl] address value
//
// Description:
//     io will let you do IO instructions on various architectures that support it.
//
package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
)

const usage = `io [inb|inw|inl] address
io [outb|outw|outl] address value`

type iod struct {
	nargs    int
	addrbits int // not all addresses are multiples of 8 in size.
	val      interface{}
	valbits  int // not all value bits are multiples of 8 in size.
	format   string
}

var (
	ios = map[string]iod{
		"inb":  {2, 16, &b, 8, "%#02x"},
		"inw":  {2, 16, &w, 16, "%#04x"},
		"inl":  {2, 16, &l, 32, "%#08x"},
		"outb": {3, 16, b, 8, ""},
		"outw": {3, 16, w, 16, ""},
		"outl": {3, 16, l, 32, ""},
	}
	b    byte
	w    uint16
	l    uint32
	addr uint64
)

func main() {
	var err error
	a := os.Args[1:]

	if len(a) == 0 {
		log.Fatal(usage)
	}

	i, ok := ios[a[0]]
	if !ok || len(a) != i.nargs {
		log.Fatal(usage)
	}

	addr, err := strconv.ParseUint(a[1], 0, i.addrbits)
	if err != nil {
		log.Fatalf("Parsing address for %d bits: %v %v", i.addrbits, a[1], err)
	}

	switch a[0] {
	case "inb", "inw", "inl":
		err = in(addr, i.val)
	case "outb", "outw", "outl":
		var v uint64
		v, err = strconv.ParseUint(a[2], 0, i.valbits)
		if err != nil {
			log.Fatalf("%v: %v", a, err)
		}
		switch t := i.val.(type) {
		case uint8:
			t = uint8(v)
		case uint16:
			t = uint16(v)
		case uint32:
			t = uint32(v)
		default:
			log.Fatalf("Can't handle %T for %v command", t, a[0])
		}
		err = out(addr, i.val)
	default:
		log.Fatalf(usage)
	}

	if err != nil {
		log.Fatalf("%v: %v", a, err)
	}

	if i.format != "" {
		switch i.val.(type) {
		case *uint8:
			fmt.Printf(i.format, *i.val.(*uint8))
		case *uint16:
			fmt.Printf(i.format, *i.val.(*uint16))
		case *uint32:
			fmt.Printf(i.format, *i.val.(*uint32))
		}

	}
}
