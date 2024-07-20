// Copyright 2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
)

var errCombine = fmt.Errorf("can't combine line and byte counts")

func run(stdout io.Writer, stderr io.Writer, bytes, count int, files ...string) error {
	if bytes > 0 && count > 0 {
		return errCombine
	}

	var printBytes bool
	var buffer []byte
	if bytes > 0 {
		printBytes = true
		buffer = make([]byte, bytes)
	}

	if count == 0 {
		count = 10
	}

	var newLineHeader bool

	for _, file := range files {
		f, err := os.Open(file)
		if err != nil {
			fmt.Fprintf(stderr, "head: %v", err)
			continue
		}
		if len(files) > 1 {
			if newLineHeader {
				fmt.Fprintf(stdout, "\n==> %s <==\n", filepath.Base(f.Name()))
			} else {
				fmt.Fprintf(stdout, "==> %s <==\n", filepath.Base(f.Name()))
				newLineHeader = true
			}
		}
		if printBytes {
			n, err := f.Read(buffer)
			if err != nil {
				fmt.Fprintf(stderr, "head: %v", err)
				continue
			}
			stdout.Write(buffer[:n])
		} else {
			var c int
			scanner := bufio.NewScanner(f)
			for scanner.Scan() {
				if c == count {
					break
				}
				fmt.Fprintln(stdout, scanner.Text())
				c++
			}
		}
	}

	return nil
}

func main() {
	var c = flag.Int("c", 0, "Print bytes of each of the specified files")
	var n = flag.Int("n", 0, "Print count lines of each of the specified files")

	flag.Parse()
	if err := run(os.Stdin, os.Stderr, *c, *n, flag.Args()...); err != nil {
		log.Fatalf("head: %v", err)
	}
}
