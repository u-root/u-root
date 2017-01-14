// Copyright 2016-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Change modifier bits of a file.
//
// Synopsis:
//     chmod MODE FILE...
//
// Desription:
//     MODE is a three character octal value.
package main

import (
	"flag"
	"fmt"
	"log"
	"strconv"
	"syscall"
)

func main() {
	flag.Parse()
	if len(flag.Args()) < 2 {
		flag.PrintDefaults()
		log.Fatalf("usage: chmod mode filepath")
	}

	octval, err := strconv.ParseUint(flag.Args()[0], 8, 32)
	if err != nil {
		log.Fatalf("Unable to decode mode. Please use an octal value. arg was %s, err was %v", flag.Args()[0], err)
	} else if octval > 0777 {
		log.Fatalf("Invalid octal value. Value larger than 777, was %o", octval)
	}

	mode := uint32(octval)

	var errors string
	for _, arg := range flag.Args()[1:] {
		if err := syscall.Chmod(arg, mode); err != nil {
			errors += fmt.Sprintf("Unable to chmod, filename was %s, err was %v\n", arg, err)
		}
	}
	if errors != "" {
		log.Fatalf(errors)
	}

}
