// Copyright 2012-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// smn: read or write registers in the System Management Network on AMD cpus
//
// Synopsis:
//
//	snm [0 or more addresses]
//
// N.B. having no addresses is an easy way to see if you can
// access PCI at all.
//
// Description:
//
//	read and write System Management Network registers
//
// Options:
//
//	-s: device glob in the form tttt:bb:dd.fn with * as needed
//	-n: number of 32-bit words to dump.
//	-v: 32-bit value to write
//	-w: write the value to the register(s)
package main

import (
	"flag"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/u-root/u-root/pkg/pci"
)

// Index/Data register pairs, as, e.g., cf8/cfc for PCI,
// are known to be a terrible idea, from almost any point of view.
// It took years to kill them on regular PCI.
// AMD brought them back for the SMN. Bummer.
const (
	regIndex = 0xa0
	regData  = 0xa4
)

var (
	devs  = flag.String("s", "0000:00:00.0", "Glob for northbridge")
	n     = flag.Uint("n", 1, "Number 32-bit words to dump/set")
	val   = flag.Uint64("v", 0, "Val to set on write")
	write = flag.Bool("w", false, "Write a value")
)

func usage() {
	log.Fatal("Usage: smn [-w] [-d glob] address [# 32-words to read | 32-bit value to write]")
}

func main() {
	flag.Parse()
	r, err := pci.NewBusReader(strings.Split(*devs, ",")...)
	if err != nil {
		log.Fatalf("%v", err)
	}

	d, err := r.Read()
	if err != nil {
		log.Fatalf("Read: %v", err)
	}

	a := flag.Args()
	if uint32(*val>>32) != 0 {
		log.Fatalf("Value:%#x is not a 32-bit number", *val)
	}
	for i := range a {
		addr, err := strconv.ParseUint(a[i], 16, 32)
		if err != nil {
			log.Fatal(err)
		}
		switch *write {
		case true:
			if err := d.WriteConfigRegister(regIndex, 32, addr); err != nil {
				log.Fatal(err)
			}
			if err := d.WriteConfigRegister(regData, 32, *val); err != nil {
				log.Fatal(err)
			}
		case false:
			for i := addr; i < addr+uint64(*n); i += 4 {
				if err := d.WriteConfigRegister(regIndex, 32, i); err != nil {
					log.Fatal(err)
				}
				// SMN data is 32 bits!
				dat, err := d.ReadConfigRegister(regData, 32)
				if err != nil {
					log.Fatal(err)
				}
				fmt.Printf("%#x:", i)
				for i := range dat {
					fmt.Printf("%s:%#x,", d[i].Addr, dat[i])
				}
				fmt.Println()
			}
		}
	}
}
