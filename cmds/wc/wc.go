// Copyright 2013-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Wc counts lines, words, runes, syntactically–invalid UTF codes.
//
// Synopsis:
//     wc [OPTIONS...] [FILES]...
//
// Description:
//     Wc counts lines, words, runes, syntactically–invalid UTF codes and bytes
//     in the named files, or in the standard input if no file is named. A word
//     is a maximal string of characters delimited by spaces, tabs or newlines.
//     The count of runes includes invalid codes. If the optional argument is
//     present, just the specified counts (lines, words, runes, broken UTF
//     codes or bytes) are selected by the letters l, w, r, b, or c. Otherwise,
//     lines, words and bytes (–lwc) are reported.
//
// Options:
//     –l: count lines
//     –w: count words
//     –r: count runes
//     –b: count broken UTF codes
//     -c: count bytes
//
// Bugs:
//     This wc differs from Plan 9's wc somewhat in word count (BSD's wc differs
//     even more significantly):
//
//     $ unicode 0x0-0x10ffff | 9 wc -w
//     2228221
//     $ unicode 0x0-0x10ffff | gowc -w
//     2228198
//     $ unicode 0x0-0x10ffff | bsdwc -w
//     2293628
//
//     This wc differs from Plan 9's wc significantly in bad rune count:
//
//     $ unicode 0x0-0x10ffff | gowc -b
//     6144
//     $ unicode 0x0-0x10ffff | 9 wc -b
//     1966080
package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
	"unicode/utf8"
)

var lines = flag.Bool("l", false, "count lines")
var words = flag.Bool("w", false, "count words")
var runes = flag.Bool("r", false, "count runes")
var broken = flag.Bool("b", false, "count broken")
var chars = flag.Bool("c", false, "count bytes (include partial UTF)")

type cnt struct {
	nline, nword, nrune, nbadr, nchar int64
}

// A modified version of utf8.Valid()
func invalidCount(p []byte) (n int64) {
	i := 0
	for i < len(p) {
		if p[i] < utf8.RuneSelf {
			i++
		} else {
			_, size := utf8.DecodeRune(p[i:])
			if size == 1 {
				// All valid runes of size 1 (those
				// below RuneSelf) were handled above.
				// This muse be a RuneError.
				n++
			}
			i += size
		}
	}
	return
}

func count(in io.Reader, fname string) (c cnt) {
	b := bufio.NewReaderSize(in, 8192)

	counted := false
	for !counted {
		line, err := b.ReadBytes('\n')
		if err != nil {
			if err == io.EOF {
				counted = true
			} else {
				fmt.Fprintf(os.Stderr, "error %s: %v", fname, err)
				return cnt{} // no partial counts; should perhaps quit altogether?
			}
		}
		if !counted {
			c.nline++
		}
		c.nword += int64(len(bytes.Fields(line)))
		c.nrune += int64(utf8.RuneCount(line))
		c.nchar += int64(len(line))
		c.nbadr += invalidCount(line)
	}
	return
}

func report(c cnt, fname string) {
	fields := []string{}
	if *lines {
		fields = append(fields, fmt.Sprintf("%d", c.nline))
	}
	if *words {
		fields = append(fields, fmt.Sprintf("%d", c.nword))
	}
	if *runes {
		fields = append(fields, fmt.Sprintf("%d", c.nrune))
	}
	if *broken {
		fields = append(fields, fmt.Sprintf("%d", c.nbadr))
	}
	if *chars {
		fields = append(fields, fmt.Sprintf("%d", c.nchar))
	}
	if fname != "" {
		fields = append(fields, fmt.Sprintf("%s", fname))
	}

	fmt.Println(strings.Join(fields, " "))
}

func main() {
	var totals cnt

	flag.Parse()

	if !(*lines || *words || *runes || *broken || *chars) {
		*lines, *words, *chars = true, true, true
	}

	if flag.NArg() == 0 {
		cnt := count(os.Stdin, "")
		report(cnt, "")
		return
	}

	for _, v := range flag.Args() {
		f, err := os.Open(v)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error opening %s: %v\n", v, err)
			os.Exit(1)
		}
		cnt := count(f, v)
		totals.nline += cnt.nline
		totals.nword += cnt.nword
		totals.nrune += cnt.nrune
		totals.nbadr += cnt.nbadr
		totals.nchar += cnt.nchar
		report(cnt, v)
	}
	if flag.NArg() > 1 {
		report(totals, "total")
	}
}
