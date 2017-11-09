// Copyright 2012-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
//     pci: show pci bus vendor ids and other info
//
// Description:
//     List the PCI bus, with names if possible.
//
// Options:
//     -n: just show numbers
package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/u-root/u-root/pkg/pci"
)

var numbers = flag.Bool("n", false, "Show numeric IDs")

func main() {
	flag.Parse()
	r, err := pci.NewBusReader()
	if err != nil {
		log.Fatalf("%v", err)
	}

	d, err := r.Read()
	if err != nil {
		log.Fatalf("Read: %v", err)
	}

	if *numbers {
		fmt.Print(d.ToString(*numbers))

	} else {
		fmt.Print(d)
	}

}
