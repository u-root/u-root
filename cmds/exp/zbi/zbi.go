// Copyright 2022 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// zbi dumps the header of a Zircon boot image.
//
// Synopsis:
//
//	zbi FILE
//
// Description:
//
//	Debugging purposes.
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/u-root/u-root/pkg/boot/zbi"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatalf("usage: %s FILE", os.Args[0])
	}

	image, err := zbi.Load(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}

	imageJSON, err := json.MarshalIndent(image, "", "  ")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s\n", string(imageJSON))
}
