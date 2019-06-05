// Copyright 2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Strings finds printable strings.
//
// Synopsis:
//     strings OPTIONS [FILES]...
//
// Description:
//     Prints all sequences of `n` or more printable characters terminated by a
//     non-printable character (or EOF).
//
//     If no files are specified, read from stdin.
//
// Options:
//     -n number: the minimum string length (default is 4)
package main

import (
	"bufio"
	"io"
	"log"
	"os"

	flag "github.com/spf13/pflag"
)

var (
	n = flag.Int("n", 4, "the minimum string length")
)

func asciiIsPrint(char byte) bool {
	return char >= 32 && char <= 126
}

func stringsIO(r *bufio.Reader, w io.Writer) error {
	var o []byte
	for {
		b, err := r.ReadByte()
		if err == io.EOF {
			if len(o) >= *n {
				w.Write(o)
				w.Write([]byte{'\n'})
			}
			return nil
		}
		if err != nil {
			return err
		}
		if !asciiIsPrint(b) {
			if len(o) >= *n {
				w.Write(o)
				w.Write([]byte{'\n'})
			}
			o = o[:0]
			continue
		}
		// Prevent the buffer from growing indefinitely.
		if len(o) >= *n+1024 {
			w.Write(o[:1024])
			o = o[1024:]
		}
		o = append(o, b)
	}
}

func stringsFile(file string, w io.Writer) error {
	f, err := os.Open(file)
	if err != nil {
		return err
	}
	defer f.Close()

	// Buffer reduces number of syscalls.
	rb := bufio.NewReader(f)
	return stringsIO(rb, w)
}

func strings(files []string, w io.Writer) error {
	if len(files) == 0 {
		rb := bufio.NewReader(os.Stdin)
		if err := stringsIO(rb, w); err != nil {
			return err
		}
	}
	for _, file := range files {
		if err := stringsFile(file, w); err != nil {
			return err
		}
	}
	return nil
}

func main() {
	flag.Parse()

	if *n < 1 {
		log.Fatalf("strings: invalid minimum string length %v", *n)
	}

	// Buffer reduces number of syscalls.
	wb := bufio.NewWriter(os.Stdout)
	defer wb.Flush()

	if err := strings(flag.Args(), wb); err != nil {
		log.Fatalf("strings: %v", err)
	}
}
