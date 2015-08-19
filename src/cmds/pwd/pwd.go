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
	"os"
)

func main() {	
	path, error := os.Getwd()

	if error == nil {
		fmt.Println(path)	
	} else {
		fmt.Printf("Error: %v",error)
	}
}
