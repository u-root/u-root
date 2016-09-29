// Copyright 2016 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// gpgv validates a signature against a file.
// It prints "OK\n" to stdout if the check succeeds and exits with 0.
// It prints an error message and exits with non-0 otherwise.
package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	"golang.org/x/crypto/openpgp"
)

var (
	verbose = flag.Bool("v", false, "verbose")
	debug   = func(string, ...interface{}) {}
)

func main() {
	var k, s, f io.Reader
	var err error
	var check = openpgp.CheckDetachedSignature

	flag.Parse()
	if flag.NArg() < 3 {
		log.Fatalf("usage: gpgv [-v] <keyring file> <signature file> <file to be verified>")
	}

	if *verbose {
		debug = log.Printf
	}

	kn, sn, fn := flag.Args()[0], flag.Args()[1], flag.Args()[2]

	if k, err = os.Open(kn); err != nil {
		log.Fatalf("Can't open key file: %v", err)
	}
	kr, err := openpgp.ReadKeyRing(k)
	if err != nil {
		log.Printf("ReadKeyRing: %v, trying Armored", err)
		if k, err = os.Open(kn); err != nil {
			log.Fatalf("reopen KeyRing: %v", err)
		}
		kr, err = openpgp.ReadArmoredKeyRing(k)
		if err != nil {
			log.Fatalf("ReadArmoredKeyRing: %v", err)
		}
		check = openpgp.CheckArmoredDetachedSignature
	}

	if s, err = os.Open(sn); err != nil {
		log.Fatalf("Can't open signature file: %v", err)
	}
	if f, err = os.Open(fn); err != nil {
		log.Fatalf("Can't open data file: %v", err)
	}

	sig, err := check(kr, f, s)
	if err != nil {
		log.Fatalf("%v", err)
	}
	debug("Signature: '%v'", sig)
	fmt.Printf("OK\n")
}
