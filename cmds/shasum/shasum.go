// Copyright 2016-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"crypto/sha1"
	"crypto/sha256"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
)

func helpPrinter() {

	fmt.Printf("Usage:\nshasum -a <algorithm> <File Name>\n")
	flag.PrintDefaults()
	os.Exit(0)
}

func versionPrinter() {
	fmt.Println("shasum utility, URoot Version.")
	os.Exit(0)
}

func GetInput(fileName string) (input []byte, err error) {

	if fileName != "" {
		return ioutil.ReadFile(fileName)
	}
	return ioutil.ReadAll(os.Stdin)
}

func ShaPrinter(algorithm int, data []byte) string {
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
	flag.IntVar(&algorithm, "algorithm", 1, "SHA algorithm, valid args are 1 and 256")
	flag.BoolVar(&help, "help", false, "Show this help and exit")
	flag.BoolVar(&version, "version", false, "Print Version")
	flag.Parse()

	if help {
		helpPrinter()
	}

	if version {
		versionPrinter()
	}
	if len(flag.Args()) == 1 {
		cliArgs = flag.Args()[0]
	}
	input, err := GetInput(cliArgs)
	if err != nil {
		fmt.Println("Error getting input.")
		os.Exit(-1)
	}
	fmt.Printf("%s ", ShaPrinter(algorithm, input))
	if cliArgs == "" {
		fmt.Printf(" -\n")
	} else {
		fmt.Printf(" %s\n", cliArgs)
	}
	os.Exit(0)
}
