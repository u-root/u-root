// Copyright 2012-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Print the system's hostname.
//
// Synopsis:
//     hostname
//
// Author:
//     Beletti <rhiguita@gmail.com>
package main

import (
	"fmt"
	"log"
	"os"
)

func main() {
	hostname, err := os.Hostname()
	if err != nil {
		log.Fatalf("could not obtain hostname: %v", err)
	}

	fmt.Println(hostname)
}
