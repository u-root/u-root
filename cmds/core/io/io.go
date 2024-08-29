// Copyright 2010-2020 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//go:build !tinygo || tinygo.enable

// io reads and writes to physical memory and ports.
//
// Synopsis:
//
//	io (r{b,w,l,q} address)...
//	io (w{b,w,l,q} address value)...
//	# x86 only:
//	io (in{b,w,l} address)
//	io (out{b,w,l} address value)
//	io (cr index}
//	io {cw index value}...
//
// Description:
//
//	io lets you read/write 1/2/4/8-bytes to memory with the {r,w}{b,w,l,q}
//	commands respectively.
//
//	On x86 platforms, {in,out}{b,w,l} allow for port io.
//
//	Use cr / cw to write to cmos registers
//
// Examples:
//
//	# Read 8-bytes from address 0x10000 and 0x10000
//	io rq 0x10000 rq 0x10008
//	# Write to the serial port on x86
//	io outb 0x3f8 50
package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/u-root/u-root/pkg/memio"
)

type (
	cmdFunc func(addr int64, data memio.UintN) error
	cmd     struct {
		f                 cmdFunc
		addrBits, valBits int
	}
)

var (
	readCmds  = map[string]*cmd{}
	writeCmds = map[string]*cmd{}
	usageMsg  string
)

func addCmd(cmds map[string]*cmd, n string, f *cmd) {
	if _, ok := cmds[n]; ok {
		log.Fatalf("Command %q is defined twice", n)
	}
	cmds[n] = f
}

func usage() {
	fmt.Print(usageMsg)
	os.Exit(1)
}

// newInt constructs a UintN with the specified value and bits.
func newInt(val uint64, bits int) memio.UintN {
	switch bits {
	case 8:
		val := memio.Uint8(int8(val))
		return &val
	case 16:
		val := memio.Uint16(uint16(val))
		return &val
	case 32:
		val := memio.Uint32(uint32(val))
		return &val
	case 64:
		val := memio.Uint64(uint64(val))
		return &val
	default:
		panic(fmt.Sprintf("invalid number of bits %d", bits))
	}
}

func main() {
	if len(os.Args) < 3 {
		usage()
	}
	os.Args = os.Args[1:]

	// To avoid the command list from being partially executed when the
	// args fail to parse, queue them up and run all at once at the end.
	queue := []func(){}

	for len(os.Args) > 0 {
		var cmdStr string
		cmdStr, os.Args = os.Args[0], os.Args[1:]
		if c, ok := readCmds[cmdStr]; ok {
			// Parse arguments.
			if len(os.Args) < 1 {
				usage()
			}
			var addrStr string
			addrStr, os.Args = os.Args[0], os.Args[1:]
			addr, err := strconv.ParseUint(addrStr, 0, c.addrBits)
			if err != nil {
				log.Fatal(err)
			}

			queue = append(queue, func() {
				// Read from addr and print.
				data := newInt(0, c.valBits)
				if err := c.f(int64(addr), data); err != nil {
					log.Fatal(err)
				}
				fmt.Printf("%s\n", data)
			})
		} else if c, ok := writeCmds[cmdStr]; ok {
			// Parse arguments.
			if len(os.Args) < 2 {
				usage()
			}
			var addrStr, dataStr string
			addrStr, dataStr, os.Args = os.Args[0], os.Args[1], os.Args[2:]
			addr, err := strconv.ParseUint(addrStr, 0, c.addrBits)
			if err != nil {
				log.Fatal(err)
			}
			value, err := strconv.ParseUint(dataStr, 0, c.valBits)
			if err != nil {
				log.Fatal(err)
			}

			queue = append(queue, func() {
				// Write data to addr.
				data := newInt(value, c.valBits)
				if err := c.f(int64(addr), data); err != nil {
					log.Fatal(err)
				}
			})
		} else {
			usage()
		}
	}

	// Run all commands.
	for _, c := range queue {
		c()
	}
}
