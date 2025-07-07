// Copyright 2012-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// cat concatenates files and prints them to stdout.
//
// Synopsis:
//
//	cat [-u] [FILES]...
//
// Description:
//
//	If no files are specified, read from stdin.
//
// Options:
//
//	-u: ignored flag
package main

import (
	"context"
	"fmt"
	"os"

	"github.com/u-root/u-root/pkg/core/cat"
)

func main() {
	cmd := cat.New()
	exitCode, err := cmd.Run(context.Background(), os.Args[1:]...)
	if err != nil {
		fmt.Fprintf(os.Stderr, "cat: %v\n", err)
	}
	os.Exit(exitCode)
}
