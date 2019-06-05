// Copyright 2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"crypto/sha1"
	"crypto/sha256"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/spf13/pflag"
)

func helpPrinter() {

	fmt.Printf("Usage:\nshasum -a <algorithm> <File Name>\n")
	pflag.PrintDefaults()
	os.Exit(0)
}

func versionPrinter() {
	fmt.Println("shasum utility, URoot Version.")
	os.Exit(0)
}

func getInput(fileName string) (input []byte, err error) {

	if fileName != "" {
		return ioutil.ReadFile(fileName)
	}
	return ioutil.ReadAll(os.Stdin)
}

//
// shaPrinter prints sha1/sha256 of given data. The
// value of algorithm is expected to be 1 for SHA1
// and 256 for SHA256
//
func shaPrinter(algorithm int, data []byte) string {
	var sha string
	if algorithm == 256 {
		sha = fmt.Sprintf("%x", sha256.Sum256(data))
	} else if algorithm == 1 {
		sha = fmt.Sprintf("%x", sha1.Sum(data))
	} else {
		fmt.Fprintf(os.Stderr, "Invalid algorithm")
		return ""
	}
	return sha
}

func main() {

	var (
		algorithm int
		help      bool
		version   bool
	)
	cliArgs := ""
	pflag.IntVarP(&algorithm, "algorithm", "a", 1, "SHA algorithm, valid args are 1 and 256")
	pflag.BoolVarP(&help, "help", "h", false, "Show this help and exit")
	pflag.BoolVarP(&version, "version", "v", false, "Print Version")
	pflag.Parse()

	if help {
		helpPrinter()
	}

	if version {
		versionPrinter()
	}
	if len(pflag.Args()) == 1 {
		cliArgs = pflag.Args()[0]
	}
	input, err := getInput(cliArgs)
	if err != nil {
		fmt.Println("Error getting input.")
		os.Exit(-1)
	}
	fmt.Printf("%s ", shaPrinter(algorithm, input))
	if cliArgs == "" {
		fmt.Printf(" -\n")
	} else {
		fmt.Printf(" %s\n", cliArgs)
	}
	os.Exit(0)
}
