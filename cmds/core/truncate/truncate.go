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

	if !size.IsSet {
		log.Println("truncate: ERROR: You need to specify -s <number>.")
		usageAndExit()
	}
	if flag.NArg() == 0 {
		log.Println("truncate: ERROR: You need to specify one or more files as argument.")
		usageAndExit()
	}

	for _, fname := range flag.Args() {
		st, err := os.Stat(fname)
		if os.IsNotExist(err) && !*create {
			if err = ioutil.WriteFile(fname, []byte{}, 0644); err != nil {
				log.Fatalf("truncate: ERROR: %v\n", err)
			}
			if st, err = os.Stat(fname); err != nil {
				log.Fatalf("truncate: ERROR: could not stat newly created file: %v\n", err)
			}
		}

		final := size.Value // base case
		if size.ExplicitSign != unit.None {
			final += st.Size() // in case of '-', size.Value is already negative
		}
		if final < 0 {
			final = 0
		}

		// intentionally ignore, like GNU truncate
		os.Truncate(fname, final)
	}
}
