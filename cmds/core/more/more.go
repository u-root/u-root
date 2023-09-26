// Copyright 2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// More pages through files without any terminal trickery.
//
// Synopsis:
//
//	more [OPTIONS] FILE
//
// Description:
//
//	Admittedly, this does not follow the conventions of GNU more. Instead,
//	it is built with the goal of not relying on any special ttys, ioctls or
//	special ANSI escapes. This is ideal when your terminal is already
//	borked. For bells and whistles, look at less.
//
// Options:
//
//	--lines NUMBER: screen size in number of lines
package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"

	flag "github.com/spf13/pflag"
)

var lines = flag.Int("lines", 40, "screen size in number of lines")
var errOnlyOneFile = fmt.Errorf("more can only take one file")
var errLinesMustBePositive = fmt.Errorf("lines must be positive")

func run(stdin io.Reader, stdout io.Writer, lines int, args []string) error {
	if len(args) != 1 {
		return errOnlyOneFile
	}
	if lines <= 0 {
		return errLinesMustBePositive
	}

	f, err := os.Open(args[0])
	if err != nil {
		return err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for i := 0; scanner.Scan(); i++ {
		if (i+1)%lines == 0 {
			fmt.Fprint(stdout, scanner.Text())
			c := make([]byte, 1)
			// We expect the OS to echo the newline character.
			if _, err := stdin.Read(c); err != nil {
				return err
			}
		} else {
			fmt.Fprintln(stdout, scanner.Text())
		}
	}
	return scanner.Err()
}

func main() {
	flag.Parse()
	if err := run(os.Stdin, os.Stderr, *lines, flag.Args()); err != nil {
		log.Fatal(err)
	}
}
