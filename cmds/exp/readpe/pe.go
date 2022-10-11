// Copyright 2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Dump the headers of a PE file.
//
// Synopsis:
//
//	pe [FILENAME]
//
// Description:
//
//	Windows and EFI executables are in the portable executable (PE) format.
//	This command prints the headers in a JSON format.
package main

import (
	"debug/pe"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
)

func main() {
	// Parse flags
	flag.Parse()
	var (
		f   *pe.File
		err error
	)
	switch flag.NArg() {
	case 0:
		f, err = pe.NewFile(os.Stdin)
	case 1:
		filename := flag.Arg(0)
		f, err = pe.Open(filename)
	default:
		log.Fatal("Usage: pe [FILENAME]")
	}
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	// Convert to JSON
	j, err := json.MarshalIndent(f, "", "\t")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(j))
}
