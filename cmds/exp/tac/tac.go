// Copyright 2012-2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// tac concatenates files and prints to stdout in reverse order,
// file by file
//
// Synopsis:
//
//	tac <file...>
//
// Description:
//
// Options:
package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"sync"
)

const ReadSize int64 = 4096

var errStdin = fmt.Errorf("can't reverse lines from stdin; can't seek")

type ReadAtSeeker interface {
	io.ReaderAt
	io.Seeker
}

func tacOne(w io.Writer, r ReadAtSeeker) error {
	var b [ReadSize]byte
	// Get current EOF. While the file may be growing, there's
	// only so much we can do.
	loc, err := r.Seek(0, io.SeekEnd)
	if err != nil {
		return err
	}
	var wg sync.WaitGroup
	wg.Add(1)
	c := make(chan byte)
	go func(r <-chan byte, w io.Writer) {
		defer wg.Done()
		line := string(<-r)
		for c := range r {
			if c == '\n' {
				if _, err := w.Write([]byte(line)); err != nil {
					log.Fatal(err)
				}
				line = ""
			}
			line = string(c) + line
		}
		if _, err := w.Write([]byte(line)); err != nil {
			log.Fatal(err)
		}
	}(c, w)

	for loc > 0 {
		n := min(loc, ReadSize)

		amt, err := r.ReadAt(b[:n], loc-int64(n))
		if err != nil && err != io.EOF {
			return err
		}
		loc -= int64(amt)
		for i := range b[:amt] {
			o := amt - i - 1
			c <- b[o]
		}
	}
	close(c)
	wg.Wait()
	return nil
}

func tac(w io.Writer, files []string) error {
	if len(files) == 0 {
		return errStdin
	}
	for _, name := range files {
		f, err := os.Open(name)
		if err != nil {
			return err
		}
		err = tacOne(w, f)
		f.Close() // Don't defer, you might get EMFILE for no good reason.
		if err != nil {
			return err
		}

	}
	return nil
}

func main() {
	flag.Parse()
	if err := tac(os.Stdout, flag.Args()); err != nil {
		log.Fatalf("tac: %v", err)
	}
}
