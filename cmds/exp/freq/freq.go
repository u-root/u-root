// Copyright 2013-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Freq reads the given files (default standard input) and prints histograms of the
// character frequencies. By default, freq counts each byte as a character; under
// the –r option it instead counts UTF sequences, that is, runes.
//
// Synopsis:
//     freq [-rdxoc] [FILES]...
//
// Description:
//     Each non–zero entry of the table is printed preceded by the byte value,
//     in decimal, octal, hex, and Unicode character (if printable). If any
//     options are given, the –d, –x, –o, –c flags specify a subset of value
//     formats: decimal, hex, octal, and character, respectively.
//
// Options:
//     –r: treat input as UTF-8
//     –d: print decimal value
//     –x: print hex value
//     –o: print octal value
//     –c: print character/UTF value
package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"unicode/utf8"
)

var utf = flag.Bool("r", false, "treat input as UTF-8")
var dec = flag.Bool("d", false, "print decimal value")
var hex = flag.Bool("x", false, "print hexadecimal value")
var oct = flag.Bool("o", false, "print octal value")
var chr = flag.Bool("c", false, "print character/rune")

var freq [utf8.MaxRune + 1]uint64

func doFreq(f *os.File) {
	b := bufio.NewReaderSize(f, 8192)

	var r rune
	var c byte
	var err error
	if *utf {
		for {
			r, _, err = b.ReadRune()
			if err != nil {
				if err != io.EOF {
					fmt.Fprintf(os.Stderr, "error reading: %v", err)
				}
				return
			}
			freq[r]++
		}
	} else {
		for {
			c, err = b.ReadByte()
			if err != nil {
				if err != io.EOF {
					fmt.Fprintf(os.Stderr, "error reading: %v", err)
				}
				return
			}
			freq[c]++
		}
	}
}

func main() {
	flag.Parse()

	if flag.NArg() > 0 {
		for _, v := range flag.Args() {
			f, err := os.Open(v)
			if err != nil {
				fmt.Fprintf(os.Stderr, "open %s: %v", v, err)
				os.Exit(1)
			}
			doFreq(f)
			f.Close()
		}
	} else {
		doFreq(os.Stdin)
	}

	if !(*dec || *hex || *oct || *chr) {
		*dec, *hex, *oct, *chr = true, true, true, true
	}

	b := bufio.NewWriterSize(os.Stdout, 8192*4)
	for i, v := range freq {
		if v == 0 {
			continue
		}

		if *dec {
			fmt.Fprintf(b, "%3d ", i)
		}
		if *oct {
			fmt.Fprintf(b, "%.3o ", i)
		}
		if *hex {
			fmt.Fprintf(b, "%.2x ", i)
		}
		if *chr {
			if i <= 0x20 || (i >= 0x7f && i < 0xa0) || (i > 0xff && !(*utf)) {
				b.WriteString("- ")
			} else {
				b.WriteRune(rune(i))
				b.WriteString(" ")
			}
		}
		fmt.Fprintf(b, "%8d\n", v)
	}
	b.Flush()
}
