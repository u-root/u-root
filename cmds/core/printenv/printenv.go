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

func printenv(w io.Writer, args ...string) {
	if len(args) == 0 {
		e := os.Environ()

		for _, v := range e {
			fmt.Fprintf(w, "%v\n", v)
		}
		return
	}

	for _, arg := range args {
		v, ok := os.LookupEnv(arg)
		if ok {
			fmt.Fprintln(w, v)
		}
	}
}

func main() {
	printenv(os.Stdout, os.Args[1:]...)
}
