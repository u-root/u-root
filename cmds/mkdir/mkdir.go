// Copyright 2012-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Mkdir makes a new directory.
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
	"syscall"
)

type modeFlag struct {
	set   bool
	value string
}

func (m *modeFlag) String() string {
	return m.value
}

func (m *modeFlag) Set(s string) error {
	m.value = s
	m.set = true
	return nil
}

const (
	cmd       = "mkdir [-m mode] [-v] [-p] <directory> [more directories]"
	StickyBit = 01000
	SgidBit   = 02000
	SuidBit   = 04000
)

var (
	mode    = &modeFlag{set: false, value: "0777"}
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

	// Complete Setting the flags
	flag.Var(mode, "m", "Directory mode")
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

	// Get Correct Access Mode
	if mode.set {
		syscall.Umask(0)
	}
	accMode64bit, err := strconv.ParseUint(mode.value, 8, 32)
	if err != nil || accMode64bit > 07777 {
		log.Fatalf("invalid mode '%s'", mode.value)
	}
	accMode := os.FileMode(accMode64bit)
	if accMode64bit&StickyBit != 0 {
		accMode |= os.ModeSticky
	}
	if accMode64bit&SgidBit != 0 {
		accMode |= os.ModeSetgid
	}
	if accMode64bit&SuidBit != 0 {
		accMode |= os.ModeSetuid
	}

	for _, name := range flag.Args() {
		if err := f(name, accMode); err != nil {
			log.Fatalf("%v: %v\n", name, err)
		} else {
			if *verbose {
				fmt.Printf("%v\n", name)
			}
			// os.Mkdir does not set up SGID and SUID correctly
			if accMode64bit&SgidBit != 0 || accMode64bit&SuidBit != 0 {
				os.Chmod(name, accMode)
			}
		}
	}
}
