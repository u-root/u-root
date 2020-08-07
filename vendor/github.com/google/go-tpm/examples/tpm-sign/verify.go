// +build !windows

// Copyright (c) 2018, Ian Haken. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"crypto/rsa"
	"crypto/x509"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
)

func verifyAction() {
	var pubKeyPath = flag.String("public-key", "publickey", "Input path of public key file")
	var hashAlgArg = flag.String("hash", "SHA256", "Hash algorithm to use when verifying the signature")
	var signaturePath = flag.String("signature", "sig.data", "Input path of previously generated signature")
	var dataPath = flag.String("data", "", "Path to the data that was signed. If empty or omitted, the data will be read from stdin.")
	flag.CommandLine.Parse(os.Args[2:])

	hashAlg, ok := hashNames[*hashAlgArg]
	if !ok {
		fmt.Fprintf(os.Stderr, "Invalid hash algorithm: %s\n", *hashAlgArg)
		return
	}

	var err error
	var data []byte
	if *dataPath == "" {
		data, err = ioutil.ReadAll(os.Stdin)
	} else {
		data, err = ioutil.ReadFile(*dataPath)
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading input data: %s\n", err)
		return
	}

	signature, err := ioutil.ReadFile(*signaturePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading signature file: %s\n", err)
		return
	}

	pubKeyBytes, err := ioutil.ReadFile(*pubKeyPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading public key file: %s\n", err)
		return
	}
	pubKey, err := x509.ParsePKIXPublicKey(pubKeyBytes)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing public key: %s\n", err)
		return
	}
	rsaPubKey, ok := pubKey.(*rsa.PublicKey)
	if !ok {
		fmt.Fprintf(os.Stderr, "Expected public key to be an RSA key, but was %T\n", pubKey)
		return
	}

	hash := hashAlg.New()
	if _, err = hash.Write(data); err != nil {
		fmt.Fprintf(os.Stderr, "Error building hash of data: %s\n", err)
		return
	}
	hashed := hash.Sum(nil)
	if err = rsa.VerifyPKCS1v15(rsaPubKey, hashAlg, hashed[:], signature); err != nil {
		fmt.Fprintf(os.Stderr, "Error from verification: %s\n", err)
		return
	}
	fmt.Printf("Signature valid.\n")
}
