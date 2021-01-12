// Copyright 2021 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"log"
	"math"

	"github.com/u-root/u-root/pkg/boot"
	"github.com/u-root/u-root/pkg/boot/uefi"
)

var (
	debug       = flag.Bool("d", false, "Print debug output")
	imageBase   = flag.Uint64("i", 0x800000, "Where to load payload image")
	serialAddr  = flag.Uint("serial_addr", 0x3f8, "Serial IO port address")
	serialWidth = flag.Uint("serial_width", 1, "Serial port reg width")
	serialHertz = flag.Uint("serial_hertz", 1843200, "Serial port input hertz")
	serialBaud  = flag.Uint("serial_baud", 115200, "Serial port baud rate")
)

var v = func(string, ...interface{}) {}

func main() {
	flag.Parse()
	if len(flag.Args()) != 1 {
		log.Fatal("Usage: uefiboot <payload>")
	}
	fv, err := uefi.New(flag.Args()[0])
	if err != nil {
		log.Fatal(err)
	}
	if *serialAddr > math.MaxUint32 {
		log.Fatal("Unsupported: serial_addr is greater than uint32.")
	}
	fv.ImageBase, fv.SerialConfig = uintptr(*imageBase), uefi.SerialPortConfig{
		Type:       uefi.SerialPortTypeIO,
		BaseAddr:   uint32(*serialAddr),
		RegWidth:   uint32(*serialWidth),
		InputHertz: uint32(*serialHertz),
		Baud:       uint32(*serialBaud),
	}

	if err := fv.Load(*debug); err != nil {
		log.Fatal(err)
	}

	if err := boot.Execute(); err != nil {
		log.Fatal(err)
	}
}
