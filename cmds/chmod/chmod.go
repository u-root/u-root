// Copyright 2016-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Chmod changes modifier bits of a file.
//
// Synopsis:
//     chmod [-R] [--reference=file] [MODE] FILE...
//
// Description:
//     MODE is a three character octal value.
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type flags struct {
	recursive bool
	reference string
}

func (f *flags) registerFlags(flag *flag.FlagSet) {
	flag.BoolVar(&f.recursive,
		"R",
		false,
		"do changes recursively")

	flag.BoolVar(&f.recursive,
		"recursive",
		false,
		"do changes recursively")

	flag.StringVar(&f.reference,
		"reference",
		"",
		"use mode from reference file")
}

func main() {
	f := flags{}
	f.registerFlags(flag.CommandLine)
	flag.Parse()

	if flag.NArg() < 1 {
		fmt.Fprintf(os.Stderr, "Usage of %s: [mode] filepath\n", os.Args[0])
		flag.PrintDefaults()
		os.Exit(1)
	}

	if flag.NArg() < 2 && len(f.reference) == 0 {
		fmt.Fprintf(os.Stderr, "Usage of %s: [mode] filepath\n", os.Args[0])
		flag.PrintDefaults()
		os.Exit(1)
	}

	if err := chmodMain(f, flag.Args()); err != nil {
		log.Fatal(err)
	}
}

func chmodMain(f flags, args []string) error {
	if len(f.reference) > 0 {
		fi, err := os.Stat(f.reference)
		if err != nil {
			return fmt.Errorf("bad reference file %q: %v", f.reference, err)
		}
		return doChmod(f.recursive, fi.Mode(), args)
	}

	modeString := args[0]
	octval, err := strconv.ParseUint(modeString, 8, 32)
	if err != nil {
		return fmt.Errorf("unable to decode mode %q: must use an octal value: %v", modeString, err)
	} else if octval > 0777 {
		return fmt.Errorf("invalid octal value %0o: value should be less than or equal to 0777", octval)
	}
	return doChmod(f.recursive, os.FileMode(octval), args[1:])
}

type errors []error

func (e *errors) Add(f error) {
	if f != nil {
		*e = append(*e, f)
	}
}

func (e errors) AnyError() error {
	if e == nil {
		return nil
	}
	// Only return e if any of the errors are non-nil.
	for _, f := range e {
		if f != nil {
			return e
		}
	}
	return nil
}

func (e errors) Error() string {
	s := make([]string, 0, len(e))
	for _, f := range e {
		if f != nil {
			s = append(s, f.Error())
		}
	}

	// Only one error? Just return that one.
	if len(s) == 1 {
		return s[0]
	}
	return fmt.Sprintf("multiple errors: %s", strings.Join(s, "; "))
}

func doChmod(recursive bool, mode os.FileMode, fileList []string) error {
	var e errors
	for _, name := range fileList {
		if recursive {
			e.Add(filepath.Walk(name, func(path string, info os.FileInfo, _ error) error {
				return os.Chmod(path, mode)
			}))
		} else {
			e.Add(os.Chmod(name, mode))
		}
	}

	return e.AnyError()
}
