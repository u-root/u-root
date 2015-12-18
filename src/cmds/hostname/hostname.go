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
	"os"
)

func hostname() (error) {

	hostname, error := os.Hostname()

	fmt.Println(hostname)
	return error
}

func main() {

	hostname()
	
}
