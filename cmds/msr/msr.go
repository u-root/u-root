// Copyright 2012-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// msr lets you read and write an msr for one or more cores.
// The cores are specified via a filepath.Glob string.
// The string should be for core number only, with no
// surrounding paths, e.g. you use 0 for core 0, not
// /dev/cpu/0/msr.
// To specify all cores, use '*'
// To specify all cores with two digits, use '??'
// To specify all odd cores, use '*[13579]'
// To specify, e.g., all the even cores, use '*[02468]'.
// Usage:
// msr r glob 32-bit-msr-number
// msr w glob 32-bit-msr-number 64-bit-value
// For each MSR operation msr will print an error if any.
// If your kernel does not have MSRs for any reason,
// this will fail due to file access. But it's quite possible
// that non-x86 architectures might someday implement MSRs,
// which on (e.g.) PPC might have a slightly different name
// (DICR) but would implement the same kinds of functions.
package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
)

const usage = `msr r glob register
msr w glob register value`

func main() {
	a := os.Args[1:]
	if len(a) < 3 {
		log.Fatal(usage)
	}

	m := msrList(a[1])

	reg, err := strconv.ParseUint(a[2], 0, 32)
	if err != nil {
		log.Fatalf("%v: %v", a[2], err)
	}
	switch a[0] {
	case "r":
		data, errs := rdmsr(m, uint32(reg))
		for i := range m {
			if errs[i] != nil {
				fmt.Printf("%v: %v\n", m[i], errs[i])
			} else {
				fmt.Printf("%v: %#016x\n", m[i], data[i])
			}
		}

	case "w":
		// Sadly, we don't get an error on write if the values
		// don't match.  Reading it back to check it is
		// not always going to work. There are many write-only
		// MSRs. If there are no errors there is still no
		// guarantee that it worked, or if we read it we would
		// see what we wrote, or that vmware did not do
		// something stupid, or the x86 did not do something really
		// stupid, or the particular implementation of the x86
		// that we are on did not do something really stupid.
		// Why is it this way? Because vendors hide proprietary
		// information in hidden MSRs, or in hidden fields in MSRs.
		// Checking is just not an option. There, feel better now?
		if len(a) < 4 {
			log.Fatal(usage)
		}
		v, err := strconv.ParseUint(a[3], 0, 64)
		if err != nil {
			log.Fatalf("%v: %v", a, err)
		}
		errs := wrmsr(m, uint32(reg), v)
		for i := range errs {
			if errs[i] != nil {
				fmt.Printf("%v: %v\n", m[i], errs[i])
			}
		}

	default:
		log.Fatalf(usage)
	}
}
