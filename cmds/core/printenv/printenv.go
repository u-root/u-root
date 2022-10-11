// Copyright 2014-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Print environment variables.
//
// Synopsis:
//
//	printenv
package main

import (
	"fmt"
	"io"
	"os"
)

func printenv(w io.Writer) {
	e := os.Environ()

	for _, v := range e {
		fmt.Fprintf(w, "%v\n", v)
	}
}

func main() {
	printenv(os.Stdout)
}
