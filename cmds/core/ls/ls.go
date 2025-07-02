// Copyright 2013-2025 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// ls prints the contents of a directory.
//
// Synopsis:
//
//	ls [OPTIONS] [DIRS]...
//
// Options:
//
//	-a: show hidden files
//	-h: show human-readable sizes
//	-d: show directories but not their contents
//	-F: append indicator (, one of */=>@|) to entries
//	-l: long form
//	-Q: quoted
//	-R: equivalent to findutil's find
//	-s: sort by size
//
// Bugs:
//
//	With the `-R` flag, directories are only ever printed once.
package main

import (
	"os"

	"github.com/u-root/u-root/pkg/ls"
)

func main() {
	wd, _ := os.Getwd()
	os.Exit(ls.Command(os.Stdout, os.Stderr, wd).Run(os.Args...))
}
