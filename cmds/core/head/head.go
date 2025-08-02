// Copyright 2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
)

var errCombine = fmt.Errorf("can't combine line and byte counts")

func run(stdin io.Reader, stdout, stderr io.Writer, bytes, count int, files ...string) error {
	if bytes > 0 && count > 0 {
		return errCombine
	}

	var printBytes bool
	var buffer []byte
	if bytes > 0 {
		printBytes = true
		buffer = make([]byte, 4096)
	}

	if count == 0 {
		count = 10
	}

	var newLineHeader bool
	var errs error

	handle := func(r io.Reader, name string) error {
		if len(files) > 1 {
			if newLineHeader {
				fmt.Fprintf(stdout, "\n==> %s <==\n", name)
			} else {
				fmt.Fprintf(stdout, "==> %s <==\n", name)
				newLineHeader = true
			}
		}
		if printBytes {
			c := bytes
			for {
				n, err := io.ReadFull(r, buffer)
				if err == io.EOF {
					break
				}
				if err != nil && err != io.ErrUnexpectedEOF {
					return err
				}

				stdout.Write(buffer[:min(c, n)])
				c -= n
				if c <= 0 {
					break
				}

				// handle the case when user request more bytes than
				// source have
				if err == io.ErrUnexpectedEOF {
					break
				}
			}
		} else {
			var c int
			scanner := bufio.NewScanner(r)
			for scanner.Scan() {
				fmt.Fprintln(stdout, scanner.Text())
				c++
				if c == count {
					break
				}
			}
		}
		return nil
	}

	// handle stdin
	if len(files) == 0 {
		return handle(stdin, "")
	}

	for _, file := range files {
		f, err := os.Open(file)
		if err != nil {
			errs = errors.Join(errs, fmt.Errorf("head: %w", err))
			continue
		}
		err = handle(f, f.Name())
		if err != nil {
			errs = errors.Join(errs, err)
		}
	}

	if errs != nil {
		fmt.Fprintf(stderr, "\n%v\n", errs)
	}
	return nil
}

func main() {
	c := flag.Int("c", 0, "Print bytes of each of the specified files")
	n := flag.Int("n", 0, "Print count lines of each of the specified files")

	flag.Parse()
	if err := run(os.Stdin, os.Stdout, os.Stderr, *c, *n, flag.Args()...); err != nil {
		log.Fatalf("head: %v", err)
	}
}
