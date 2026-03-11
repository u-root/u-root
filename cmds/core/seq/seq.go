// Copyright 2013-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Print a sequence of numbers.
//
// Synopsis:
//
//	seq [-f FORMAT] [-w] [-s SEPARATOR] [START [STEP [END]]]
//
// Examples:
//
//	% seq -s=' ' 3
//	1 2 3
//	% seq -s=' ' 2 4
//	2 3 4
//	% seq -s=' ' 3 2 7
//	3 5 7
//
// Options:
//
//	-f: use printf style floating-point FORMAT (default: %v)
//	-s: use STRING to separate numbers (default: \n)
//	-w: equalize width by padding with leading zeroes (default: false)
package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
)

const cmd = "seq [-f format] [-w] [-s separator] [start [step [end]]]"

var (
	format     = flag.String("f", "%v", "use printf style floating-point FORMAT")
	separator  = flag.String("s", "\n", "use STRING to separate numbers")
	widthEqual = flag.Bool("w", false, "equalize width by padding with leading zeroes")
)

func init() {
	defUsage := flag.Usage
	flag.Usage = func() {
		os.Args[0] = cmd
		defUsage()
	}
}

var (
	errUsage       = errors.New(cmd)
	errNegativeDec = errors.New("needs negative decrement")
	errPositiveInc = errors.New("needs positive increment")
	errZeroDec     = errors.New("zero decrement")
)

type s struct {
	first  float64
	incr   float64
	last   float64
	format string
}

func parse(format string, args []string) (s, error) {
	s1 := s{format: format, first: 1, incr: 1}
	argc := len(args)
	switch argc {
	case 3: // first, incr, last
		first, err := strconv.ParseFloat(args[0], 64)
		if err != nil {
			return s1, err
		}
		s1.first = first

		incr, err := strconv.ParseFloat(args[1], 64)
		if err != nil {
			return s1, err
		}
		s1.incr = incr

		if s1.incr-float64(int(s1.incr)) > 0 && format == "%v" {
			d := len(fmt.Sprintf("%v", s1.incr-float64(int(s1.incr)))) - 2 // get the nums of y.xx decimal part
			s1.format = fmt.Sprintf("%%.%df", d)
		}

		last, err := strconv.ParseFloat(args[2], 64)
		if err != nil {
			return s1, err
		}
		s1.last = last

		// handle wrong inc errors
		if s1.incr == 0 {
			return s1, errZeroDec
		}
		if s1.first > s1.last && s1.incr >= 0 {
			return s1, errNegativeDec
		}
		if s1.first < s1.last && s1.incr <= 0 {
			return s1, errPositiveInc
		}
	case 2: // first, last
		first, err := strconv.ParseFloat(args[0], 64)
		if err != nil {
			return s1, err
		}
		s1.first = first
		last, err := strconv.ParseFloat(args[1], 64)
		if err != nil {
			return s1, err
		}
		s1.last = last
	case 1: // last
		last, err := strconv.ParseFloat(args[0], 64)
		if err != nil {
			return s1, err
		}
		s1.last = last
	default:
		return s1, errUsage
	}

	return s1, nil
}

func seq(w io.Writer, format string, separator string, widthEqual bool, args []string) error {
	s1, err := parse(format, args)
	if err != nil {
		return err
	}

	// set default inc if needed
	if s1.first > s1.last && s1.incr == 1 {
		s1.incr = -1
	}

	format = strings.Replace(s1.format, "%", "%0*", 1) // support widthEqual
	var width int
	if widthEqual {
		width = len(fmt.Sprintf(format, 0, s1.last))
	}

	if s1.first < s1.last {
		for s1.first <= s1.last {
			fmt.Fprintf(w, format, width, s1.first)
			s1.first += s1.incr
			if s1.first <= s1.last { // print only between the values
				fmt.Fprint(w, separator)
			}
		}
	} else if s1.first > s1.last {
		for s1.first >= s1.last {
			fmt.Fprintf(w, format, width, s1.first)
			s1.first += s1.incr
			if s1.first >= s1.last { // print only between the values
				fmt.Fprint(w, separator)
			}
		}
	} else {
		fmt.Fprintf(w, format, width, s1.first)
	}

	fmt.Fprint(w, "\n") // last char is always '\n'
	return nil
}

func main() {
	flag.Parse()
	if err := seq(os.Stdout, *format, *separator, *widthEqual, flag.Args()); err != nil {
		log.Fatalf("seq: %v", err)
	}
}
