// Copyright 2012-2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build !plan9

// rsdp allows to determine the ACPI RSDP structure address which could
// be passed to the boot command later on
// It must be executed at the system init as it relies on scanning
// the kernel messages which could be quickly filled up in some cases
//
// Synopsis:
//
//	rsdp [-f file]
//
// Description:
//
//	Look for rsdp value in a file, default /dev/kmsg
//
// Example:
//
//	rsdp
//	rsdp -f /path/to/file
package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/u-root/u-root/pkg/acpi"
)

var cmdUsage = "Usage: rsdp"

func usage() {
	log.Fatal(cmdUsage)
}

func main() {
	flag.Parse()
	if flag.NArg() != 0 {
		usage()
	}
	rsdp, err := acpi.GetRSDP()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf(" acpi_rsdp=%#x \n", rsdp.RSDPAddr())
}
