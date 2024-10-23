// Copyright 2012-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// pci: show pci bus vendor ids and other info
//
// Description:
//
//	List the PCI bus, with names if possible.
//
// Options:
//
//	-n: just show numbers
//	-c: dump config space
//	-s: specify glob for choosing devices.
package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/u-root/u-root/pkg/pci"
)

type cmd struct {
	w         io.Writer
	osargs    []string
	args      []string
	numbers   bool
	devGlobs  string
	dumpJSON  bool
	verbosity int
	hexdump   int
	readJSON  string
	flags     *flag.FlagSet
}

// errBadJSON is for any JSON Unmarshal error.
// The JSON package predates wrapped errors. This error is a placeholder for all
// "something went wrong" JSON parsing errors.
var errBadJSON = errors.New("JSON parsing failed")

func command(w io.Writer, args ...string) *cmd {
	f := flag.NewFlagSet("pci", flag.ExitOnError)

	c := &cmd{
		w:      w,
		osargs: args,
		flags:  f,
	}

	f.BoolVar(&c.numbers, "n", false, "Show numeric IDs")
	f.StringVar(&c.devGlobs, "s", "*", ",-seperated list of globs in /sys/bus/pci/devices/ glob, e.g. 0000:*,0001:*")
	f.BoolVar(&c.dumpJSON, "j", false, "Dump the bus in JSON")
	f.IntVar(&c.verbosity, "v", 0, "verbosity")
	f.IntVar(&c.hexdump, "x", 0, "hexdump the config space")
	f.StringVar(&c.readJSON, "J", "", "Read JSON in instead of /sys")
	f.Parse(c.osargs)
	c.args = f.Args()
	return c
}

var format = map[int]string{
	32: "%08x:%08x",
	16: "%08x:%04x",
	8:  "%08x:%02x",
}

// maybe we need a better syntax than the standard pcitools?
func registers(d pci.Devices, cmds ...string) error {
	var justCheck bool
	var err error
	for _, c := range cmds {
		// TODO: replace this nonsense with a state machine.
		// Split into register and value
		rv := strings.Split(c, "=")
		if len(rv) != 1 && len(rv) != 2 {
			log.Printf("%v: only one = allowed. Due to this error no more commands will be issued", c)
			err = errors.Join(err, fmt.Errorf("%v:only one = allowed.%w", c, strconv.ErrSyntax))
			justCheck = true
			continue
		}

		// Split into register offset and size
		rs := strings.Split(rv[0], ".")
		if len(rs) != 1 && len(rs) != 2 {
			log.Printf("%v: only one . allowed. Due to this error no more commands will be issued", rv[1])
			err = errors.Join(err, fmt.Errorf("%v:only one . allowed.%w", c, strconv.ErrSyntax))
			justCheck = true
			continue
		}
		s := 32
		if len(rs) == 2 {
			switch rs[1] {
			default:
				log.Printf("Bad size: %v. Due to this error no more commands will be issued", rs[1])
				err = errors.Join(err, fmt.Errorf("%v:bad size.%w", rs[1], strconv.ErrSyntax))
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
		reg, e := strconv.ParseUint(rs[0], 0, 16)
		if e != nil {
			log.Printf("%v:%v. Due to this error no more commands will be issued", rs[0], e)
			err = errors.Join(err, fmt.Errorf("%v:%w", rs[0], e))
			justCheck = true
			continue
		}
		if len(rv) == 1 {
			v, e := d.ReadConfigRegister(int64(reg), int64(s))
			if e != nil {
				log.Printf("%v:%v. Due to this error no more commands will be issued", rv[0], e)
				err = errors.Join(err, fmt.Errorf("%v:%w", c, e))
				justCheck = true
				continue
			}
			// Should this go in the package somewhere? Not sure.
			for i := range v {
				d[i].ExtraInfo = append(d[i].ExtraInfo, fmt.Sprintf(format[s], reg, v[i]))
			}
		}
		if len(rv) == 2 {
			val, e := strconv.ParseUint(rv[1], 0, s)
			if e != nil {
				log.Printf("%v. Due to this error no more commands will be issued", e)
				err = errors.Join(err, fmt.Errorf("%w", e))
				justCheck = true
				continue
			}
			if e := d.WriteConfigRegister(int64(reg), int64(s), val); e != nil {
				log.Printf("%v:%v. Due to this error no more commands will be issued", rv[1], e)
				err = errors.Join(err, fmt.Errorf("%v:%w", rv[1], e))
				justCheck = true
				continue
			}
		}

	}
	return err
}

func (c *cmd) run() error {
	var dumpSize int
	switch c.hexdump {
	case 4:
		dumpSize = 4096
	case 3:
		dumpSize = 256
	case 2: // lspci disallows this value
		dumpSize = 256
	case 1:
		dumpSize = 64
	}
	r, err := pci.NewBusReader(strings.Split(c.devGlobs, ",")...)
	if err != nil {
		return err
	}

	var d pci.Devices
	if len(c.readJSON) != 0 {
		b, err := os.ReadFile(c.readJSON)
		if err != nil {
			return err
		}
		if err := json.Unmarshal(b, &d); err != nil {
			return fmt.Errorf("%w:%w", err, errBadJSON)
		}

	} else {
		if d, err = r.Read(); err != nil {
			return err
		}
	}

	if !c.numbers || c.dumpJSON {
		d.SetVendorDeviceName(pci.IDs)
	}
	if len(c.args) > 0 {
		if err := registers(d, c.args...); err != nil {
			return err
		}
	}
	if c.dumpJSON {
		o, err := json.MarshalIndent(d, "", "\t")
		if err != nil {
			return err
		}
		fmt.Fprintf(c.w, "%s", string(o))
		return nil
	}
	if err := d.Print(c.w, c.verbosity, dumpSize); err != nil {
		return err
	}
	return nil
}

func main() {
	c := command(os.Stdout, os.Args[1:]...)
	if err := c.run(); err != nil {
		log.Fatal(err)
	}
}
