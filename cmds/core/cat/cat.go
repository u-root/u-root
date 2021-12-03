// Copyright 2012-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// cat concatenates files and prints them to stdout.
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
	"fmt"
	"io"
	"log"
	"os"
)

var _ = flag.Bool("u", false, "ignored")

func cat(reader io.Reader, writer io.Writer) error {
	if _, err := io.Copy(writer, reader); err != nil {
		return fmt.Errorf("error concatenating stdin to stdout: %v", err)
	}
	return nil
}

func run(args []string, stdin io.Reader, stdout io.Writer) error {
	if len(args) == 0 {
		if err := cat(stdin, stdout); err != nil {
			return fmt.Errorf("error concatenating stdin to stdout: %v", err)
		}
	}
	for _, file := range args {
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
	if err := run(os.Args[1:], os.Stdin, os.Stdout); err != nil {
		log.Fatalf("cat failed with: %v", err)
	}
}
