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

func extendPcrAction() {
	var tpmname = flag.String("tpm", "/dev/tpm0", "The path to the TPM device to use")
	var pcrNum = flag.Int("pcr", 16, "PCR number to extend")
	var reset = flag.Bool("reset", false, "Reset the PCR rather than extending it")
	var dataPath = flag.String("data", "", "Path to the data that will be used to extend the PCR. If empty or omitted, the data will be read from stdin.")
	flag.CommandLine.Parse(os.Args[2:])

	rwc, err := tpm.OpenTPM(*tpmname)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Couldn't open the TPM file %s: %s\n", *tpmname, err)
		return
	}
	defer rwc.Close()

	if *reset {
		if err = tpm.PcrReset(rwc, []int{*pcrNum}); err != nil {
			fmt.Fprintf(os.Stderr, "Unable to reset PCR: %s\n", err)
			return
		}
	} else {
		var data []byte
		if *dataPath == "" {
			data, err = ioutil.ReadAll(os.Stdin)
		} else {
			data, err = ioutil.ReadFile(*dataPath)
		}
		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to read input: %s\n", err)
			return
		}
		if _, err = tpm.PcrExtend(rwc, uint32(*pcrNum), sha1.Sum(data)); err != nil {
			fmt.Fprintf(os.Stderr, "Error extending PCR: %s\n", err)
			return
		}
	}
}
