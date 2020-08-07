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
	"crypto/sha1"
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/google/go-tpm/tpm"
)

func signAction() {
	var tpmname = flag.String("tpm", "/dev/tpm0", "The path to the TPM device to use")
	var keyblobPath = flag.String("keyblob", "keyblob", "Input path of the keyblob to use")
	var signaturePath = flag.String("signature", "sig.data", "Output path of the signature")
	var hashAlgArg = flag.String("hash", "SHA256", "Hash algorithm to use when generating the signature")
	var dataPath = flag.String("data", "", "Path to the data that will be signed. If empty or omitted, the data will be read from stdin.")
	flag.CommandLine.Parse(os.Args[2:])

	hashAlg, ok := hashNames[*hashAlgArg]
	if !ok {
		fmt.Fprintf(os.Stderr, "Invalid hash algorithm: %s\n", *hashAlgArg)
		return
	}

	rwc, err := tpm.OpenTPM(*tpmname)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Couldn't open the TPM file %s: %s\n", *tpmname, err)
		return
	}
	defer rwc.Close()

	// Compute the auth values as needed.
	var srkAuth [20]byte
	srkInput := os.Getenv(srkAuthEnvVar)
	if srkInput != "" {
		sa := sha1.Sum([]byte(srkInput))
		copy(srkAuth[:], sa[:])
	}

	var usageAuth [20]byte
	usageInput := os.Getenv(usageAuthEnvVar)
	if usageInput != "" {
		ua := sha1.Sum([]byte(usageInput))
		copy(usageAuth[:], ua[:])
	}

	var migrationAuth [20]byte
	migrationInput := os.Getenv(migrationAuthEnvVar)
	if migrationInput != "" {
		ma := sha1.Sum([]byte(migrationInput))
		copy(migrationAuth[:], ma[:])
	}

	keyblob, err := ioutil.ReadFile(*keyblobPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading keyblob file: %s\n", err)
		return
	}
	keyHandle, err := tpm.LoadKey2(rwc, keyblob, srkAuth[:])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not load keyblob: %s\n", err)
		return
	}
	defer tpm.CloseKey(rwc, keyHandle)

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

	hash := hashAlg.New()
	if _, err = hash.Write(data); err != nil {
		fmt.Fprintf(os.Stderr, "Error building hash of data: %s\n", err)
		return
	}
	hashed := hash.Sum(nil)
	signature, err := tpm.Sign(rwc, usageAuth[:], keyHandle, hashAlg, hashed[:])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not perform sign operation: %s\n", err)
		return
	}
	fmt.Printf("Writing signature to %s\n", *signaturePath)
	if err = ioutil.WriteFile(*signaturePath, signature, 0644); err != nil {
		fmt.Fprintf(os.Stderr, "Unable to write signature to file: %s\n", err)
		return
	}
}
