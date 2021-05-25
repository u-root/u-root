// Copyright 2013-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Truncate - shrink or extend the size of a file to the specified size
//
// Synopsis:
//     truncate [OPTIONS] [FILE]...
//
// Options:
//     -s: size in bytes
//     -r: reference file for size
//     -c: do not create any files
//
// Author:
//     Roland Kammerer <dev.rck@gmail.com>
package main

import (
	"flag"
	"io/ioutil"
	"log"
	"os"

	"github.com/rck/unit"
)

const cmd = "truncate [-c] -s size file..."

var (
	create = flag.Bool("c", false, "Do not create files.")
	size   = unit.MustNewUnit(unit.DefaultUnits).MustNewValue(1, unit.None)
	rfile  = flag.String("r", "", "Reference file for size")
)

func init() {
	flag.Var(size, "s", "Size in bytes, prefixes +/- are allowed")

	defUsage := flag.Usage
	flag.Usage = func() {
		os.Args[0] = cmd
		defUsage()
	}
}

func usageAndExit() {
	flag.Usage()
	os.Exit(1)
}

func main() {
	flag.Parse()

	if !size.IsSet && *rfile == "" {
		log.Println("truncate: ERROR: You need to specify size via -s <number> or -r <rfile>.")
		usageAndExit()
	}
	if size.IsSet && *rfile != "" {
		log.Println("truncate: ERROR: You need to specify size via -s <number> or -r <rfile>.")
		usageAndExit()
	}
	if flag.NArg() == 0 {
		log.Println("truncate: ERROR: You need to specify one or more files as argument.")
		usageAndExit()
	}

	for _, fname := range flag.Args() {

		var final int64
		st, err := os.Stat(fname)
		if os.IsNotExist(err) && !*create {
			if err = ioutil.WriteFile(fname, []byte{}, 0644); err != nil {
				log.Fatalf("truncate: ERROR: %v\n", err)
			}
			if st, err = os.Stat(fname); err != nil {
				log.Fatalf("truncate: ERROR: could not stat newly created file: %v\n", err)
			}
		}
		if *rfile != "" {
			if st, err = os.Stat(*rfile); err != nil {
				log.Fatalf("truncate: ERROR: could not stat reference file: %v\n", err)
			}
			final = st.Size()
		} else if size.IsSet {
			final = size.Value // base case
			if size.ExplicitSign != unit.None {
				final += st.Size() // in case of '-', size.Value is already negative
			}
			if final < 0 {
				final = 0
			}
		}

		// intentionally ignore, like GNU truncate
		os.Truncate(fname, final)
	}
}
