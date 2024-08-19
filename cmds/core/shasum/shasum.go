// Copyright 2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bufio"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"flag"
	"fmt"
	"hash"
	"io"
	"log"
	"os"
)

var algorithm int

var usage = "Usage:\nshasum -a <algorithm> <File Name>"

// shaPrinter prints sha1/sha256/sha512 of given data. The
// value of algorithm is expected to be 1 for SHA1
// 256 for SHA256
// and 512 for SHA512
func shaGenerator(w io.Writer, r io.Reader, algo int) ([]byte, error) {
	var h hash.Hash
	switch algo {
	case 1:
		h = sha1.New()
	case 256:
		h = sha256.New()
	case 512:
		h = sha512.New()
	default:
		return nil, fmt.Errorf("invalid algorithm, only 1, 256 or 512 are valid:%w", os.ErrInvalid)
	}
	if _, err := io.Copy(h, r); err != nil {
		return nil, err
	}
	return h.Sum(nil), nil
}

func shasum(w io.Writer, r io.Reader, args ...string) error {
	var hashbytes []byte
	var err error
	if len(args) == 0 {
		buf := bufio.NewReader(r)
		if hashbytes, err = shaGenerator(w, buf, algorithm); err != nil {
			return err
		}
		fmt.Fprintf(w, "%x -\n", hashbytes)
		return nil
	}
	for _, arg := range args {
		file, err := os.Open(arg)
		if err != nil {
			return err
		}
		defer file.Close()
		if hashbytes, err = shaGenerator(w, file, algorithm); err != nil {
			return err
		}
		fmt.Fprintf(w, "%x %s\n", hashbytes, arg)
	}
	return nil
}

func main() {
	flag.IntVar(&algorithm, "algorithm", 1, "SHA algorithm, valid args are 1, 256 and 512")
	flag.IntVar(&algorithm, "a", 1, "SHA algorithm, valid args are 1, 256 and 512")

	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), usage+"\n")
		flag.PrintDefaults()
	}

	flag.Parse()
	if err := shasum(os.Stdout, os.Stdin, flag.Args()...); err != nil {
		log.Fatal(err)
	}
}
