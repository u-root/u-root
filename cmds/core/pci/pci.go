// Copyright 2012-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// pci: show pci bus vendor ids and other info
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
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"

	flag "github.com/pborman/getopt/v2"
	"github.com/u-root/u-root/pkg/pci"
)

var (
	format = map[int]string{
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
				log.Printf("%v:%v. Due to this error no more commands will be issued", rv[0], err)
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
	var dumpSize int
	numbers := flag.Bool('n', "Show numeric IDs")
	devs := flag.StringLong("select", 's', "*", "Devices to match")
	j := flag.BoolLong("json", 'j', "Dump the bus in JSON")
	v := flag.Counter('v', "verbosity")
	x := flag.Counter('x', "hexdump the config space")
	readJSON := flag.StringLong("JSON", 'J', "", "Read JSON in instead of /sys")

	flag.Parse()

	switch *x {
	case 4:
		dumpSize = 4096
	case 3:
		dumpSize = 256
	case 2: // lspci disallows this value
		dumpSize = 256
	case 1:
		dumpSize = 64
	}
	r, err := pci.NewBusReader(strings.Split(*devs, ",")...)
	if err != nil {
		log.Fatalf("%v", err)
	}

	var d pci.Devices
	if len(*readJSON) != 0 {
		b, err := ioutil.ReadFile(*readJSON)
		if err != nil {
			log.Fatal(err)
		}
		if err := json.Unmarshal(b, &d); err != nil {
			log.Fatal(err)
		}

	} else {
		if d, err = r.Read(); err != nil {
			log.Fatal(err)
		}
	}

	if !*numbers || *j {
		d.SetVendorDeviceName()
	}
	if len(flag.Args()) > 0 {
		registers(d, flag.Args()...)
	}
	if *j {
		o, err := json.MarshalIndent(d, "", "\t")
		if err != nil {
			log.Fatalf("%v", err)
		}
		fmt.Printf("%s", string(o))
		os.Exit(0)
	}
	if err := d.Print(os.Stdout, *v, dumpSize); err != nil {
		log.Fatal(err)
	}
}
