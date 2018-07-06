// Copyright 2012-2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Basename return name with leading path information removed.
//
// Synopsis:
//     basename NAME [SUFFIX]
//
// Options:
//     -s: optional flag for removing suffix
package main

import (
	"flag"
	"fmt"
	"os"
	"path"
	"strings"
)

var (
	suffix = flag.String("s", "", "Strip prefix from file")
)

func init() {
	// Usage Definition
	defUsage := flag.Usage
	flag.Usage = func() {
		os.Args[0] = "basename [-s] NAME [SUFFIX]"
		defUsage()
	}
}

func main() {
	flag.Parse()
	if len(flag.Args()) < 1 {
		flag.Usage()
		os.Exit(1)
	}

	if len(*suffix) > 0 {
		for i := 0; i < len(flag.Args()); i++ {
			_, fileName := path.Split(flag.Arg(i))
			fileName = strings.TrimSuffix(fileName, *suffix)
			fmt.Println(fileName)
		}
	} else if len(flag.Args()) > 1 {
		_, fileName := path.Split(flag.Arg(0))
		fileName = strings.TrimSuffix(fileName, flag.Arg(1))
		fmt.Println(fileName)
	} else {
		_, fileName := path.Split(flag.Arg(0))
		fmt.Println(fileName)
	}
}
