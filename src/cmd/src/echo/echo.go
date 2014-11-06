// Copyright 2013 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
Echo writes its arguments separated by blanks and terminated by a newline on the standard output.

The options are:
	–n		suppress newline.
*/

package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

var nonewline = flag.Bool("n", false, "suppress newline")

func main() {
	flag.Parse()

	_, err := fmt.Printf("%s", strings.Join(flag.Args(), " "))
	if err != nil {
		os.Exit(1) // "write error" on Plan 9
	}

	if !*nonewline {
		fmt.Print("\n")
	}
}
