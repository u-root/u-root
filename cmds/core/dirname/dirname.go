// Copyright 2016-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// dirname prints out the directory name of one or more args.
// If no arg is given it returns an error and prints a message which,
// per the man page, is incorrect, but per the standard, is correct.
package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatalf("dirname: missing operand")
	}

	for _, n := range os.Args[1:] {
		fmt.Println(filepath.Dir(n))
	}
}
