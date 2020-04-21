// Copyright 2010-2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// reads and writes to valid CMOS registers
//
// Description:
//	cmos gives you an an easy way to read and write to CMOS registers
//	using port io on x86 platforms.
//
// Examples:
//	# Read register 15
//	cmos read 15
//	# Write to register 15
//	cmos write 15 1
package main

import (
	"strconv"
	"os"
	"log"
	"fmt"
	"github.com/u-root/u-root/pkg/memio"
)

var usageMsg = `usage: cmos read <index>...
cmos write <index> <data>...
`

func usage() {
	fmt.Print(usageMsg)
	os.Exit(1)
}

func read(reg uint64, data memio.UintN) error {
	regVal := memio.Uint8(reg)
	if err := memio.Out(0x70, &regVal); err != nil {
		return err
	}
	return memio.In(0x71, data)
}

func write(reg uint64, data memio.UintN) error {
	regVal := memio.Uint8(reg)
	if err := memio.Out(0x70, &regVal); err != nil {
		return err
	}
	return memio.Out(0x71, data)
}

func main() {
	if len(os.Args) == 3 && os.Args[1] == "read" {
		val, regStr := memio.Uint8(0), os.Args[2]
		data := &val
		reg, err := strconv.ParseUint(regStr, 10, 7)
		if err != nil {
			fmt.Println("bytes above 127 inaccessable")
			log.Fatal(err)
		}
		if reg < 14 {
			fmt.Println("can't read bytes below 14")
			os.Exit(1)
		}
		if err := read(reg, data); err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%s\n", data)
	} else if len(os.Args) == 4 && os.Args[1] == "write" {
		reg, err := strconv.ParseUint(os.Args[2], 10, 7)
		if err != nil {
			fmt.Println("bytes above 127 inaccessable")
			log.Fatal(err)
		}
		if reg < 14 {
			fmt.Println("can't write to bytes below 14")
			os.Exit(1)
		}
		value, err := strconv.ParseUint(os.Args[3], 10, 8)
		if err != nil {
			log.Fatal(err)
		}

		val := memio.Uint8(int8(value))
		data := &val
		if err := write(reg, data); err != nil {
			log.Fatal(err)
		}
	} else {
		usage()
	}
}
