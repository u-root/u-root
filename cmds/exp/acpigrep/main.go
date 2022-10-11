// Copyright 2020 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// grep a stream of ACPI tables by regexp
//
// Synopsis:
//
//	acpigrep [-v] [-d] regexp
//
// Description:
//
//		Read tables from stdin and write tables with ACPI signatures
//		matching a pattern (or not matching, with -v, as in grep)
//		to stdout.
//
//		Read all tables from sysfs and discard MADT
//		sudo cat /sys/firmware/acpi/tables/[A-Z]* | ./acpigrep -v MADT
//
//		Read a large blob and print out its tables
//		acpigrep -d '.*' >/dev/null < blob
//
//	     Read all the files in /sys and discard any DSDT. Useful for coreboot work.
//		sudo cat /sys/firmware/acpi/tables/[A-Z]* | ./acpigrep -v DSDT > nodsdt.bin
//
//		Read all the files, keeping only SRAT and MADT
//		sudo cat /sys/firmware/acpi/tables/[A-Z]* | ./acpigrep 'MADT|SRAT' > madtsrat.bin
//
//		Read all the files, keeping only SRAT and MADT, and print what is done
//		sudo cat /sys/firmware/acpi/tables/[A-Z]* | ./acpigrep -d 'MADT|SRAT' > madtsrat.bin
//
// Options:
//
//	-d print debug information about what is kept and what is discarded.
//	-v reverse the sense of the match to "discard is matching"
package main

import (
	"flag"
	"log"
	"os"
	"regexp"

	"github.com/u-root/u-root/pkg/acpi"
)

var (
	v     = flag.Bool("v", false, "Only non-matching signatures will be kept")
	d     = flag.Bool("d", false, "Print debug messages")
	debug = func(string, ...interface{}) {}
)

func main() {
	flag.Parse()
	if *d {
		debug = log.Printf
	}
	if len(flag.Args()) != 1 {
		log.Fatal("Usage: acpigrep [-v] [-d] pattern")
	}
	r := regexp.MustCompile(flag.Args()[0])
	tabs, err := acpi.RawFromFile(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}
	for _, t := range tabs {
		m := r.MatchString(t.Sig())
		if m == *v {
			debug("Dropping %s", acpi.String(t))
			continue
		}
		debug("Keeping %s", acpi.String(t))
		os.Stdout.Write(t.Data())
	}
}
