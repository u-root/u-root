// Copyright 2012 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
 print name of current/working directory
 created by Beletti (rhiguita@gmail.com)
*/

package main

import (
	"fmt"
	"log"
	"os"
)

func main() {
	if path, error := os.Getwd(); error != nil {
		log.Fatalf("%v", error)
	} else {
		fmt.Println(path)
	}
}
