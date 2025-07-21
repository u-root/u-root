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

func seq(w io.Writer, format string, separator string, widthEqual bool, args []string) error {
	var (
		stt   = 1.0
		stp   = 1.0
		end   float64
		width int
	)

	argv, argc := args, len(args)
	if argc < 1 || argc > 4 {
		return fmt.Errorf("mismatch n args; got %v, wants 1 >= n args >= 3", argc)
	}

	// loading step value if args is <start> <step> <end>
	if argc == 3 {
		_, err := fmt.Sscanf(argv[1], "%v", &stp)
		if stp-float64(int(stp)) > 0 && format == "%v" {
			d := len(fmt.Sprintf("%v", stp-float64(int(stp)))) - 2 // get the nums of y.xx decimal part
			format = fmt.Sprintf("%%.%df", d)
		}
		if stp == 0.0 {
			return errors.New("step value should be != 0")
		}

		if err != nil {
			return err
		}
	}

	if argc >= 2 { // cases: start + end || start + step + end
		if _, err := fmt.Sscanf(argv[0]+" "+argv[argc-1], "%v %v", &stt, &end); err != nil {
			return err
		}
	} else { // only <end>
		if _, err := fmt.Sscanf(argv[0], "%v", &end); err != nil {
			return err
		}
	}

	format = strings.Replace(format, "%", "%0*", 1) // support widthEqual
	if widthEqual {
		width = len(fmt.Sprintf(format, 0, end))
	}

	defer fmt.Fprint(w, "\n") // last char is always '\n'
	for stt <= end {
		fmt.Fprintf(w, format, width, stt)
		stt += stp
		if stt <= end { // print only between the values
			fmt.Fprint(w, separator)
		}
	}

	return nil
}

func main() {
	flag.Parse()
	if err := seq(os.Stdout, *format, *separator, *widthEqual, flag.Args()); err != nil {
		log.Fatalf("seq: %v", err)
	}
}
