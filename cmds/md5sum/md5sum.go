// Copyright 2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"crypto/md5"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
)

func getInput(fileName string) (input []byte, err error) {

	if fileName != "" {
		return ioutil.ReadFile(fileName)
	}
	return ioutil.ReadAll(os.Stdin)
}

func helpPrinter() {

	fmt.Printf("Usage:\nmd5sum <File Name>\n")
	flag.PrintDefaults()
	os.Exit(0)
}

func versionPrinter() {
	fmt.Println("md5sum utility, URoot Version.")
	os.Exit(0)
}

func calculateMd5Sum(data []byte) string {
	return fmt.Sprintf("%x", md5.Sum(data))
}

func main() {
	var (
		help    bool
		version bool
	)
	cliArgs := ""
	flag.BoolVar(&help, "help", false, "Show this help and exit")
	flag.BoolVar(&version, "version", false, "Print Version")
	flag.Parse()

	if help {
		helpPrinter()
	}

	if version {
		versionPrinter()
	}

	if len(os.Args) >= 2 {
		cliArgs = os.Args[1]
	}
	input, err := getInput(cliArgs)
	if err != nil {
		fmt.Println("Error getting input.")
		os.Exit(-1)
	}
	fmt.Printf("%s ", calculateMd5Sum(input))
	if cliArgs == "" {
		fmt.Printf(" -\n")
	} else {
		fmt.Printf(" %s\n", cliArgs)
	}
	os.Exit(0)
}
