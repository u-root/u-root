// Copyright 2012 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
Cat reads each file from its arguments in sequence and writes it on the standard output.
*/

package main

import (
	"fmt"
	"io"
	"os"
)

func main() {
	if len(os.Args) == 1 {
		io.Copy(os.Stdout, os.Stdin)
	} else {
		for _, name := range os.Args[1:] {
			f, err := os.Open(name)
			if err != nil {
				fmt.Fprintf(os.Stderr, "can't open %s: %v\n", name, err)
				os.Exit(1)
			}

			_, err = io.Copy(os.Stdout, f)
			if err != nil {
				fmt.Fprintf(os.Stderr, "error %s: %v", name, err)
				os.Exit(1)
			}

			f.Close()
		}
	}
}
