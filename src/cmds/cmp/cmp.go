// Copyright 2013 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
Cmp compares the two files and prints a message if the contents differ.

cmp [ –lLs ] file1 file2 [ offset1 [ offset2 ] ]

The options are:
	–l    Print the byte number (decimal) and the differing bytes (hexadecimal) for each difference.
	–L    Print the line number of the first differing byte.
	–s    Print nothing for differing files, but set the exit status.

-If offsets are given, comparison starts at the designated byte position of the corresponding file.
-Offsets that begin with 0x are hexadecimal; with 0, octal; with anything else, decimal.
*/

package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strconv"
)

var long = flag.Bool("l", false, "print the byte number (decimal) and the differing bytes (hexadecimal) for each difference")
var line = flag.Bool("L", false, "print the line number of the first differing byte")
var silent = flag.Bool("s", false, "print nothing for differing files, but set the exit status")

func emit(f *os.File, c chan byte, offset int64) {
	if offset > 0 {
		f.Seek(offset, 0)
	}

	b := bufio.NewReader(f)
	for {
		b, err := b.ReadByte()
		if err != nil {
			close(c)
			return
		}
		c <- b
	}
}

func main() {
	flag.Parse()
	var offset1, offset2 int64
	var err error

	fnames := flag.Args()
	if len(fnames) != 2 && len(fnames) != 4 {
		fmt.Fprintf(os.Stderr, "expected two filenames (and two optional offsets), got %d", len(fnames))
		os.Exit(1)
	}
	if len(fnames) == 4 {
		offset1, err = strconv.ParseInt(fnames[2], 0, 64)
		if err != nil {
			fmt.Fprintf(os.Stderr, "bad offset1: %s: %v\n", fnames[2], err)
			return
		}
		offset2, err = strconv.ParseInt(fnames[3], 0, 64)
		if err != nil {
			fmt.Fprintf(os.Stderr, "bad offset2: %s: %v\n", fnames[3], err)
			return
		}
	}

	c1 := make(chan byte, 8192)
	c2 := make(chan byte, 8192)

	f, err := os.Open(fnames[0])
	if err != nil {
		fmt.Fprintf(os.Stderr, "error opening %s: %v", fnames[0], err)
		os.Exit(1)
	}
	go emit(f, c1, offset1)

	f, err = os.Open(fnames[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "error opening %s: %v", fnames[1], err)
		os.Exit(2)
	}
	go emit(f, c2, offset2)

	lineno, charno := int64(1), int64(1)
	var b1, b2 byte
	for {
		b1 = <-c1
		b2 = <-c2

		if b1 != b2 {
			if *silent {
				os.Exit(1)
			}
			if *line {
				fmt.Fprintf(os.Stderr, "%s %s differ: char %d line %d\n", fnames[0], fnames[1], charno, lineno)
				os.Exit(1)
			}
			if *long {
				if b1 == '\u0000' {
					fmt.Fprintf(os.Stderr, "EOF on %s\n", fnames[0])
					os.Exit(1)
				}
				if b2 == '\u0000' {
					fmt.Fprintf(os.Stderr, "EOF on %s\n", fnames[1])
					os.Exit(1)
				}
				fmt.Fprintf(os.Stderr, "%8d %#.2x %#.2x\n", charno, b1, b2)
				goto skip
			}
			fmt.Fprintf(os.Stderr, "%s %s differ: char %d\n", fnames[0], fnames[1], charno)
			os.Exit(1)
		}
	skip:
		charno++
		if b1 == '\n' {
			lineno++
		}
		if b1 == '\u0000' && b2 == '\u0000' {
			os.Exit(0)
		}
	}
}
