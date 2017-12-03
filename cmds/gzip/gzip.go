// Copyright 2012-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"os"

	"github.com/spf13/pflag"
	"github.com/u-root/u-root/pkg/gzip"
)

func main() {
	var opts gzip.Options
	if err := opts.ParseArgs(); err != nil {
		os.Exit(gzip.ErrorHandler(err, os.Stdout, os.Stderr))
	}

	var input []string
	if opts.Stdin {
		input = []string{"/dev/stdin"}
	} else {
		input = pflag.Args()
	}

	for _, path := range input {
		f := gzip.File{Path: path, Options: &opts}
		if err := f.CheckPath(); err != nil {
			if !opts.Quiet {
				gzip.ErrorHandler(err, os.Stdout, os.Stderr)
			}
			continue
		}

		if err := f.CheckOutputPath(); err != nil {
			if !opts.Quiet {
				gzip.ErrorHandler(err, os.Stdout, os.Stderr)
			}
			continue
		}

		if err := f.Process(); err != nil {
			if !opts.Quiet {
				os.Exit(gzip.ErrorHandler(err, os.Stdout, os.Stderr))
			}
		}

		if err := f.Cleanup(); err != nil {
			if !opts.Quiet {
				gzip.ErrorHandler(err, os.Stdout, os.Stderr)
			}
			continue
		}
	}
}
