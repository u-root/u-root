// Copyright 2013 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
Rm removes the named files.

The options are:
*/

package main

import (
	"fmt"
	"os"
)

func main() {
	for _,v := range(os.Args) {
		err := os.Remove(v)
		if err != nil {
			fmt.Printf("%v: %v\n", v, err)
		}
	}
}
