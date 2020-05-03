// Copyright 2012-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// mkdir makes a new directory.
//
// Synopsis:
//     mkdir [-m mode] [-v] [-p] DIRECTORY...
//
// Options:
//     -m: make all needed directories in the path
//     -v: directory mode (ex: 666)
//     -p: print each directory as it is made
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
)

const (
	cmd                 = "mkdir [-m mode] [-v] [-p] <directory> [more directories]"
	DefaultCreationMode = 0777
	StickyBit           = 01000
	SgidBit             = 02000
	SuidBit             = 04000
)

var (
	mode    = flag.String("m", "", "Directory mode")
	mkall   = flag.Bool("p", false, "Make all needed directories in the path")
	verbose = flag.Bool("v", false, "Print each directory as it is made")
	f       = os.Mkdir
)

func init() {
	// Usage Definition
	defUsage := flag.Usage
	flag.Usage = func() {
		os.Args[0] = cmd
		defUsage()
	}
}

func main() {
	flag.Parse()
	if len(flag.Args()) < 1 {
		flag.Usage()
		os.Exit(1)
	}
	if *mkall {
		f = os.MkdirAll
	}

	// Get Correct Creation Mode
	var m uint64
	var err error
	if *mode == "" {
		m = DefaultCreationMode
	} else {
		m, err = strconv.ParseUint(*mode, 8, 32)
		if err != nil || m > 07777 {
			log.Fatalf("invalid mode '%s'", *mode)
		}
	}
	createMode := os.FileMode(m)
	if m&StickyBit != 0 {
		createMode |= os.ModeSticky
	}
	if m&SgidBit != 0 {
		createMode |= os.ModeSetgid
	}
	if m&SuidBit != 0 {
		createMode |= os.ModeSetuid
	}

	for _, name := range flag.Args() {
		if err := f(name, createMode); err != nil {
			log.Printf("%v: %v\n", name, err)
			continue
		}
		if *verbose {
			fmt.Printf("%v\n", name)
		}
		if *mode != "" {
			os.Chmod(name, createMode)
		}
	}
}
