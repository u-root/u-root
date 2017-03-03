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
	"bytes"
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
	var b = &bytes.Buffer{}
	for {
		b.Reset()
		n, err := io.CopyN(b, r, ibs)
		if n == 0 {
			break
		}
		if err != nil && err != io.EOF {
			log.Fatalf("Reading from input: %v", err)
		}
		var out *bytes.Buffer
		switch conv {
		case "ucase":
			out = bytes.NewBuffer(bytes.ToUpper(b.Bytes()))
		case "lcase":
			out = bytes.NewBuffer(bytes.ToLower(b.Bytes()))
		default:
			out = bytes.NewBuffer(b.Bytes())
		}
		for n := int64(0); n < int64(out.Len()); {
			amt, err := io.CopyN(w, out, obs)
			if err != nil {
				fmt.Fprintf(os.Stderr, "pass: %v\n", err)
				break
			}
			n += int64(amt)
		}
		if err == io.EOF {
			break
		}
	}
}

func splitArgs() []string {
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

func seekOrRead(r io.ReadSeeker, bs, cnt int64) {
	if bs*cnt == 0 {
		return
	}
	// I tried to make a NewSectionReader but, sadly, most OSes
	// get upset when you pread even if it does not involve a seek.
	// I wish I were making that up.
	if _, err := r.Seek(1, int(cnt*bs)); err == nil {
		return
	}

	// the only choice is to eat it.
	var b = &bytes.Buffer{}

	for i := int64(0); i < cnt*bs; {
		amt, err := io.CopyN(b, r, bs)
		if err != nil {
			return
		}
		i += amt
	}
}

func openFiles() (io.Reader, io.Writer) {
	i := io.ReadSeeker(os.Stdin)
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

	// I tried to make a NewSectionReader but, sadly, most OSes
	// get upset when you pread even if it does not involve a seek.
	// I wish I were making that up.
	seekOrRead(i, *ibs, *skip)
	return io.LimitReader(i, maxRead), o
}

// rather than, in essence, recreating all the apparatus of flag.xxxx with the if= bits,
// including dup checking, conversion, etc. we just convert the arguments and then
// run flag.Parse. Gross, but hey, it works.
func main() {
	os.Args = splitArgs()
	flag.Parse()
	i, o := openFiles()
	pass(i, o, *obs, *obs, *conv)
}
