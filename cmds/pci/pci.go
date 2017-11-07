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

var numbers = flag.Bool("n", false, "Just show numbers")

func main() {
	flag.Parse()
	r, err := pci.NewBusReader(numbers)
	if err != nil {
		log.Fatalf("%v", err)
	}
	devs, err := r.Read()
	if err != nil {
		log.Fatalf("Read: %v", err)
	}
	for _, p := range devs {
		fmt.Printf("%v\n", p)
	}

}
