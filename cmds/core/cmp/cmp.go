// Copyright 2013-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Cmp compares two files and prints a message if their contents differ.
//
// Synopsis:
//     cmp [–lLs] FILE1 FILE2 [OFFSET1 [OFFSET2]]
//
// Description:
//     If offsets are given, comparison starts at the designated byte position
//     of the corresponding file.
//
//     Offsets that begin with 0x are hexadecimal; with 0, octal; with anything
//     else, decimal.
//
// Options:
//     –l: Print the byte number (decimal) and the differing bytes (octal) for
//         each difference.
//     –L: Print the line number of the first differing byte.
//     –s: Print nothing for differing files, but set the exit status.
package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/rck/unit"
)

var (
	long   = flag.Bool("l", false, "print the byte number (decimal) and the differing bytes (hexadecimal) for each difference")
	line   = flag.Bool("L", false, "print the line number of the first differing byte")
	silent = flag.Bool("s", false, "print nothing for differing files, but set the exit status")
)

func emit(rs io.ReadSeeker, c chan byte, offset int64) error {
	if offset > 0 {
		if _, err := rs.Seek(offset, 0); err != nil {
			log.Fatalf("%v", err)
		}
	}

	b := bufio.NewReader(rs)
	for {
		b, err := b.ReadByte()
		if err != nil {
			close(c)
			return err
		}
		c <- b
	}
}

func openFile(name string) (*os.File, error) {
	var f *os.File
	var err error

	if name == "-" {
		f = os.Stdin
	} else {
		f, err = os.Open(name)
	}

	return f, err
}

// cmp is defined to fail with exit code 2
func failf(format string, a ...interface{}) {
	fmt.Fprintf(os.Stderr, format, a...)
	os.Exit(2)
}

func main() {
	flag.Parse()
	var offset [2]int64
	var f *os.File
	var err error

	fnames := flag.Args()

	cmpUnits := unit.DefaultUnits

	off, err := unit.NewUnit(cmpUnits)
	if err != nil {
		failf("Could not create unit based on mapping: %v\n", err)
	}

	var v *unit.Value
	switch len(fnames) {
	case 2:
	case 3:
		if v, err = off.ValueFromString(fnames[2]); err != nil {
			failf("bad offset1: %s: %v\n", fnames[2], err)
		}
		offset[0] = v.Value
	case 4:
		if v, err = off.ValueFromString(fnames[2]); err != nil {
			failf("bad offset1: %s: %v\n", fnames[2], err)
		}
		offset[0] = v.Value

		if v, err = off.ValueFromString(fnames[3]); err != nil {
			failf("bad offset2: %s: %v\n", fnames[3], err)
		}
		offset[1] = v.Value
	default:
		failf("expected two filenames (and one to two optional offsets), got %d", len(fnames))
	}

	c := make([]chan byte, 2)

	for i := 0; i < 2; i++ {
		if f, err = openFile(fnames[i]); err != nil {
			failf("Failed to open %s: %v", fnames[i], err)
		}
		c[i] = make(chan byte, 8192)
		go emit(f, c[i], offset[i])
	}

	lineno, charno := int64(1), int64(1)
	var b1, b2 byte
	for {
		b1 = <-c[0]
		b2 = <-c[1]

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
				fmt.Fprintf(os.Stderr, "%8d %#.2o %#.2o\n", charno, b1, b2)
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
