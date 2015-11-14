// Copyright 2013 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
Rm removes the named files.

The options are:
*/

package main

import (
	"fmt"
	"os"
)

func printErr(file string, err error) {
	if err != nil {
		fmt.Printf("%v: %v\n", file, err)
	}
}

func main() {
	start := 1
	if len(os.Args) > 1 {
		if os.Args[1] == "-r" || os.Args[1] == "-R" {
			start = 2
		}
	}

	for _,v := range(os.Args[start:]) {
		if start == 2 && os.Args[1] == "-r" || os.Args[1] == "-R" {
			printErr(v, os.RemoveAll(v))
		} else {
			printErr(v, os.Remove(v))
		}
	}
}
