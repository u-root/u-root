// Copyright 2013-2023 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Wc counts lines, words, runes, syntactically–invalid UTF codes.
//
// Synopsis:
//
//	wc [OPTIONS...] [FILES]...
//
// Description:
//
//	Wc counts lines, words, runes, syntactically–invalid UTF codes and bytes
//	in the named files, or in the standard input if no file is named. A word
//	is a maximal string of characters delimited by spaces, tabs or newlines.
//	The count of runes includes invalid codes. If the optional argument is
//	present, just the specified counts (lines, words, runes, broken UTF
//	codes or bytes) are selected by the letters l, w, r, b, or c. Otherwise,
//	lines, words and bytes (–lwc) are reported.
//
// Options:
//
//	–l: count lines
//	–w: count words
//	–r: count runes
//	–b: count broken UTF codes
//	-c: count bytes
//
// Bugs:
//
//	This wc differs from Plan 9's wc somewhat in word count (BSD's wc differs
//	even more significantly):
//
//	$ unicode 0x0-0x10ffff | 9 wc -w
//	2228221
//	$ unicode 0x0-0x10ffff | gowc -w
//	2228198
//	$ unicode 0x0-0x10ffff | bsdwc -w
//	2293628
//
//	This wc differs from Plan 9's wc significantly in bad rune count:
//
//	$ unicode 0x0-0x10ffff | gowc -b
//	6144
//	$ unicode 0x0-0x10ffff | 9 wc -b
//	1966080
package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"unicode/utf8"
)

var (
	lines  = flag.Bool("l", false, "count lines")
	words  = flag.Bool("w", false, "count words")
	runes  = flag.Bool("r", false, "count runes")
	broken = flag.Bool("b", false, "count broken")
	chars  = flag.Bool("c", false, "count bytes (include partial UTF)")
)

type params struct {
	lines  bool
	words  bool
	runes  bool
	broken bool
	chars  bool
}

type cmd struct {
	stdin  io.Reader
	stdout io.Writer
	stderr io.Writer
	args   []string
	params
}

func command(stdin io.Reader, stdout io.Writer, stderr io.Writer, p params, args []string) *cmd {
	return &cmd{
		stdin:  stdin,
		stdout: stdout,
		stderr: stderr,
		params: p,
		args:   args,
	}
}

func (c *cmd) run() error {
	var totals cnt
	if !c.lines && !c.words && !c.runes && !c.broken && !c.chars {
		c.lines, c.words, c.chars = true, true, true
	}

	if len(c.args) == 0 {
		res := c.count(c.stdin, "")
		c.report(res, "")
		return nil
	}

	for _, v := range c.args {
		f, err := os.Open(v)
		if err != nil {
			fmt.Fprintf(c.stderr, "wc: %s: %v\n", v, err)
			continue
		}
		res := c.count(f, v)
		totals.lines += res.lines
		totals.words += res.words
		totals.runes += res.runes
		totals.badRunes += res.badRunes
		totals.chars += res.chars
		c.report(res, v)
	}
	if len(c.args) > 1 {
		c.report(totals, "total")
	}
	return nil
}

type cnt struct {
	lines, words, runes, badRunes, chars int64
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

func (c *cmd) count(in io.Reader, fname string) cnt {
	b := bufio.NewReaderSize(in, 8192)
	counted := false
	count := cnt{}
	for !counted {
		line, err := b.ReadBytes('\n')
		if err != nil {
			if err == io.EOF {
				counted = true
			} else {
				fmt.Fprintf(c.stderr, "wc: %s: %v\n", fname, err)
				return cnt{} // no partial counts; should perhaps quit altogether?
			}
		}
		if !counted {
			count.lines++
		}
		count.words += int64(len(bytes.Fields(line)))
		count.runes += int64(utf8.RuneCount(line))
		count.chars += int64(len(line))
		count.badRunes += invalidCount(line)
	}
	return count
}

func (c *cmd) report(count cnt, fname string) {
	fields := []string{}
	if c.lines {
		fields = append(fields, fmt.Sprintf("%d", count.lines))
	}
	if c.words {
		fields = append(fields, fmt.Sprintf("%d", count.words))
	}
	if c.runes {
		fields = append(fields, fmt.Sprintf("%d", count.runes))
	}
	if c.broken {
		fields = append(fields, fmt.Sprintf("%d", count.badRunes))
	}
	if c.chars {
		fields = append(fields, fmt.Sprintf("%d", count.chars))
	}
	if fname != "" {
		fields = append(fields, fname)
	}

	fmt.Fprintln(c.stdout, strings.Join(fields, " "))
}

func main() {
	flag.Parse()
	p := params{lines: *lines, words: *words, runes: *runes, broken: *broken, chars: *chars}
	if err := command(os.Stdin, os.Stdout, os.Stderr, p, flag.Args()).run(); err != nil {
		log.Fatal(err)
	}
}
