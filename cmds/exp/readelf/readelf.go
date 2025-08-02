// Copyright 2023 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Dump the headers of an ELF file, from stdin or a set of files
//
// Synopsis:
//
//	readelf [file...]
//
// Description:
//
//	read the ELF header and print information.
package main

import (
	"debug/elf"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
)

func main() {
	flag.Parse()
	if err := run(os.Stdin, os.Stdout, flag.Args()...); err != nil {
		log.Fatal(err)
	}
}

func run(in io.ReaderAt, out io.Writer, args ...string) error {
	var (
		files []*elf.File
		err   error
	)
	switch len(args) {
	case 0:
		files = make([]*elf.File, 1)
		files[0], err = elf.NewFile(in)
	default:
		files = make([]*elf.File, len(args))
		for i, n := range args {
			f, e := elf.Open(n)
			if e != nil {
				// Until CI gets to go1.21, just break
				// on an error.
				// err = errors.Join(err, e)
				// continue
				return e
			}
			files[i] = f
			defer f.Close()
		}
	}

	if err != nil {
		return err
	}

	for _, f := range files {
		j, err := json.MarshalIndent(f, "", "\t")
		if err != nil {
			return err
		}
		fmt.Fprintf(out, "%s\n", j)
	}
	return nil
}
