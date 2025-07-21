// Copyright 2013-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// cmp compares two files and prints a message if their contents differ.
//
// Synopsis:
//
//	cmp [–lLs] FILE1 FILE2 [OFFSET1 [OFFSET2]]
//
// Description:
//
//	If offsets are given, comparison starts at the designated byte position
//	of the corresponding file.
//
//	Offsets that begin with 0x are hexadecimal; with 0, octal; with anything
//	else, decimal.
//
// Options:
//
//	–l: Print the byte number (decimal) and the differing bytes (octal) for
//	    each difference.
//	–L: Print the line number of the first differing byte.
//	–s: Print nothing for differing files, but set the exit status.
//
// What is an error, what goes on stderr, and what goes on stdout in cmp
// is fairly ad-hoc, but go something like this:
// invocation error: unparseable integer: error return from cmp()
// IO error, file too small: output on stderr, error return from cmp()
// Files are different: print difference info on stdout, no error return
// Files are same: no output, no error
package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/rck/unit"
)

const (
	usage = "usage:[-options] file1 file 2 [offset1 [offset 2]]"
)

var (
	long   = flag.Bool("l", false, "print the byte number (decimal) and the differing bytes (hexadecimal) for each difference")
	line   = flag.Bool("L", false, "print the line number of the first differing byte")
	silent = flag.Bool("s", false, "print nothing for differing files, but set the exit status")

	ErrArgCount  = errors.New("arg count")
	ErrBadOffset = errors.New("bad offset")
	ErrDiffer    = errors.New("files differ")
)

func readFileOrStdin(stdin *os.File, name string) (*os.File, error) {
	var f *os.File
	var err error

	if name == "-" {
		f = stdin
	} else {
		f, err = os.Open(name)
	}

	return f, err
}

func cmp(stdout, stderr io.Writer, long, line, silent bool, args ...string) error {
	var offset [2]int64
	var f *os.File
	var err error

	cmpUnits := unit.DefaultUnits

	off, err := unit.NewUnit(cmpUnits)
	if err != nil {
		return fmt.Errorf("could not create unit based on mapping: %w", err)
	}

	var v *unit.Value
	switch len(args) {
	case 2:
	case 3:
		if v, err = off.ValueFromString(args[2]); err != nil {
			fmt.Fprintf(stderr, "bad offset1: %s: %v", args[2], err)
			return fmt.Errorf("%w:%w", err, ErrBadOffset)
		}
		offset[0] = v.Value
	case 4:
		if v, err = off.ValueFromString(args[2]); err != nil {
			fmt.Fprintf(stderr, "bad offset1: %s: %v", args[2], err)
			return fmt.Errorf("%w:%w", err, ErrBadOffset)
		}
		offset[0] = v.Value

		if v, err = off.ValueFromString(args[3]); err != nil {
			fmt.Fprintf(stderr, "bad offset2: %s: %v", args[3], err)
			return fmt.Errorf("%w:%w", err, ErrBadOffset)
		}
		offset[1] = v.Value
	default:
		fmt.Fprint(stderr, usage)
		return ErrArgCount
	}

	c := make([]io.Reader, 2)

	for i := range 2 {
		if f, err = readFileOrStdin(os.Stdin, args[i]); err != nil {
			return fmt.Errorf("failed to open %s: %w", args[i], err)
		}
		if _, err := f.Seek(offset[i], 0); err != nil {
			return fmt.Errorf("%w:%w", err, ErrBadOffset)
		}
		c[i] = bufio.NewReader(f)
	}

	lineno, charno := int64(1), int64(1)

	for {
		var b [2]byte
		_, err1 := c[0].Read(b[:1])
		_, err2 := c[1].Read(b[1:2])

		if err1 != nil || err2 != nil {
			if err1 == io.EOF && err2 == io.EOF {
				return nil
			}
			if err1 != nil {
				fmt.Fprintf(stderr, "%s:%v", args[0], err1)
				return err1
			}
			if err2 != nil {
				fmt.Fprintf(stderr, "%s:%v", args[1], err2)
				return err2
			}
		}

		b1, b2 := b[0], b[1]
		if b1 != b2 {
			if silent {
				return nil
			}
			if line {
				fmt.Fprintf(stdout, "%s %s: char %d line %d", args[0], args[1], charno, lineno)
				return ErrDiffer
			}
			if long {
				fmt.Fprintf(stdout, "%8d %#.2o %#.2o\n", charno, b1, b2)
			} else {
				fmt.Fprintf(stdout, "%s %s: char %d", args[0], args[1], charno)
				return ErrDiffer
			}
		}
		charno++
		if b1 == '\n' {
			lineno++
		}
	}
}

// cmp is defined to fail with exit code 2
func main() {
	flag.Parse()
	if err := cmp(os.Stdout, os.Stderr, *long, *line, *silent, flag.Args()...); err != nil {
		os.Exit(2)
	}
}
