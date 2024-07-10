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
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/u-root/u-root/pkg/ts"
	"github.com/u-root/u-root/pkg/uroot/unixflag"
)

func run(stdin io.Reader, stdout io.Writer, args []string) error {
	var first, relative bool
	f := flag.NewFlagSet(args[0], flag.ExitOnError)
	f.BoolVar(&first, "f", false, "All timestamps are relative to the first character")
	f.BoolVar(&relative, "R", false, "Timestamps are relative to the previous timestamp")

	f.Usage = func() {
		fmt.Printf("Usage: ts [options]\n")
		f.PrintDefaults()
	}

	f.Parse(unixflag.ArgsToGoArgs(args[1:]))
	if f.NArg() != 0 {
		f.Usage()
		return fmt.Errorf("invalid use")
	}

	t := ts.New(stdin)
	t.ResetTimeOnNextRead = first
	if relative {
		t.Format = ts.NewRelativeFormat()
	}

	_, err := io.Copy(stdout, t)
	return err
}

func main() {
	if err := run(os.Stdin, os.Stdout, os.Args); err != nil {
		log.Fatal(err)
	}
}
