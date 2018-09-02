// Copyright 2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"crypto/md5"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/spf13/pflag"
)

func getInput(fileName string) (input []byte, err error) {

	if fileName != "" {
		return ioutil.ReadFile(fileName)
	}
	return ioutil.ReadAll(os.Stdin)
}

func helpPrinter() {

	fmt.Printf("Usage:\nmd5sum <File Name>\n")
	pflag.PrintDefaults()
	os.Exit(0)
}

func versionPrinter() {
	fmt.Println("md5sum utility, URoot Version.")
	os.Exit(0)
}

func calculateMd5Sum(data []byte) string {
	return fmt.Sprintf("%x", md5.Sum(data))
}

func checksum(hasher hash.Hash, r io.Reader) {
	if _, err := io.Copy(hasher, r); err != nil {
		log.Fatal(err)
	}
	sum := hasher.Sum(nil)
	fmt.Println(hex.EncodeToString(sum))
}

func main() {
	var (
		help    bool
		version bool
	)
	cliArgs := ""
	pflag.BoolVarP(&help, "help", "h", false, "Show this help and exit")
	pflag.BoolVarP(&version, "version", "v", false, "Print Version")
	pflag.Parse()

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
