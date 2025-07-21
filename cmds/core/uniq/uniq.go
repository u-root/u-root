// Copyright 2013-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Uniq removes repeated lines.
//
// Synopsis:
//
//	uniq [OPTIONS...] [FILES]...
//
// Description:
//
//	Uniq copies the input file, or the standard input, to the standard
//	output, comparing adjacent lines. In the normal case, the second and
//	succeeding copies of repeated lines are removed. Repeated lines must be
//	adjacent in order to be found.
//
// Options:
//
//	–u:      Print unique lines.
//	–d:      Print (one copy of) duplicated lines.
//	–c:      Prefix a repetition count and a tab to each output line.
//	         Implies –u and –d.
//	-i:      Case insensitive comparison of lines.
//	–f num:  The first num fields together with any blanks before each are
//	         ignored. A field is defined as a string of non–space, non–tab
//	         characters separated by tabs and spaces from its neighbors.
//	-cn num: The first num characters are ignored. Fields are skipped before
//	         characters.
package main

// TODO(aam): -num and +num are not implemented. they're easy to do, just not exactly the
// way that the plan9 uniq does them as we want to avoid polluting the flag parsing libs with
// outdated flags.

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
)

var (
	unique     = flag.Bool("u", false, "print unique lines")
	duplicates = flag.Bool("d", false, "print one copy of duplicated lines")
	count      = flag.Bool("c", false, "prefix a repetition count and a tab for each output line")
	ignoreCase = flag.Bool("i", false, "case insensitive comparison of lines")
)

// var fnum = flag.Int("f", 0, "ignore num fields from beginning of line")
// var cnum = flag.Int("cn", 0, "ignore num characters from beginning of line")

func uniq(r io.Reader, w io.Writer, unique, duplicates, count bool, equal func(a, b []byte) bool) {
	br := bufio.NewReader(r)

	var err error
	var oline, line []byte
	cnt := 1
	isLast := false
	for {
		line, err = br.ReadBytes('\n')
		line = bytes.TrimSuffix(line, []byte{'\n'})
		if err == io.EOF {
			isLast = true
		} else if err != nil {
			log.Printf("Can't read the %v line of %v file: %v", line, r, err)
		}
		if oline == nil {
			oline = line
			continue
		}
		if !equal(line, oline) {
			if count {
				fmt.Fprintf(w, "%d\t%s\n", cnt, oline)
				goto skip
			}
			if cnt > 1 && unique {
				goto skip
			}
			if cnt == 1 && duplicates {
				goto skip
			}
			fmt.Fprintf(w, "%s\n", oline)
		skip:
			oline = line
			cnt = 1
		} else {
			cnt++
		}
		if isLast {
			break
		}
	}
	if cnt == 1 && duplicates {
		return
	}
	if len(line) == 0 && cnt == 1 {
		return
	}
	if count {
		if len(line) == 0 {
			cnt--
		}
		fmt.Fprintf(w, "%d\t%s\n", cnt, line)
		return
	}
	fmt.Fprintf(w, "%s\n", line)
}

func run(stdin io.Reader, stdout io.Writer, unique, duplicates, count, ignoreCase bool, args []string) error {
	var eq func(a, b []byte) bool
	if ignoreCase {
		eq = bytes.EqualFold
	} else {
		eq = bytes.Equal
	}
	if len(args) == 0 {
		uniq(stdin, stdout, unique, duplicates, count, eq)
		return nil
	}
	for _, fn := range args {
		f, err := os.Open(fn)
		if err != nil {
			log.Printf("open %s: %v\n", fn, err)
			return err
		}
		uniq(f, stdout, unique, duplicates, count, eq)
		f.Close()
	}
	return nil
}

func main() {
	flag.Parse()
	if err := run(os.Stdin, os.Stdout, *unique, *duplicates, *count, *ignoreCase, flag.Args()); err != nil {
		log.Fatal(err)
	}
}
