// Copyright 2013-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Uniq removes repeated lines.
//
// Synopsis:
//     uniq [OPTIONS...] [FILES]...
//
// Description:
//     Uniq copies the input file, or the standard input, to the standard
//     output, comparing adjacent lines. In the normal case, the second and
//     succeeding copies of repeated lines are removed. Repeated lines must be
//     adjacent in order to be found.
//
// Options:
//     –u:      Print unique lines.
//     –d:      Print (one copy of) duplicated lines.
//     –c:      Prefix a repetition count and a tab to each output line.
//              Implies –u and –d.
//     –f num:  The first num fields together with any blanks before each are
//              ignored. A field is defined as a string of non–space, non–tab
//              characters separated by tabs and spaces from its neighbors.
//     -cn num: The first num characters are ignored. Fields are skipped before
//              characters.
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

var uniques = flag.Bool("u", false, "print unique lines")
var duplicates = flag.Bool("d", false, "print one copy of duplicated lines")
var count = flag.Bool("c", false, "prefix a repetition count and a tab for each output line")

//var fnum = flag.Int("f", 0, "ignore num fields from beginning of line")
//var cnum = flag.Int("cn", 0, "ignore num characters from beginning of line")

func uniq(f *os.File) {
	br := bufio.NewReader(f)

	var err error
	var oline, line []byte
	cnt := 1
	isLast := false
	for {
		line, err = br.ReadBytes('\n')
		if err == io.EOF {
			isLast = true
		} else if err != nil {
			log.Printf("Can't read the %v line of %v file: %v", line, f, err)
		}
		if oline == nil {
			oline = line
			continue
		}
		if !bytes.Equal(line, oline) {
			if *count {
				fmt.Printf("%d\t%s", cnt, oline)
				goto skip
			}
			if cnt > 1 && *uniques {
				goto skip
			}
			if cnt == 1 && *duplicates {
				goto skip
			}
			fmt.Printf("%s", oline)
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
	if cnt > 1 && *uniques {
		return
	}
	if cnt == 1 && *duplicates {
		return
	}
	fmt.Printf("%s", line)
}

func main() {
	flag.Parse()

	if flag.NArg() > 0 {
		for _, fn := range flag.Args() {
			f, err := os.Open(fn)
			if err != nil {
				log.Printf("open %s: %v\n", fn, err)
				os.Exit(1)
			}
			uniq(f)
			f.Close()
		}
	} else {
		uniq(os.Stdin)
	}
}
