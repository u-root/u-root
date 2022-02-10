// Copyright 2016-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// dirname prints out the directory name of one or more args.
// If no arg is given it returns an error and prints a message which,
// per the man page, is incorrect, but per the standard, is correct.
package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
)

func runDirname(w io.Writer, dirname *string, args []string) error {
	if dirname == nil {
		return fmt.Errorf("dirname: missing operand")
	}
	fmt.Fprintf(w, "%s\n", filepath.Dir(*dirname))
	for _, n := range args {
		if _, err := fmt.Fprintf(w, "%s\n", filepath.Dir(n)); err != nil {
			return err
		}
	}
	return nil
}

func main() {
	if len(os.Args) == 0 {
		log.Fatal()
	}
	var dirname *string
	if len(os.Args[1:]) > 1 {
		*dirname = os.Args[1]
	}
	if err := runDirname(os.Stdout, dirname, os.Args[2:]); err != nil {
		log.Fatal(err)
	}
}
