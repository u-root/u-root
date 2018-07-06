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
	"flag"
	"fmt"
	"os"
	"path"
	"strings"
)

func init() {
	// Usage Definition
	defUsage := flag.Usage
	flag.Usage = func() {
		os.Args[0] = "basename NAME [SUFFIX]"
		defUsage()
	}
}

func main() {
	flag.Parse()
	if len(flag.Args()) < 1 {
		flag.Usage()
		os.Exit(1)
	}

	if len(flag.Args()) == 2 {
		_, fileName := path.Split(flag.Arg(0))
		if fileName != flag.Arg(1) {
			fileName = strings.TrimSuffix(fileName, flag.Arg(1))
		}
		fmt.Println(fileName)
	} else {
		_, fileName := path.Split(flag.Arg(0))
		fmt.Println(fileName)
	}
}
