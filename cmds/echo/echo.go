// Copyright 2013-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Echo writes its arguments separated by blanks and terminated by a newline on
// the standard output.
//
// Synopsis:
//     echo [STRING]...
package main

import (
	"fmt"
	"io"
	"os"
	"strings"

	flag "github.com/juju/gnuflag"
)

var nonewline = flag.Bool("n", false, "suppress newline")

func echo(w io.Writer, s ...string) error {
	_, err := fmt.Fprintf(w, "%s", strings.Join(s, " "))

	if !*nonewline {
		fmt.Fprint(w, "\n")
	}

	return err
}

func main() {

	echo(os.Stdout, flag.Args()...)
}
