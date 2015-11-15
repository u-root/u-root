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
	"flag"
)

var recursive = flag.Bool("R", false, "Remove file hierarchies")
var recursive_too = flag.Bool("r", false, "Equivalent to -R.")

func main() {
	flag.Parse()

	var f = os.Remove
	if *recursive || *recursive_too {
		f = os.RemoveAll
	}

	for _,arg := range(os.Args[1:]) {
		if arg == "-r" || arg == "-R" {
			continue
		}

		err := f(arg)
		if err != nil {
			fmt.Printf("%v: %v\n", arg, err)
		}
	}
}
