// Copyright 2013 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"log"
	"os"
)

func main() {

	a, s, e := os.Args, 0, 0

	if len(a) != 3 {
		log.Fatal("Usage: seq <start> <end>")
	}

	if _, err := fmt.Sscanf(a[1] + " " + a[2], "%v %v", &s, &e); err != nil {
		log.Fatal("Usage: seq <start> <end>")
	} else {
		for s <= e {
			fmt.Printf("%d\n", s)
			s++
		}
	}
}
