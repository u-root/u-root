// Copyright 2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// zimage dumps the header of a zImage.
//
// Synopsis:
//
//	zimage FILE
//
// Description:
//
//	This is mainly for debugging purposes.
package main

import (
	"fmt"
	"log"
	"os"

	"github.com/u-root/u-root/pkg/boot/zimage"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatalf("usage: %s FILE", os.Args[0])
	}

	f, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}

	header, err := zimage.Parse(f)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%#v\n", header)
}
