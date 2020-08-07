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
	"crypto/x509"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	"github.com/google/go-tpm/tpm"
)

func generateAction() {
	var tpmname = flag.String("tpm", "/dev/tpm0", "The path to the TPM device to use")
	var keyblobPath = flag.String("keyblob", "keyblob", "Output path of the generated keyblob")
	var pubKeyPath = flag.String("public-key", "publickey", "Output path of the generated keyblob's public key")
	var pcrsStr = flag.String("pcrs", "", "A comma-separated list of PCR numbers against which the generated key will be bound. If blank, it will not be bound to any PCR values.")
	flag.CommandLine.Parse(os.Args[2:])

	var pcrs []int
	if *pcrsStr != "" {
		for _, pcr := range strings.Split(*pcrsStr, ",") {
			pcrNum, err := strconv.Atoi(pcr)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Bad value in pcrs argument: %s\n", pcr)
				return
			}
			pcrs = append(pcrs, pcrNum)
		}
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

	keyblob, err := tpm.CreateWrapKey(rwc, srkAuth[:], usageAuth, migrationAuth, pcrs)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Couldn't make a new signing key: %s\n", err)
		return
	}
	fmt.Printf("Writing keyblob to %s\n", *keyblobPath)
	if err = ioutil.WriteFile(*keyblobPath, keyblob, 0644); err != nil {
		fmt.Fprintf(os.Stderr, "Error writing keyblob file: %s\n", err)
		return
	}

	pubKey, err := tpm.UnmarshalRSAPublicKey(keyblob)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not get public key: %s\n", err)
		return
	}

	pubKeyBytes, err := x509.MarshalPKIXPublicKey(pubKey)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not marshal public key: %s\n", err)
		return
	}
	fmt.Printf("Writing public key to %s\n", *pubKeyPath)
	if err = ioutil.WriteFile(*pubKeyPath, pubKeyBytes, 0644); err != nil {
		fmt.Fprintf(os.Stderr, "Error writing public key file: %s\n", err)
		return
	}
}
