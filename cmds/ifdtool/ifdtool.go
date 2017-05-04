// Copyright 2012-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// ifdtool is the coreboot ifdtool in Go. It reads an ifd from stdin
// and writes a JSON version of it to stdout. Or something.
package main

import (
	"encoding/binary"
	"flag"
	"fmt"
	"log"
	"os"
)

var ()

func dump(c *chip) {
	fmt.Printf("%s", c.Data.String())
}

func readImage(c *chip) error {
	b := binary.LittleEndian
	var m uint32

	for m != magic {
		if err := binary.Read(c, b, &m); err != nil {
			return fmt.Errorf("Reading image: no magic: %v", err)
		}
	}
	log.Printf("Got the magic number")

	if err := binary.Read(c, b, &c.Data); err != nil {
		return fmt.Errorf("%v", err)
	}
	return nil

}

func usage() {
	log.Fatalf("Usage: yeah right")
}

func main() {
	var err error
	flag.Parse()

	if len(flag.Args()) < 1 {
		usage()
	}
	switch flag.Args()[0] {
	default:
		usage()
	case "dump":
	}

	c := &chip{
		Reader: os.Stdin,
	}

	if err = readImage(c); err != nil {
		log.Fatalf("Can't read image: %v", err)
	}

	switch flag.Args()[0] {
	case "dump":
		dump(c)
	}

}
