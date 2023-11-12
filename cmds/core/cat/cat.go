// Copyright 2012-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// cat concatenates files and prints them to stdout.
//
// Synopsis:
//
//	cat [-u] [FILES]...
//
// Description:
//
//	If no files are specified, read from stdin.
//
// Options:
//
//	-u: ignored flag
package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
)

var _ = flag.Bool("u", false, "ignored")
var errCopy = fmt.Errorf("error concatenating stdin to stdout")

func cat(reader io.Reader, writer io.Writer) error {
	if _, err := io.Copy(writer, reader); err != nil {
		return errCopy
	}
	return nil
}

func run(stdin io.Reader, stdout io.Writer, args ...string) error {
	if len(args) == 0 {
		if err := cat(stdin, stdout); err != nil {
			return err
		}
	}
	for _, file := range args {
		if file == "-" {
			err := cat(stdin, stdout)
			if err != nil {
				return err
			}
			continue
		}
		f, err := os.Open(file)
		if err != nil {
			return err
		}
		if err := cat(f, stdout); err != nil {
			return fmt.Errorf("failed to concatenate file %s to given writer", f.Name())
		}
		f.Close()
	}
	return nil
}

func main() {
	flag.Parse()
	if err := run(os.Stdin, os.Stdout, flag.Args()...); err != nil {
		log.Fatalf("cat failed with: %v", err)
	}
}
