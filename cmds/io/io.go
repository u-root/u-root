// Copyright 2010-2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// io reads and writes to physical memory and ports.
//
// Synopsis:
//     io r{b,w,l,q} address
//     io w{b,w,l,q} address value
//     # x86 only:
//     io in{b,w,l} address
//     io out{b,w,l} address value
//
// Description:
//     io lets you read/write 1/2/4/8-bytes to memory with the {r,w}{b,w,l,q}
//     commands respectively.
//
//     On x86 platforms, {in,out}{b,w,l} allow for port io.
//
// Examples:
//     # Read 8-bytes from address 0x10000
//     io rq 0x10000
//     # Write to the serial port on x86
//     io outb 0x3f8 50
package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/u-root/u-root/pkg/io"
)

type cmd struct {
	f                 func(addr int64, data io.UintN) error
	addrBits, valBits int
}

var (
	readCmds = map[string]cmd{
		"rb": {io.Read, 64, 8},
		"rw": {io.Read, 64, 16},
		"rl": {io.Read, 64, 32},
		"rq": {io.Read, 64, 64},
	}
	writeCmds = map[string]cmd{
		"wb": {io.Write, 64, 8},
		"ww": {io.Write, 64, 16},
		"wl": {io.Write, 64, 32},
		"wq": {io.Write, 64, 64},
	}
)

var usageMsg = `io r{b,w,l,q} address
io w{b,w,l,q} address value
`

func usage() {
	fmt.Print(usageMsg)
	os.Exit(1)
}

// newInt constructs a UintN with the specified value and bits.
func newInt(val uint64, bits int) io.UintN {
	switch bits {
	case 8:
		val := io.Uint8(int8(val))
		return &val
	case 16:
		val := io.Uint16(uint16(val))
		return &val
	case 32:
		val := io.Uint32(uint32(val))
		return &val
	case 64:
		val := io.Uint64(uint64(val))
		return &val
	default:
		panic(fmt.Sprintf("invalid number of bits %d", bits))
	}
}

func main() {
	if len(os.Args) < 3 {
		usage()
	}

	if c, ok := readCmds[os.Args[1]]; ok {
		if len(os.Args) != 3 {
			usage()
		}
		addr, err := strconv.ParseUint(os.Args[2], 0, c.addrBits)
		if err != nil {
			log.Fatal(err)
		}
		data := newInt(0, c.valBits)
		if err := c.f(int64(addr), data); err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%s\n", data)
	} else if c, ok := writeCmds[os.Args[1]]; ok {
		if len(os.Args) != 4 {
			usage()
		}
		addr, err := strconv.ParseUint(os.Args[2], 0, c.addrBits)
		if err != nil {
			log.Fatal(err)
		}
		value, err := strconv.ParseUint(os.Args[3], 0, c.valBits)
		if err != nil {
			log.Fatal(err)
		}
		data := newInt(value, c.valBits)
		if err := c.f(int64(addr), data); err != nil {
			log.Fatal(err)
		}
	} else {
		usage()
	}
}
