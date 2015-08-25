// Copyright 2012 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
 show the system's hostname
 created by Beletti (rhiguita@gmail.com)
*/

package main

import (
	"fmt"
	"log"
	"os"
)

func main() {
	if hostname, error := os.Hostname(); error != nil {
		log.Fatalf("%v", error)
	} else {
		fmt.Println(hostname)
	}
}
