// Copyright 2012-2020 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//go:build !tinygo || tinygo.enable

// lockmsrs locks important intel MSRs.
//
// All MSRs are specified in the Intel Software developer's manual.
// This seems like a good set of bits to lock down when booting through NERF/LINUXBOOT
// to some other OS. When locked, these MSRs generally prevent
// further modifications until reset.
package main

import (
	"flag"
	"log"
	"os"

	"github.com/u-root/u-root/pkg/msr"
)

var (
	verbose = flag.Bool("v", false, "verbose mode")
	verify  = flag.Bool("V", false, "Verify, do not write")
	debug   = func(string, ...interface{}) {}
)

func main() {
	flag.Parse()

	if *verbose {
		debug = log.Printf
	}

	msr.Debug = debug

	if *verify {
		if err := msr.Locked(); err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}

	cpus, err := msr.AllCPUs()
	if err != nil {
		log.Fatal(err)
	}
	for _, m := range msr.LockIntel {
		debug("Lock MSR %s on cpus %v, clearmask %#08x, setmask %#08x", m.String(), cpus, m.Clear, m.Set)
		var errs []error
		if m.WriteOnly {
			errs = m.Addr.Write(cpus, m.Set)
		} else {
			errs = m.Addr.TestAndSet(cpus, m.Clear, m.Set)
		}

		for i, e := range errs {
			if e != nil {
				// Hope no one ever modifies this slice.
				log.Printf("Error locking msr %v on cpu %v: %v\n", m.Addr.String(), cpus[i], e)
			}
		}
	}
}
