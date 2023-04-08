// Copyright 2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// hexdump prints file content in hexadecimal.
//
// Synopsis:
//
//	hexdump [FILES]...
//
// Description:
//
//	Concatenate the input files into a single hexdump. If there are no
//	arguments, stdin is read.
package main

import (
	"encoding/hex"
	"flag"
	"io"
	"log"
	"os"
)

func hexdump(filenames []string, reader io.Reader, writer io.Writer) error {
	var readers []io.Reader

	if len(filenames) == 0 {
		readers = []io.Reader{reader}
	} else {
		readers = make([]io.Reader, 0, len(filenames))

		for _, filename := range filenames {
			f, err := os.Open(filename)
			if err != nil {
				return err
			}
			defer f.Close()
			readers = append(readers, f)
		}
	}

	r := io.MultiReader(readers...)
	w := hex.Dumper(writer)
	defer w.Close()

	if _, err := io.Copy(w, r); err != nil {
		return err
	}

	return nil
}

func main() {
	flag.Parse()
	if err := hexdump(flag.Args(), os.Stdin, os.Stdout); err != nil {
		log.Fatal(err)
	}
}
