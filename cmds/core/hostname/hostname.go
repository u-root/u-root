// Copyright 2012-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// hostname prints or changes the system's hostname.
//
// Synopsis:
//     hostname [HOSTNAME]
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
	a := os.Args
	switch len(a) {
	case 2:
		if err := Sethostname(a[1]); err != nil {
			log.Fatalf("could not set hostname: %v", err)
		}
	case 1:
		hostname, err := os.Hostname()
		if err != nil {
			log.Fatalf("could not obtain hostname: %v", err)
		}
		fmt.Println(hostname)
	default:
		log.Fatalf("usage: hostname [HOSTNAME]")
	}
}
