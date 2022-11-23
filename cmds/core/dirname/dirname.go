// Copyright 2016-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// dirname prints out the directory name of one or more args.
// If no arg is given it returns an error and prints a message which,
// per the man page, is incorrect, but per the standard, is correct.
package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
)

var ErrNoArg = errors.New("missing operand")

func run(out io.Writer, args []string) error {
	if len(args) < 1 {
		return ErrNoArg
	}

	for _, n := range args {
		fmt.Fprintln(out, filepath.Dir(n))
	}
	return nil
}

func main() {
	if err := run(os.Stdout, os.Args[1:]); err != nil {
		log.Fatalf("dirname: %v", err)
	}
}
