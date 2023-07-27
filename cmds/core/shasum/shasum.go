// Copyright 2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bufio"
	"crypto/sha1"
	"crypto/sha256"
	"fmt"
	"hash"
	"io"
	"log"
	"os"

	"github.com/spf13/pflag"
)

var (
	algorithm = pflag.IntP("algorithm", "a", 1, "SHA algorithm, valid args are 1 and 256")
	help      = pflag.BoolP("help", "h", false, "Show this help and exit")
)
var usage = "Usage:\nshasum -a <algorithm> <File Name>"

func helpPrinter() {
	fmt.Println(usage)
	pflag.PrintDefaults()
}

// shaPrinter prints sha1/sha256 of given data. The
// value of algorithm is expected to be 1 for SHA1
// and 256 for SHA256
func shaGenerator(w io.Writer, r io.Reader, algo int) ([]byte, error) {
	var h hash.Hash
	switch algo {
	case 1:
		h = sha1.New()
	case 256:
		h = sha256.New()
	default:
		return nil, fmt.Errorf("invalid algorithm, only 1 or 256 are valid")
	}
	if _, err := io.Copy(h, r); err != nil {
		return nil, err
	}
	return h.Sum(nil), nil
}

func shasum(w io.Writer, r io.Reader, args ...string) error {
	if *help {
		helpPrinter()
		return nil
	}
	var hashbytes []byte
	var err error
	if len(args) == 0 {
		buf := bufio.NewReader(r)
		if hashbytes, err = shaGenerator(w, buf, *algorithm); err != nil {
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
		if hashbytes, err = shaGenerator(w, file, *algorithm); err != nil {
			return err
		}
		fmt.Fprintf(w, "%x %s\n", hashbytes, arg)
	}
	return nil
}

func main() {
	pflag.Parse()
	if err := shasum(os.Stdout, os.Stdin, pflag.Args()...); err != nil {
		log.Fatal(err)
	}
}
