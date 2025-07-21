// Copyright 2013-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Freq reads the given files (default standard input) and prints histograms of the
// character frequencies. By default, freq counts each byte as a character; under
// the –r option it instead counts UTF sequences, that is, runes.
//
// Synopsis:
//
//	freq [-rdxoc] [FILES]...
//
// Description:
//
//	Each non–zero entry of the table is printed preceded by the byte value,
//	in decimal, octal, hex, and Unicode character (if printable). If any
//	options are given, the –d, –x, –o, –c flags specify a subset of value
//	formats: decimal, hex, octal, and character, respectively.
//
// Options:
//
//	–r: treat input as UTF-8
//	–d: print decimal value
//	–x: print hex value
//	–o: print octal value
//	–c: print character/UTF value
package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"slices"
	"unicode"
)

type params struct {
	utf bool
	dec bool
	hex bool
	oct bool
	chr bool
}

type cmd struct {
	stdin  io.Reader
	stdout io.Writer
	stderr io.Writer
	freq   map[rune]uint64
	args   []string
	params
}

func command(stdin io.Reader, stdout io.Writer, stderr io.Writer, p params, args ...string) *cmd {
	if !p.dec && !p.hex && !p.oct && !p.chr {
		p.dec, p.hex, p.oct, p.chr = true, true, true, true
	}

	return &cmd{
		stdin:  stdin,
		stdout: stdout,
		stderr: stderr,
		freq:   make(map[rune]uint64),
		params: p,
		args:   args,
	}
}

func (c *cmd) run() error {
	if len(c.args) > 0 {
		for _, v := range c.args {
			f, err := os.Open(v)
			if err != nil {
				return fmt.Errorf("open %s: %w", v, err)
			}
			c.doFreq(f)
			f.Close()
		}
	} else {
		c.doFreq(c.stdin)
	}

	keys := make([]rune, 0, len(c.freq))
	for k := range c.freq {
		keys = append(keys, k)
	}
	slices.Sort(keys)

	b := bufio.NewWriterSize(c.stdout, 8192*4)
	for _, r := range keys {
		if c.dec {
			fmt.Fprintf(b, "%3d ", r)
		}
		if c.oct {
			fmt.Fprintf(b, "%.3o ", r)
		}
		if c.hex {
			fmt.Fprintf(b, "%.2x ", r)
		}
		if c.chr {
			if unicode.IsPrint(r) {
				b.WriteRune(r)
				b.WriteString(" ")
			} else {
				b.WriteString("- ")
			}
		}
		fmt.Fprintf(b, "%8d\n", c.freq[r])
	}
	return b.Flush()
}

func (c *cmd) doFreq(f io.Reader) {
	b := bufio.NewReaderSize(f, 8192)

	var r rune
	var ch byte
	var err error
	if c.utf {
		for {
			r, _, err = b.ReadRune()
			if err != nil {
				if err != io.EOF {
					fmt.Fprintf(c.stderr, "error reading: %v", err)
				}
				return
			}
			c.freq[r]++
		}
	} else {
		for {
			ch, err = b.ReadByte()
			if err != nil {
				if err != io.EOF {
					fmt.Fprintf(c.stderr, "error reading: %v", err)
				}
				return
			}
			c.freq[rune(ch)]++
		}
	}
}

func main() {
	utf := flag.Bool("r", false, "treat input as UTF-8")
	dec := flag.Bool("d", false, "print decimal value")
	hex := flag.Bool("x", false, "print hexadecimal value")
	oct := flag.Bool("o", false, "print octal value")
	chr := flag.Bool("c", false, "print character/rune")
	flag.Parse()
	p := params{*utf, *dec, *hex, *oct, *chr}
	if err := command(os.Stdin, os.Stderr, os.Stdout, p, flag.Args()...).run(); err != nil {
		log.Fatal(err)
	}
}
