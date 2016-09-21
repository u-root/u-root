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

	armored    = flag.Bool("a", false, "signature is ASCII armored")
	sumfile    = flag.String("i", "", "checksum file")
	alg        = flag.String("alg", "", "algorithms to check")
	verbose    = flag.Bool("v", false, "verbose")
	debug      = func(string, ...interface{}) {}
	try, tried []string
)

func init() {
	for v := range algs {
		try = append(try, v)
	}
}

func one(n string, b []byte, sig string) bool {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Hash %v did not get linked in: %v", n, r)
		}
	}()
	debug("Check alg %v", n)
	checker := algs[n].New()
	checker.Write(b)
	r := checker.Sum([]byte{})
	tried = append(tried, n)

	// There has to be a better way.
	sumText := ""
	for _, v := range r {
		sumText += fmt.Sprintf("%02x", v)
	}
	debug("Compare to %v", sumText)
	if sumText == sig {
		return true
	}
	return false
}

func main() {
	flag.Parse()
	if flag.NArg() < 2 {
		log.Fatalf("Need at least a file to be validated and one public key")
	}

	if *verbose {
		debug = log.Printf
	}

	v, f := flag.Args()[0], flag.Args()[1]

	sigData, err := ioutil.ReadFile(v)
	if err != nil {
		log.Fatalf("%v", err)
	}

	if *alg != "" {
		try = []string{}
		for _, v := range strings.Split(*alg, ",") {
			try = append(try, v)
		}
	}

	log.Printf("Try %v", try)

	sig := strings.Split(string(sigData), " ")

	debug("Signature is %v len %v", sig[0], len(sig[0]))

	b, err := ioutil.ReadFile(f)
	if err != nil {
		log.Fatalf("%s: %v", f, err)
	}
	for i := range try {
		debug("Check %v", try[i])
		if one(try[i], b, sig[0]) {
			fmt.Printf("%v\n", try[i])
			os.Exit(0)
		}
		// Sometimes it's not a file in the standard format, but some binary thing.
		// Check that too.
		if one(try[i], b, string(sigData)) {
			fmt.Printf("%v\n", try[i])
			os.Exit(0)
		}
		debug("not ok")
	}
	log.Fatalf("No matches found for %v", tried)
}
