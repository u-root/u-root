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

func openFiles() ([]io.Reader, error) {
	readers := []io.Reader{os.Stdin}
	if flag.NArg() > 0 {
		readers = []io.Reader{}
		for _, filename := range flag.Args() {
			f, err := os.Open(filename)
			if err != nil {
				return nil, err
			}
			readers = append(readers, f)
		}
	}
	return readers, nil
}

func main() {
	flag.Parse()

	// Create a reader.
	readers, err := openFiles()
	if err != nil {
		log.Fatal(err)
	}
	r := io.MultiReader(readers...)

	// Dump hex to stdout.
	w := hex.Dumper(os.Stdout)
	if _, err := io.Copy(w, r); err != nil {
		log.Fatal(err)
	}
	w.Close()
}
