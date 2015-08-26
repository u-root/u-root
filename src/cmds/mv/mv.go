// Copyright 2012 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
 move (rename) files
 created by Beletti (rhiguita@gmail.com)
*/

package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	flag.Parse()
	a := flag.Args()
	if len(a) < 2 {
		fmt.Printf("mv - missing file operand\n")
	} else {
		error := os.Rename(a[0], a[1])
		if error != nil {
			fmt.Printf("%v\n", error)
		}
	}
}
