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

func stringsIO(r io.Reader, w io.Writer) error {
	// At least `n` bytes must be held in memory at a given time.
	rb := bufio.NewReaderSize(r, *n)

	// This processes the file one byte at a time. It might be inefficient,
	// but it works for now.
outerLoop:
	for {
		// Discard bytes from the buffer until the first `n` bytes are
		// all printable.
		peek, err := rb.Peek(*n)
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		for i := *n - 1; i >= 0; i-- {
			if !asciiIsPrint(peek[i]) {
				rb.Discard(i + 1)
				continue outerLoop
			}
		}

		// Write the first `n` bytes of the buffer.
		w.Write(peek)
		rb.Discard(*n)

		// Keep writing bytes until a non-printable byte is encountered.
		for {
			b, err := rb.ReadByte()
			if err == io.EOF {
				w.Write([]byte{'\n'})
				return nil
			}
			if err != nil {
				return err
			}
			if !asciiIsPrint(b) {
				w.Write([]byte{'\n'})
				break
			}
			w.Write([]byte{b})
		}
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
