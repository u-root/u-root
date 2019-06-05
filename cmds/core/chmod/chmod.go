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
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
)

var (
	recursive bool
	reference string
)

func init() {
	flag.BoolVar(&recursive,
		"R",
		false,
		"do changes recursively")

	flag.BoolVar(&recursive,
		"recursive",
		false,
		"do changes recursively")

	flag.StringVar(&reference,
		"reference",
		"",
		"use mode from reference file")
}

func main() {
	flag.Parse()
	if len(flag.Args()) < 1 {
		fmt.Fprintf(os.Stderr, "Usage of %s: [mode] filepath\n", os.Args[0])
		flag.PrintDefaults()
		os.Exit(1)
	}

	if len(flag.Args()) < 2 && reference == "" {
		fmt.Fprintf(os.Stderr, "Usage of %s: [mode] filepath\n", os.Args[0])
		flag.PrintDefaults()
		os.Exit(1)
	}

	var mode os.FileMode
	var fileList []string

	if reference != "" {
		fi, err := os.Stat(reference)
		if err != nil {
			log.Fatalf("bad reference file: %v", err)

		}
		mode = fi.Mode()
		fileList = flag.Args()
	} else {
		modeString := flag.Args()[0]
		octval, err := strconv.ParseUint(modeString, 8, 32)
		if err != nil {
			log.Fatalf("Unable to decode mode %q. Please use an octal value: %v", modeString, err)
		} else if octval > 0777 {
			log.Fatalf("Invalid octal value %0o. Value should be less than or equal to 0777.", octval)
		}
		mode = os.FileMode(octval)
		fileList = flag.Args()[1:]
	}

	var exitError bool
	for _, name := range fileList {
		if recursive {
			err := filepath.Walk(name, func(path string,
				info os.FileInfo,
				err error) error {
				return os.Chmod(path, mode)
			})
			if err != nil {
				log.Printf("%v", err)
				exitError = true
			}
		} else {
			if err := os.Chmod(name, mode); err != nil {
				log.Printf("%v", err)
				exitError = true
			}
		}
	}
	if exitError {
		os.Exit(1)
	}
}
