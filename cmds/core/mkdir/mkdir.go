// Copyright 2012-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// mkdir makes one or more new directories.
//
// Synopsis:
//
//	mkdir [-m mode] [-v] [-p] DIRECTORY...
//
// Options:
//
//	-m: directory mode (ex: 755)
//	-v: print each directory as it is made
//	-p: make all needed directories in the path
package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"

	"github.com/u-root/u-root/pkg/uroot/unixflag"
)

const (
	defaultCreationMode = 0o777
	stickyBit           = 0o1000
	sgidBit             = 0o2000
	suidBit             = 0o4000
)

var errNoDirs = errors.New("no directories specified")
var errInvalidMode = errors.New("invalid mode")

type flags struct {
	mode    string
	mkall   bool
	verbose bool
}

func mkdir(stdout, stderr io.Writer, args []string) error {
	var f flags

	fs := flag.NewFlagSet("mkdir", flag.ContinueOnError)
	fs.SetOutput(stderr)

	fs.StringVar(&f.mode, "m", "", "Directory mode")
	fs.BoolVar(&f.mkall, "p", false, "Make all needed directories in the path")
	fs.BoolVar(&f.verbose, "v", false, "Print each directory as it is made")

	fs.Usage = func() {
		fmt.Fprintf(fs.Output(), "Usage: mkdir [-m mode] [-v] [-p] DIRECTORY...\n\n")
		fmt.Fprintf(fs.Output(), "mkdir makes one or more new directories.\n\n")
		fmt.Fprintf(fs.Output(), "Options:\n")
		fs.PrintDefaults()
	}

	if err := fs.Parse(unixflag.ArgsToGoArgs(args)); err != nil {
		return err
	}

	if len(fs.Args()) < 1 {
		fs.Usage()
		return errNoDirs
	}

	fm := os.Mkdir
	if f.mkall {
		fm = os.MkdirAll
	}

	// Get Correct Creation Mode
	var m uint64
	var err error
	if f.mode == "" {
		m = defaultCreationMode
	} else {
		m, err = strconv.ParseUint(f.mode, 8, 32)
		if err != nil || m > 0o7777 {
			return fmt.Errorf("%w: %q", errInvalidMode, f.mode)
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

	for _, name := range fs.Args() {
		if err := fm(name, createMode); err != nil {
			fmt.Fprintf(stderr, "%v: %v\n", name, err)
			continue
		}
		if f.verbose {
			fmt.Fprintf(stdout, "%v\n", name)
		}
		if f.mode != "" {
			os.Chmod(name, createMode)
		}
	}
	return nil
}

func main() {
	if err := mkdir(os.Stdout, os.Stderr, os.Args[1:]); err != nil {
		log.Fatal(err)
	}
}
