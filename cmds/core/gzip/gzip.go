// Copyright 2012-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// gzip compresses files using gzip compression.
package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/u-root/u-root/pkg/gzip"
)

var cmdLine = flag.CommandLine

func usage() {
	fmt.Fprintf(os.Stderr, "Usage of %s:\n", filepath.Base(os.Args[0]))
	cmdLine.PrintDefaults()
}

func main() {
	var opts gzip.Options

	cmdLine.Usage = usage

	if err := opts.ParseArgs(os.Args, cmdLine); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		cmdLine.Usage()
		os.Exit(2)
	}

	var input []gzip.File
	args := cmdLine.Args()

	if len(args) == 0 {
		// no args given, compress stdin to stdout
		input = append(input, gzip.File{Options: &opts})
	} else {
		for _, arg := range args {
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
			os.Exit(1)
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
}
