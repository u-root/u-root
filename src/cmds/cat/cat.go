/* Copyright 2012 the u-root Authors. All rights reserved
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 *
 *
 * cat - concatenate and print files
 */

package main

import (
	"flag"
	"fmt"
	"io"
	"os"
)

var u_flag = flag.Bool("u", false, "Write bytes from the input file to the standard output without delay as each is read.")
var help_flag = flag.Bool("h", false,	"Display cat's help.")

func help() {
	fmt.Println("cat usage: 'cat [-u] [file ...]'")
	os.Exit(1)
}

func cat() {
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

func main() {
	flag.Parse()
	if len(os.Args) == 1 {
		io.Copy(os.Stdout, os.Stdin)

	} else {
		if *help_flag == true {
			help()
		} else if *u_flag == true {
			// treat -u flag!
			cat()
		} else {
			cat()
		}
	}
}
