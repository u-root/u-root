// Copyright 2021 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// base64 - encode and decode base64 from stdin or file to stdout

// Synopsis:
//     base64 [-d] [FILE]

// Description:
//    Encode or decode a file to or from base64 encoding.
//    -d   decode data (default is to encode)

package main

import (
	"encoding/base64"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
)

var decode = flag.Bool("d", false, "Decode")

func do(r io.Reader, w io.Writer, decode bool) error {
	op := "decoding"
	if decode {
		r = base64.NewDecoder(base64.RawStdEncoding, r)
	} else {
		op = "encoding"
		w = base64.NewEncoder(base64.RawStdEncoding, w)
	}

	if _, err := io.Copy(w, r); err != nil {
		return fmt.Errorf("error %s the data: %v", op, err)
	}
	return nil
}

func run(name string, stdin io.Reader, stdout io.Writer, decode bool) error {
	if name != "-" && len(name) > 0 {
		f, err := os.Open(name)
		if err != nil {
			return err
		}
		stdin = f
	}

	return do(stdin, stdout, decode)
}

func main() {
	flag.Parse()
	var name string
	if len(flag.Args()) > 1 {
		name = flag.Args()[0]
	}

	if err := run(name, os.Stdin, os.Stdout, *decode); err != nil {
		log.Fatalf("base64: %v", err)
	}
}
