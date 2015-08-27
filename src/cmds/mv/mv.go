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
		err := os.Rename(a[0], a[1])
		if err != nil {
			if os.IsNotExist(err) {
				fmt.Printf("file does not exist\n")
			} else {
				if os.IsPermission(err) {
					fmt.Printf("you do not have permission\n")
				} else {
					fmt.Printf("unknow error\n")
				}
			}
		}
	}
}
