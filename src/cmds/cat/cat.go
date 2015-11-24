// Copyright 2012 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
Cat reads each file from its arguments in sequence and writes it on the standard output.
*/

package main

import (
	"flag"
	"fmt"
	"io"
	"os"
)

var (
	flags struct {
		u bool
	}
	cmd = "cat [-u] [file ...]"
)

func usage() {
	fmt.Fprintln(os.Stderr, "Usage:", cmd)
	flag.PrintDefaults()
	os.Exit(1)
}

func init() {
	flag.BoolVar(&flags.u, "u", false, "Write bytes from the input file to the standard output without delay as each is read.")
	flag.Parse()
	flag.Usage = usage
}

func cat(files []string) error {
	for _, name := range files {
		f, err := os.Open(name)
		if err != nil {
			return err
		}

		_, err = io.Copy(os.Stdout, f)
		if err != nil {
			return err
		}
		f.Close()
	}

	return nil
}

func main() {
	if len(os.Args) == 1 {
		io.Copy(os.Stdout, os.Stdin)
	} else if len(flag.Args()) == 0 { // user has passed the -u flag but didn't added files in the args.
		usage()
		os.Exit(1)
	}

	err := cat(flag.Args())
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}
