// Copyright 2012-2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
)

func run(stdout io.Writer, args ...string) error {
	var errs error

	for _, arg := range args {
		absPath, err := filepath.Abs(arg)
		if err != nil {
			errs = errors.Join(errs, err)
			continue
		}

		realPath, err := filepath.EvalSymlinks(absPath)
		if err != nil {
			errs = errors.Join(errs, err)
			continue
		}

		fmt.Fprintf(stdout, "%s\n", filepath.Clean(realPath))
	}

	return errs
}

func main() {
	q := flag.Bool("q", false, "quiet mode")
	flag.Parse()
	if err := run(os.Stdout, flag.Args()...); err != nil {
		if *q {
			os.Exit(1)
		}
		log.Fatal(err)
	}
}
