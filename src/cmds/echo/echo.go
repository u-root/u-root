// Copyright 2013 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//Echo writes its arguments separated by blanks and terminated by a newline on the standard output.
package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
)

var nonewline = flag.Bool("n", false, "suppress newline")

func echo(s string, w io.Writer) error {

	_, err := fmt.Fprintf(w, "%s", s)

	if !*nonewline {
		fmt.Fprint(w, "\n")
	}

	return err

}

func main() {

	flag.Parse()

	echo(strings.Join(flag.Args(), " "), os.Stdout)

}
