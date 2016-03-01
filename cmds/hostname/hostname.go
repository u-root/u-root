// Copyright 2012 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//show the system's hostname
//created by Beletti (rhiguita@gmail.com)
package main

import (
	"fmt"
	"io"
	"os"
)

func hostname(w io.Writer) error {

	hostname, error := os.Hostname()

	fmt.Fprintf(w, "%v", hostname)
	return error
}

func main() {

	hostname(os.Stdout)
	fmt.Println()

}
