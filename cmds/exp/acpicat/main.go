// Copyright 2020 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// acpicat cats ACPI tables from the kernel.
// The default method is "files", commonly provided in Linux via /sys.
// Other methods are available depending on the platform.
// Further selection of which tables are used can be done with acpigrep.
package main

import (
	"flag"
	"log"
	"os"

	"github.com/u-root/u-root/pkg/acpi"
)

var (
	source = flag.String("s", acpi.DefaultMethod, "source of the tables")
	debug  = flag.Bool("d", false, "Enable debug prints")
)

func main() {
	flag.Parse()
	if *debug {
		acpi.Debug = log.Printf
	}
	t, err := acpi.ReadTables(*source)
	if err != nil {
		log.Fatal(err)
	}
	if len(t) == 0 {
		log.Fatalf("%s: no tables read", *source)
	}
	if err := acpi.WriteTables(os.Stdout, t[0], t[1:]...); err != nil {
		log.Fatal(err)
	}
}
