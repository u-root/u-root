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
	"github.com/u-root/u-root/pkg/cmos"
	"github.com/u-root/u-root/pkg/memio"
)

var usageMsg = `usage: cmos read <index>...
cmos write <index> <data>...
`

func usage() {
	fmt.Print(usageMsg)
	os.Exit(1)
}

func processRegStr(regStr string) uint64 {
	reg, err := strconv.ParseUint(regStr, 10, 7)
	if err != nil || reg < 14 {
		log.Fatal("accessible bytes are only between 14-127")
	}
	return reg
}

func main() {
	switch os.Args[1] {
	case "read":
		if len(os.Args) != 3 {
			usage()
		}
		val := memio.Uint8(0)
		data := &val
		reg := processRegStr(os.Args[2])
		if err := cmos.Read(reg, data); err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%s\n", data)
	case "write":
		if len(os.Args) != 4 {
			usage()
		}
		reg := processRegStr(os.Args[2])
		value, err := strconv.ParseUint(os.Args[3], 10, 8)
		if err != nil {
			log.Fatal(err)
		}

		val := memio.Uint8(int8(value))
		data := &val
		if err := cmos.Write(reg, data); err != nil {
			log.Fatal(err)
		}
	default:
		usage()
	}
}
