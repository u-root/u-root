// Copyright 2012-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// mkdir makes a new directory.
//
// Synopsis:
//
//	mkdir [-m mode] [-v] [-p] DIRECTORY...
//
// Options:
//
//	-m: make all needed directories in the path
//	-v: directory mode (ex: 666)
//	-p: print each directory as it is made
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/u-root/u-root/pkg/uroot/util"
)

const (
	cmd                 = "mkdir [-m mode] [-v] [-p] <directory> [more directories]"
	defaultCreationMode = 0o777
	stickyBit           = 0o1000
	sgidBit             = 0o2000
	suidBit             = 0o4000
)

var (
	mode    = flag.String("m", "", "Directory mode")
	mkall   = flag.Bool("p", false, "Make all needed directories in the path")
	verbose = flag.Bool("v", false, "Print each directory as it is made")
)

func init() {
	flag.Usage = util.Usage(flag.Usage, cmd)
}

func mkdir(mode string, mkall, verbose bool, args []string) error {
	f := os.Mkdir
	if mkall {
		f = os.MkdirAll
	}

	// Get Correct Creation Mode
	var m uint64
	var err error
	if mode == "" {
		m = defaultCreationMode
	} else {
		m, err = strconv.ParseUint(mode, 8, 32)
		if err != nil || m > 0o7777 {
			return fmt.Errorf("invalid mode %q", mode)
		}
	}
	createMode := os.FileMode(m)
	if m&stickyBit != 0 {
		createMode |= os.ModeSticky
	}
	if m&sgidBit != 0 {
		createMode |= os.ModeSetgid
	}
	if m&suidBit != 0 {
		createMode |= os.ModeSetuid
	}

	for _, name := range args {
		if err := f(name, createMode); err != nil {
			log.Printf("%v: %v\n", name, err)
			continue
		}
		if verbose {
			fmt.Printf("%v\n", name)
		}
		if mode != "" {
			os.Chmod(name, createMode)
		}
	}
	return nil
}

func main() {
	flag.Parse()
	if len(flag.Args()) < 1 {
		flag.Usage()
		os.Exit(1)
	}
	if err := mkdir(*mode, *mkall, *verbose, flag.Args()); err != nil {
		log.Fatal(err)
	}
}
