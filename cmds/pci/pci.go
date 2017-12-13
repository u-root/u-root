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
//     -c: dump config space
//     -s: specify glob for choosing devices.
package main

import (
	"flag"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/u-root/u-root/pkg/pci"
)

var (
	numbers    = flag.Bool("n", false, "Show numeric IDs")
	dumpConfig = flag.Bool("c", false, "Dump config space")
	devs       = flag.String("s", "*", "Devices to match")
	format     = map[int]string{
		32: "%08x:%08x",
		16: "%08x:%04x",
		8:  "%08x:%02x",
	}
)

// maybe we need a better syntax than the standard pcitools?
func registers(d pci.Devices, cmds ...string) {
	var justCheck bool
	for _, c := range cmds {
		// TODO: replace this nonsense with a state machine.
		// Split into register and value
		rv := strings.Split(c, "=")
		if len(rv) != 1 && len(rv) != 2 {
			log.Printf("%v: only one = allowed. Due to this error no more commands will be issued", c)
			justCheck = true
			continue
		}

		// Split into register offset and size
		rs := strings.Split(rv[0], ".")
		if len(rs) != 1 && len(rs) != 2 {
			log.Printf("%v: only one . allowed. Due to this error no more commands will be issued", rv[1])
			justCheck = true
			continue
		}
		s := 32
		if len(rs) == 2 {
			switch rs[1] {
			default:
				log.Printf("Bad size: %v. Due to this error no more commands will be issued", rs[1])
				justCheck = true
				continue
			case "l":
			case "w":
				s = 16
			case "b":
				s = 8
			}
		}
		if justCheck {
			continue
		}
		reg, err := strconv.ParseUint(rs[0], 0, 16)
		if err != nil {
			log.Printf("%v:%v. Due to this error no more commands will be issued", rs[0], err)
			justCheck = true
			continue
		}
		if len(rv) == 1 {
			v, err := d.ReadConfigRegister(int64(reg), int64(s))
			if err != nil {
				log.Printf("%v:%v. Due to this error no more commands will be issued", rv[1], err)
				justCheck = true
				continue
			}
			// Should this go in the package somewhere? Not sure.
			for i := range v {
				d[i].ExtraInfo = append(d[i].ExtraInfo, fmt.Sprintf(format[s], reg, v[i]))
			}
		}
		if len(rv) == 2 {
			val, err := strconv.ParseUint(rv[1], 0, s)
			if err != nil {
				log.Printf("%v. Due to this error no more commands will be issued", err)
				justCheck = true
				continue
			}
			if err := d.WriteConfigRegister(int64(reg), int64(s), val); err != nil {
				log.Printf("%v:%v. Due to this error no more commands will be issued", rv[1], err)
				justCheck = true
				continue
			}
		}

	}
}
func main() {
	flag.Parse()
	r, err := pci.NewBusReader(*devs)
	if err != nil {
		log.Fatalf("%v", err)
	}

	d, err := r.Read()
	if err != nil {
		log.Fatalf("Read: %v", err)
	}

	if !*numbers {
		d.SetVendorDeviceName()
	}
	if len(flag.Args()) > 0 {
		registers(d, flag.Args()...)
	}
	if *dumpConfig {
		d.ReadConfig()
	}
	fmt.Print(d)
}
