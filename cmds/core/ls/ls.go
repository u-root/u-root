// Copyright 2013-2017 the u-root Authors. All rights reserved
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
	"context"
	"fmt"
	"os"

	"github.com/u-root/u-root/pkg/core/ls"
)

func main() {
	cmd := ls.New()
	err := cmd.Run(context.Background(), os.Args[1:]...)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ls: %v\n", err)
	}
}
