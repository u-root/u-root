// Copyright 2012 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"io"
	"os"
)

var a = flag.Bool("a", false, "What?")

func main() {
	flag.Parse()
	p, err := os.Getenv("PATH")
	if err != nil {
		fmt.Printf("No path! %v\n", err)
		os.Exit(1)
	}

	p := strings.Split(p, ":")

		for _, name := range flag.Args() {
			
			for i := range(p) {
				f := path.Join(p[i], name)
				if s, err := os.Stat(f); err == nil {
					if 
				
				}
		}
	}
}
