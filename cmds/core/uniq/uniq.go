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
package main

// TODO: implement these flags:
//	–f num:  The first num fields together with any blanks before each are
//	         ignored. A field is defined as a string of non–space, non–tab
//	         characters separated by tabs and spaces from its neighbors.
//	-s num: The first num characters are ignored. Fields are skipped before
//	         characters.

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

func shouldPrint(unique, duplicates bool, cnt int) bool {
	if unique && cnt > 1 || duplicates && cnt == 1 {
		return false
	}
	return true
}

func printLine(w io.Writer, line []byte, count bool, cnt int) {
	if count {
		_, _ = fmt.Fprintf(w, "%d\t%s\n", cnt, line)
	} else {
		_, _ = fmt.Fprintf(w, "%s\n", line)
	}
}

func uniq(r io.Reader, w io.Writer, unique, duplicates, count bool, equal func(a, b []byte) bool) error {
	bs := bufio.NewScanner(r)

	if !bs.Scan() {
		return bs.Err()
	}
	prevLine := bytes.Clone(bs.Bytes())
	cnt := 1

	for bs.Scan() {
		line := bytes.Clone(bs.Bytes())
		if equal(line, prevLine) {
			cnt++
		} else {
			if shouldPrint(unique, duplicates, cnt) {
				printLine(w, prevLine, count, cnt)
			}
			cnt = 1
			prevLine = line
		}
	}
	if err := bs.Err(); err != nil {
		return err
	}

	if shouldPrint(unique, duplicates, cnt) {
		printLine(w, prevLine, count, cnt)
	}
	return nil
}

func run(stdin io.Reader, stdout io.Writer, unique, duplicates, count, ignoreCase bool, args []string) error {
	eq := bytes.Equal
	if ignoreCase {
		eq = bytes.EqualFold
	}
	if len(args) == 0 {
		return uniq(stdin, stdout, unique, duplicates, count, eq)
	}
	for _, fn := range args {
		f, err := os.Open(fn)
		if err != nil {
			return err
		}
		if err := uniq(f, stdout, unique, duplicates, count, eq); err != nil {
			return err
		}
		_ = f.Close()
	}
	return nil
}

func main() {
	flag.Parse()
	if err := run(os.Stdin, os.Stdout, *unique, *duplicates, *count, *ignoreCase, flag.Args()); err != nil {
		log.Fatal(err)
	}
}
