// Copyright 2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Prints files in hexadecimal.
//
// Synopsis:
//     hexdump [FILES]...
//
// Description:
//     Concatenate the input files into a single hexdump. If there are no
//     arguments, stdin is read.
package main

import (
	"encoding/hex"
	"flag"
	"io"
	"log"
	"os"
)

func main() {
	flag.Parse()

	var readers []io.Reader

	if flag.NArg() == 0 {
		readers = []io.Reader{os.Stdin}
	} else {
		readers = make([]io.Reader, 0, flag.NArg())

		for _, filename := range flag.Args() {
			f, err := os.Open(filename)
			if err != nil {
				log.Fatal(err)
			}
			defer f.Close()
			readers = append(readers, f)
		}
	}

	r := io.MultiReader(readers...)
	w := hex.Dumper(os.Stdout)
	defer w.Close()

	if _, err := io.Copy(w, r); err != nil {
		log.Fatal(err)
	}
}
