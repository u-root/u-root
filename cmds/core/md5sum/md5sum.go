// Copyright 2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// md5sum prints an md5 hash generated from file contents.
package main

import (
	"crypto/md5"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"

	"github.com/spf13/pflag"
)

func getInput() (input []byte, err error) {
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

func calculateMd5Sum(fileName string, data []byte) string {
	if len(data) > 0 {
		return fmt.Sprintf("%x", md5.Sum(data))
	}

	fileDesc, err := os.Open(fileName)
	if err != nil {
		log.Fatal(err)
	}
	defer fileDesc.Close()

	md5Generator := md5.New()
	if _, err := io.Copy(md5Generator, fileDesc); err != nil {
		log.Fatal(err)
	}

	md5Sum := fmt.Sprintf("%x", md5Generator.Sum(nil))
	return md5Sum
}

func main() {
	var (
		help    bool
		version bool
		input   []byte
		err     error
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
	if cliArgs == "" {
		input, err = getInput()
		if err != nil {
			fmt.Println("Error getting input.")
			os.Exit(-1)
		}
	}
	fmt.Printf("%s ", calculateMd5Sum(cliArgs, input))
	if cliArgs == "" {
		fmt.Printf(" -\n")
	} else {
		fmt.Printf(" %s\n", cliArgs)
	}
	os.Exit(0)
}
