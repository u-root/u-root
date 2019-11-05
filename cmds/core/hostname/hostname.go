// Copyright 2012-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Print or set the system's hostname.
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

	"golang.org/x/sys/unix"
)

func main() {
	if len(os.Args) == 2 {
		newHostname := os.Args[1]

		err := unix.Sethostname([]byte(newHostname))
		if err != nil {
			log.Fatalf("could not set hostname: %v", err)
		}
	} else if len(os.Args) == 1 {
		hostname, err := os.Hostname()
		if err != nil {
			log.Fatalf("could not obtain hostname: %v", err)
		}

		fmt.Println(hostname)
	} else {
		log.Fatalf("usage: hostname [HOSTNAME]")
	}
}
