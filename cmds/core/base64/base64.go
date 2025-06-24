// Copyright 2021 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// base64 - encode and decode base64 from stdin or file to stdout

// Synopsis:
//     base64 [-d] [FILE]

// Description:
//    Encode or decode a file to or from base64 encoding.
//    -d   decode data (default is to encode)
//    For stdin, on standard Unix systems, you can use /dev/stdin

package main

import (
	"encoding/base64"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
)

var (
	decode      = flag.Bool("d", false, "Decode")
	errBadUsage = errors.New("usage: base64 [-d] [file]")
)

func do(r io.Reader, w io.Writer, decode bool) error {
	if decode {
		r = base64.NewDecoder(base64.RawStdEncoding, r)
		if _, err := io.Copy(w, r); err != nil {
			return fmt.Errorf("base64: error decoding %w", err)
		}
		return nil
	}

	// WriteCloser is important here, from NewEncoder documentation:
	// when finished writing, the caller must Close the returned encoder
	// to flush any partially written blocks.
	wc := base64.NewEncoder(base64.RawStdEncoding, w)
	defer wc.Close()
	if _, err := io.Copy(wc, r); err != nil {
		return fmt.Errorf("base64: error encoding %w", err)
	}
	return nil
}

// run runs the base64 command. Why use ...string?
// makes testing a tad easier (so we don't have an if in main()
// allows us, should we wish, in future, to go with using
// names[1] as out. base64 commands are very nonstandard.
func run(stdin io.Reader, stdout io.Writer, decode bool, names ...string) error {
	switch len(names) {
	case 0:
	case 1:
		f, err := os.Open(names[0])
		if err != nil {
			return err
		}
		stdin = f
	default:
		return errBadUsage
	}

	return do(stdin, stdout, decode)
}

func main() {
	flag.Parse()
	if err := run(os.Stdin, os.Stdout, *decode, flag.Args()...); err != nil {
		log.Fatalf("base64: %v", err)
	}
}
