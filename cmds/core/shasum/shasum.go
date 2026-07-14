// Copyright 2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// shasum computes SHA checksums of files.
//
// Synopsis:
//
//	shasum -a <algorithm> <File Name>
//
// Description:
//
//	shasum computes SHA checksums of files using the specified algorithm.
//	If no files are specified, read from stdin.
//
// Options:
//
//	-a, -algorithm: SHA algorithm, valid args are 1, 224, 256, 384, 512, 512224 and 512256
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

	"github.com/u-root/u-root/pkg/uroot/unixflag"
)

// shaGenerator generates SHA hash of given data. The
// value of algorithm is expected to be 1 for SHA1
// 256 for SHA256
// and 512 for SHA512
func shaGenerator(r io.Reader, algo int) ([]byte, error) {
	var h hash.Hash
	switch algo {
	case 1:
		h = sha1.New()
	case 224:
		h = sha256.New224()
	case 256:
		h = sha256.New()
	case 384:
		h = sha512.New384()
	case 512:
		h = sha512.New()
	case 512224:
		h = sha512.New512_224()
	case 512256:
		h = sha512.New512_256()
	default:
		return nil, fmt.Errorf("invalid algorithm, only 1, 224, 256, 384, 512, 512224 and 512256 are valid:%w", os.ErrInvalid)
	}
	if _, err := io.Copy(h, r); err != nil {
		return nil, err
	}
	return h.Sum(nil), nil
}

func run(w io.Writer, r io.Reader, args []string) error {
	f := flag.NewFlagSet("shasum", flag.ExitOnError)

	var algorithm int
	f.IntVar(&algorithm, "algorithm", 1, "SHA algorithm, valid args are 1, 224, 256, 384, 512, 512224 and 512256")
	f.IntVar(&algorithm, "a", 1, "SHA algorithm, valid args are 1, 224, 256, 384, 512, 512224 and 512256")

	f.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "%s\n", "Usage:\nshasum -a <algorithm> <File Name>")
		flag.PrintDefaults()
	}

	f.Parse(unixflag.ArgsToGoArgs(args))
	args = f.Args()

	var hashbytes []byte
	var err error
	if len(args) == 0 {
		buf := bufio.NewReader(r)
		if hashbytes, err = shaGenerator(buf, algorithm); err != nil {
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
		if hashbytes, err = shaGenerator(file, algorithm); err != nil {
			return err
		}
		_ = file.Close()
		fmt.Fprintf(w, "%x %s\n", hashbytes, arg)
	}
	return nil
}

func main() {
	if err := run(os.Stdout, os.Stdin, os.Args[1:]); err != nil {
		log.Fatal(err)
	}
}
