// Copyright 2013 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
Ls reads the directories in the command line and prints out the names.

The options are:
	–l		Long form.
*/

package main

import (
	"flag"
	"fmt"
	"io/ioutil"
)

var long = flag.Bool("l", false, "Long form")

func dir(path string) error {
	ents, err := ioutil.ReadDir(path)
	if err != nil {
		return err
	}
	for _,v := range(ents) {
		fmt.Printf("%v\n", v)
	}
	return nil
}
	
func main() {
	flag.Parse()

	dirs := flag.Args()

	if len(dirs) == 0 {
		dirs = []string{"."}
	}
	for _,v := range(dirs) {
		err := dir(v)
		if err != nil {
			fmt.Printf("%v: %v\n", v, err)
		}
	}
}
