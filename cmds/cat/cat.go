// Copyright 2012-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Cat concatenates files and prints to stdout.
//
// Synopsis:
//     cp [-u] [FILES]...
//
// Description:
//     If no files are specified, read from stdin.
//
// Options:
//     -u: ignored flag
package main

import (
	"flag"
	"fmt"
	"io"
	"os"
)

var (
	_ = flag.Bool("u", false, "ignored")
)

func cat(writer io.Writer, files []string) error {
	for _, name := range files {
		f, err := os.Open(name)
		if err != nil {
			return err
		}

		_, err = io.Copy(writer, f)
		if err != nil {
			return err
		}
		f.Close()
	}

	return nil
}

func main() {
	flag.Parse()

	if len(os.Args) == 1 {
		io.Copy(os.Stdout, os.Stdin)
	}

	err := cat(os.Stdout, flag.Args())
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}
