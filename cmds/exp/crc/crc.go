// Copyright 2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Prints crc checksum of a file.
//
// Synopsis:
//
//	crc OPTIONS [FILE]
//
// Description:
//
//	One of the crc types must be specified. If there is no file, stdin is
//	read.
//
// Options:
//
//	-f: CRC function to use. May be one of the following:
//	    crc32-ieee:       CRC-32 IEEE standard (default)
//	    crc32-castognoli: CRC-32 Castognoli standard
//	    crc32-koopman:    CRC-32 Koopman standard
//	    crc64-ecma:       CRC-64 ECMA standard
//	    crc64-iso:        CRC-64 ISO standard
package main

import (
	"flag"
	"fmt"
	"hash"
	"hash/crc32"
	"hash/crc64"
	"io"
	"log"
	"os"
	"strings"
)

func run(stdin io.Reader, stdout io.Writer, function string, args []string) error {
	functions := map[string]hash.Hash{
		"crc32-ieee":       crc32.New(crc32.MakeTable(crc32.IEEE)),
		"crc32-castognoli": crc32.New(crc32.MakeTable(crc32.Castagnoli)),
		"crc32-koopman":    crc32.New(crc32.MakeTable(crc32.Koopman)),
		"crc64-ecma":       crc64.New(crc64.MakeTable(crc64.ECMA)),
		"crc64-iso":        crc64.New(crc64.MakeTable(crc64.ISO)),
	}

	h, ok := functions[function]
	if !ok {
		var k []string
		for key := range functions {
			k = append(k, key)
		}
		return fmt.Errorf("%w: %q, expected one of: %s", os.ErrInvalid, function, strings.Join(k, " "))
	}

	var r io.Reader
	switch len(args) {
	case 0:
		r = stdin
	case 1:
		f, err := os.Open(args[0])
		if err != nil {
			return err
		}
		defer f.Close()
		r = f
	default:
		return fmt.Errorf("expected 0 or 1 positional args")
	}

	if _, err := io.Copy(h, r); err != nil {
		return err
	}

	_, err := fmt.Fprintf(stdout, "%x\n", h.Sum([]byte{}))
	return err
}

func main() {
	function := flag.String("f", "crc32-ieee", "CRC function")
	flag.Parse()

	if err := run(os.Stdin, os.Stdout, *function, flag.Args()); err != nil {
		log.Fatal(err)
	}
}
