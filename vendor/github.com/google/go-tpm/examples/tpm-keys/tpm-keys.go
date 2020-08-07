// +build !windows

// Copyright (c) 2016, Kevin Walsh.  All rights reserved.
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

// Package main implements a program to clear key handles from a TPM.
package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/google/go-tpm/tpm"
)

func main() {
	var tpmname = flag.String("tpm", "/dev/tpm0", "The path to the TPM device to use")
	var closekey = flag.Bool("close", false, "Close (unload) all existing key handles")
	flag.Parse()

	rwc, err := tpm.OpenTPM(*tpmname)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Couldn't open the TPM file %s: %s\n", *tpmname, err)
		return
	}

	handles, err := tpm.GetKeys(rwc)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Couldn't enumerate loaded TPM keys: %s\n", err)
		return
	}

	fmt.Printf("%d keys loaded in the TPM\n", len(handles))
	for i, h := range handles {
		fmt.Printf("  (%d) Key handle %d\n", i+1, h)
		if *closekey {
			if err = tpm.CloseKey(rwc, h); err != nil {
				fmt.Fprintf(os.Stderr, "Couldn't close TPM key handle %d\n", h)
			} else {
				fmt.Printf("    Closed handle %d\n", h)
			}
		}
	}

	return
}
