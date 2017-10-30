// Copyright 2016-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Change modifier bits of a file.
//
// Synopsis:
//     chmod MODE FILE...
//
// Desription:
//     MODE is a three character octal value.
package main

import (
	"flag"
	"log"
	"os"
	"strconv"
)

func main() {
	flag.Parse()
	if len(flag.Args()) < 2 {
		flag.PrintDefaults()
		log.Fatalf("usage: chmod mode filepath")
	}

	mode := flag.Args()[0]

	octval, err := strconv.ParseUint(mode, 8, 32)
	if err != nil {
		log.Fatalf("Unable to decode mode %q. Please use an octal value: %v", mode, err)
	} else if octval > 0777 {
		log.Fatalf("Invalid octal value %0o. Value should be less than or equal to 0777.", octval)
	}

	var exitError bool
	for _, name := range flag.Args()[1:] {
		if err := os.Chmod(name, os.FileMode(octval)); err != nil {
			log.Printf("%v", err)
			exitError = true
		}
	}
	if exitError {
		os.Exit(1)
	}
}
