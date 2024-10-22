// Copyright 2013-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Truncate - shrink or extend the size of a file to the specified size
//
// Synopsis:
//
//	truncate [OPTIONS] [FILE]...
//
// Options:
//
//	-s: size in bytes
//	-r: reference file for size
//	-c: do not create any files
//
// Author:
//
//	Roland Kammerer <dev.rck@gmail.com>
package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/rck/unit"
	"github.com/u-root/u-root/pkg/uroot/util"
)

const usage = "truncate [-c] -s size file..."

var (
	create = flag.Bool("c", false, "Do not create files.")
	size   = unit.MustNewUnit(unit.DefaultUnits).MustNewValue(1, unit.None)
	rfile  = flag.String("r", "", "Reference file for size")
)

func init() {
	flag.Var(size, "s", "Size in bytes, prefixes +/- are allowed")
	flag.Usage = util.Usage(flag.Usage, usage)
}

func truncate(args ...string) error {
	if !size.IsSet && *rfile == "" {
		return fmt.Errorf("you need to specify size via -s <number> or -r <rfile>")
	}
	if size.IsSet && *rfile != "" {
		return fmt.Errorf("you need to specify size via -s <number> or -r <rfile>")
	}
	if len(args) == 0 {
		return fmt.Errorf("you need to specify one or more files as argument")
	}

	for _, fname := range args {

		var final int64
		st, err := os.Stat(fname)
		if os.IsNotExist(err) && !*create {
			if err = os.WriteFile(fname, []byte{}, 0o644); err != nil {
				return fmt.Errorf("%w", err)
			}
			if st, err = os.Stat(fname); err != nil {
				return fmt.Errorf("could not stat newly created file: %w", err)
			}
		}
		if *rfile != "" {
			if st, err = os.Stat(*rfile); err != nil {
				return fmt.Errorf("could not stat reference file: %w", err)
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
	return nil
}

func main() {
	flag.Parse()
	if err := truncate(flag.Args()...); err != nil {
		flag.Usage()
		log.Fatal(err)
	}
}
