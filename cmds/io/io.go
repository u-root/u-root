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
	"io"
	"log"
	"os"
	"strconv"
)

const usage = `io [inb|inw|inl|rb|rw|rl|rq] address
io [outb|outw|outl|wb|ww|wl|wq] address value`

type IoIntf interface {
	Fd() uintptr
	io.Reader
	io.Seeker
	io.Writer
}

type Accessor func(IoIntf, uint64, interface{}) error

type iod struct {
	nargs    int
	addrbits int // not all addresses are multiples of 8 in size.
	val      interface{}
	valbits  int // not all value bits are multiples of 8 in size.
	format   string
	dev      string
	mode     int
	access   Accessor
}

var (
	ios = map[string]iod{
		"inb":  {2, 16, &b, 8, "%#02x", "/dev/port", os.O_RDONLY, inp},
		"inw":  {2, 16, &w, 16, "%#04x", "/dev/port", os.O_RDONLY, inp},
		"inl":  {2, 16, &l, 32, "%#08x", "/dev/port", os.O_RDONLY, inp},
		"outb": {3, 16, b, 8, "", "/dev/port", os.O_WRONLY, outp},
		"outw": {3, 16, w, 16, "", "/dev/port", os.O_WRONLY, outp},
		"outl": {3, 16, l, 32, "", "/dev/port", os.O_WRONLY, outp},
		"rb":   {2, 64, &b, 8, "%#02x", "/dev/mem", os.O_RDONLY, in},
		"rw":   {2, 64, &w, 16, "%#04x", "/dev/mem", os.O_RDONLY, in},
		"rl":   {2, 64, &l, 32, "%#08x", "/dev/mem", os.O_RDONLY, in},
		"rq":   {2, 64, &q, 64, "%#16x", "/dev/mem", os.O_RDONLY, in},
		"wb":   {3, 64, b, 8, "", "/dev/mem", os.O_WRONLY, out},
		"ww":   {3, 64, w, 16, "", "/dev/mem", os.O_WRONLY, out},
		"wl":   {3, 64, l, 32, "", "/dev/mem", os.O_WRONLY, out},
		"wq":   {3, 64, q, 64, "", "/dev/mem", os.O_WRONLY, out},
	}
	b    byte
	w    uint16
	l    uint32
	q    uint64
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

	f, err := os.OpenFile(i.dev, i.mode, 0)
	if err != nil {
		log.Fatalf("%v", err)
	}

	addr, err := strconv.ParseUint(a[1], 0, i.addrbits)
	if err != nil {
		log.Fatalf("Parsing address for %d bits: %v %v", i.addrbits, a[1], err)
	}

	switch a[0][0] {
	case 'i', 'r':
		err = i.access(f, addr, i.val)
	case 'o', 'w':
		var v uint64
		v, err = strconv.ParseUint(a[2], 0, i.valbits)
		if err != nil {
			log.Fatalf("%v: %v", a, err)
		}
		switch i.valbits {
		case 8:
			err = i.access(f, addr, uint8(v))
		case 16:
			err = i.access(f, addr, uint16(v))
		case 32:
			err = i.access(f, addr, uint32(v))
		case 64:
			err = i.access(f, addr, uint64(v))
		}
	default:
		log.Fatalf(usage)
	}

	if err != nil {
		log.Fatalf("%v: %v", a, err)
	}

	if i.format != "" {
		switch i.val.(type) {
		case *uint8:
			fmt.Printf(i.format+"\n", *i.val.(*uint8))
		case *uint16:
			fmt.Printf(i.format+"\n", *i.val.(*uint16))
		case *uint32:
			fmt.Printf(i.format+"\n", *i.val.(*uint32))
		case *uint64:
			fmt.Printf(i.format+"\n", *i.val.(*uint64))
		}

	}
}
