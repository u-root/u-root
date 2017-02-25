// Copyright 2013-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Convert and copy a file.
//
// Synopsis:
//     dd [OPTIONS...] [-inName FILE] [-outName FILE]
//
// Description:
//     dd is modeled after dd(1).
//
// Options:
//     -ibs n:   input block size (default=1)
//     -obs n:   output block size (default=1)
//     -bs n:    input and output block size (default=0)
//     -skip n:  skip n bytes before reading (default=0)
//     -seek n:  seek output when writing (default=0)
//     -conv s:  Convert the file on a specific way, like notrunc
//     -count n: max output of data to copy
//     -inName:  defaults to stdin
//     -outName: defaults to stdout
package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"strings"
)

var (
	ibs     = flag.Int64("ibs", 1, "Default input block size")
	obs     = flag.Int64("obs", 1, "Default output block size")
	bs      = flag.Int64("bs", 0, "Default input and output block size")
	skip    = flag.Int64("skip", 0, "skip n bytes before reading")
	seek    = flag.Int64("seek", 0, "seek output when writing")
	conv    = flag.String("conv", "", "Convert the file on a specific way, like notrunc")
	count   = flag.Int64("count", math.MaxInt64, "Max output of data to copy")
	inName  = flag.String("if", "", "Input file")
	outName = flag.String("of", "", "Output file")
)

func pass(r io.Reader, w io.Writer, ibs, obs int64, conv string) {
	b := make([]byte, ibs)
	for {
		n, err := io.ReadFull(r, b)
		if n == 0 || (err != nil && err != io.EOF) {
			break
		}
		var out []byte
		switch conv {
		case "ucase":
			out = []byte(strings.ToUpper(string(b)))
		case "lcase":
			out = []byte(strings.ToLower(string(b)))
		default:
			out = b
		}
		for n := int64(0); n < int64(len(out)); n += obs {
			if _, err := w.Write(out[n:obs]); err != nil {
				fmt.Fprintf(os.Stderr, "pass: %v\n", err)
				break
			}
		}
	}
}

func SplitArgs() []string {
	// EVERYTHING in dd follows x=y. So blindly split and convert sleep well
	arg := []string{}
	for _, v := range os.Args {
		l := strings.SplitN(v, "=", 2)
		// We only fix the exact case for x=y.
		if len(l) == 2 {
			l[0] = "-" + l[0]
			arg = append(arg, l...)
		} else {
			arg = append(arg, l...)
		}
	}
	return arg
}

func OpenFiles() (io.Reader, io.Writer) {
	i := io.ReaderAt(os.Stdin)
	o := io.Writer(os.Stdout)
	var err error

	if *inName != "" {
		if i, err = os.Open(*inName); err != nil {
			log.Fatal(err)
		}
	}
	if *outName != "" {
		if o, err = os.OpenFile(*outName, os.O_CREATE|os.O_WRONLY, 0644); err != nil {
			log.Fatal(err)
		}
	}

	// bs = both 'ibs' and 'obs' (IEEE Std 1003.1 - 2013)
	if *bs > 0 {
		*ibs = *bs
		*obs = *bs
	}

	var maxRead int64 = math.MaxInt64
	if *count != math.MaxInt64 {
		maxRead = *count * *ibs
	}

	return io.NewSectionReader(i, *skip**ibs, maxRead), o
}

// rather than, in essence, recreating all the apparatus of flag.xxxx with the if= bits,
// including dup checking, conversion, etc. we just convert the arguments and then
// run flag.Parse. Gross, but hey, it works.
func main() {
	os.Args = SplitArgs()
	flag.Parse()
	i, o := OpenFiles()
	pass(i, o, *obs, *obs, *conv)
}
