// Copyright 2021 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build plan9

// ls prints the contents of a directory.
//
// Synopsis:
//     ls [OPTIONS] [DIRS]...
//
// Options:
//     -l: long form
//     -Q: quoted
//     -R: equivalent to findutil's find
//     -F: append indicator (one of */=>@|) to entries
//
// Bugs:
//     With the `-R` flag, directories are only ever printed once.
package main

import (
	flag "github.com/spf13/pflag"
)

var (
	final = flag.BoolP("print-last", "p", false, "Print only the final path element of each file name")
)
