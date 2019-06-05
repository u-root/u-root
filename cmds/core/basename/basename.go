// Copyright 2012-2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Basename return name with leading path information removed.
//
// Synopsis:
//     basename NAME [SUFFIX]
//
package main

import (
	"fmt"
	"log"
	"path/filepath"
	"strings"

	flag "github.com/spf13/pflag"
)

func usage() {
	log.Fatal("Usage: basename NAME [SUFFIX]")
}

func main() {
	flag.Parse()

	args := flag.Args()
	switch len(args) {
	case 2:
		fileName := filepath.Base(args[0])
		if fileName != args[1] {
			fileName = strings.TrimSuffix(fileName, args[1])
		}
		fmt.Println(fileName)
	case 1:
		fileName := filepath.Base(args[0])
		fmt.Println(fileName)
	default:
		usage()
	}
}
