// Copyright 2020 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// acpi2coreboot builds coreboot C code from ACPI tables
//
// Synopsis:
//     acpi2coreboot
//
// Description:
//	Read ACPI tables from stdin and write coreboot C code to stdout
//      for those tables we support.
//
//      For MADT, one might to this:
//      acpicat | acpigrep APIC | acpi2coreboot
//	or, on Linux:
//	sudo cat  /sys/firmware/acpi/tables/APIC | ./acpi2coreboot
//
// Options:
package main

import (
	"flag"
	"log"
	"os"

	"github.com/u-root/u-root/pkg/acpi"
	"github.com/u-root/u-root/pkg/acpi/coreboot"
)

func main() {
	flag.Parse()
	tabs, err := acpi.RawFromFile(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}
	for _, t := range tabs {
		cb, err := coreboot.NewCorebooter(t)
		if err != nil {
			log.Fatal(err)
		}
		if err := cb.Coreboot(os.Stdout); err != nil {
			log.Fatal(err)
		}
	}
}
