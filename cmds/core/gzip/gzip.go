// Copyright 2012-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// gzip compresses files using gzip compression.
package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/u-root/u-root/pkg/gzip"
)

func run(args []string) error {
	fs := flag.NewFlagSet("gzip", flag.ContinueOnError)
	fs.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of gzip\n")
		fs.PrintDefaults()
	}

	var opts gzip.Options

	if err := opts.ParseArgs("gzip", args, fs); err != nil {
		if errors.Is(err, gzip.ErrHelp) {
			fs.Usage()
			return nil
		}
		return err
	}

	var input []gzip.File
	if len(fs.Args()) == 0 {
		// no args given, compress stdin to stdout
		input = append(input, gzip.File{Options: &opts})
	} else {
		for _, arg := range fs.Args() {
			input = append(input, gzip.File{Path: arg, Options: &opts})
		}
	}

	for _, f := range input {
		if err := f.CheckPath(); err != nil {
			if !opts.Quiet {
				fmt.Fprintf(os.Stderr, "%s\n", err)
			}
			continue
		}

		if err := f.CheckOutputStdout(); err != nil {
			if !opts.Quiet {
				fmt.Fprintf(os.Stderr, "%s\n", err)
			}
			return err
		}

		if err := f.CheckOutputPath(); err != nil {
			if !opts.Quiet {
				fmt.Fprintf(os.Stderr, "%s\n", err)
			}
			continue
		}

		if err := f.Process(); err != nil {
			if !opts.Quiet {
				fmt.Fprintf(os.Stderr, "%s\n", err)
			}
		}

		if err := f.Cleanup(); err != nil {
			if !opts.Quiet {
				fmt.Fprintf(os.Stderr, "%s\n", err)
			}
			continue
		}
	}

	return nil
}

func main() {
	if err := run(os.Args[1:]); err != nil {
		log.Fatalf("gzip: %v", err)
	}
}
