// Copyright 2020 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// ts prepends each line of stdin with a timestamp.
//
// Synopsis:
//
//	ts
package main

import (
	"io"
	"log"
	"os"

	flag "github.com/spf13/pflag"
	"github.com/u-root/u-root/pkg/ts"
)

var (
	first    = flag.Bool("f", false, "All timestamps are relative to the first character")
	relative = flag.Bool("R", false, "Timestamps are relative to the previous timestamp")
)

func run(stdin io.Reader, stdout io.Writer, first, relative bool) error {
	t := ts.New(stdin)
	t.ResetTimeOnNextRead = first
	if relative {
		t.Format = ts.NewRelativeFormat()
	}

	_, err := io.Copy(stdout, t)
	return err
}

func main() {
	flag.Parse()
	if flag.NArg() != 0 {
		log.Fatal("Usage: ts")
	}

	if err := run(os.Stdin, os.Stdout, *first, *relative); err != nil {
		log.Fatal(err)
	}
}
