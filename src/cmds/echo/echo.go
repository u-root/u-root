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
	"strings"
)

var nonewline = flag.Bool("n", false, "suppress newline")

func echo(s string) error {

	_, err :=fmt.Printf("%s", s)

	if !*nonewline {
		fmt.Print("\n")
	}

	return err
	
}

func main() {
	
	flag.Parse()
	
	echo(strings.Join(flag.Args(), " "))
	
}
