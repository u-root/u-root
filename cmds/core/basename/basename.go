// Copyright 2012-2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Basename return name with leading path information removed.
//
// Synopsis:
//
//	basename NAME [SUFFIX]
package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

const (
	usageString = "Usage: basename NAME [SUFFIX]"
)

func usage(w io.Writer) {
	fmt.Fprintf(w, "%s", usageString)
}

func run(w io.Writer, args []string) {
	switch len(args) {
	case 2:
		fileName := filepath.Base(args[0])
		if fileName != args[1] {
			fileName = strings.TrimSuffix(fileName, args[1])
		}
		fmt.Fprintf(w, "%s\n", fileName)
	case 1:
		fileName := filepath.Base(args[0])
		fmt.Fprintf(w, "%s\n", fileName)
	default:
		usage(w)
	}
}

func main() {
	run(os.Stdout, os.Args[1:])
}
