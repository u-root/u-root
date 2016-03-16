// Copyright 2016 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//This program validates a file by verifying a checksum file and a signature file
//Exit status: 0-OK, 1-Any error, 2-Bad signature, 3-Bad checksum
package main

import (
	"crypto"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	// TODO _ "golang.org/x/crypto/openpgp"
	// TODO _ "golang.org/x/crypto/md4"
	_ "crypto/md5"
	_ "crypto/sha1"
	_ "crypto/sha256"
	_ "crypto/sha512"

	//	_ "golang.org/x/crypto/ripemd160"
	//	_ "golang.org/x/crypto/sha3"
	_ "crypto/sha512"
)

var (
	algs = map[string]crypto.Hash{
		"MD4":       crypto.MD4,
		"MD5":       crypto.MD5,
		"SHA1":      crypto.SHA1,
		"SHA224":    crypto.SHA224,
		"SHA256":    crypto.SHA256,
		"SHA384":    crypto.SHA384,
		"SHA512":    crypto.SHA512,
		"RIPEMD160": crypto.RIPEMD160,
		"SHA3_224":  crypto.SHA3_224,
		"SHA3_256":  crypto.SHA3_256,
		"SHA3_384":  crypto.SHA3_384,
		"SHA3_512":  crypto.SHA3_512,
	}

	armored = flag.Bool("a", false, "signature is ASCII armored")
	sumfile = flag.String("i", "", "checksum file")
	alg     = flag.String("alg", "MD5", "md5sum")
	verbose = flag.Bool("v", false, "verbose")
	debug   = func(string, ...interface{}) {}
)

func main() {
	flag.Parse()
	if flag.NArg() < 2 {
		log.Fatalf("Need at least a file to be validated and one public key")
	}
	if *verbose {
		debug = log.Printf
	}

	// TODO: read in the file with the validation.
	// The second args will be flag.Args()[1], of course!
	f, v := flag.Args()[0], flag.Args()[1]

	sigData, err := ioutil.ReadFile(v)
	if err != nil {
		log.Fatalf("%v", err)
	}

	sig := strings.Split(string(sigData), " ")
	debug("Signature is %v len %v", sig[0], len(sig[0]))

	b, err := ioutil.ReadFile(f)
	if err != nil {
		log.Fatalf("%s: %v", f, err)
	}
	for _, h := range strings.Split(*alg, ",") {
		debug("Check %v", h)
		h, ok := algs[h]
		if !ok {
			debug("%s is not in %v", h, algs)
		}

		checker := h.New()
		checker.Write(b)
		r := checker.Sum([]byte{})

		// There has to be a better way.
		sumText := ""
		for _, v := range r {
			sumText += fmt.Sprintf("%02x", v)
		}
		debug("Compare to %v", sumText)
		if sumText == sig[0] {
			debug("ok")
			os.Exit(0)
		}
		debug("not ok")
	}
	log.Fatalf("No matches found for *alg")
}
