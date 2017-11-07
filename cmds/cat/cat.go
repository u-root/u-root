// Copyright 2012-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Cat concatenates files and prints to stdout.
//
// Synopsis:
//     cat [-u] [FILES]...
//
// Description:
//     If no files are specified, read from stdin.
//
// Options:
//     -u: ignored flag
package main

import (
	"flag"
	"io"
	"log"
	"os"
)

var (
	_ = flag.Bool("u", false, "ignored")
)

func catFile(w io.Writer, file string) error {
	f, err := os.Open(file)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = io.Copy(w, f)
	return err
}

func cat(w io.Writer, files []string) error {
	for _, name := range files {
		if err := catFile(w, name); err != nil {
			return err
		}
	}
	return nil
}

func main() {
	flag.Parse()

	if len(os.Args) == 1 {
		if _, err := io.Copy(os.Stdout, os.Stdin); err != nil {
			log.Fatalf("error concatenating stdin to stdout: %v", err)
		}
	}

	if err := cat(os.Stdout, flag.Args()); err != nil {
		log.Fatalf("cat: %v", err)
	}
}
