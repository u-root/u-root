// Copyright 2021 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// dumpebda reads and prints the Extended BIOS Data Area.
package main

import (
	"encoding/hex"
	"fmt"
	"log"
	"os"

	"github.com/u-root/u-root/pkg/boot/ebda"
)

func main() {
	f, err := os.OpenFile("/dev/mem", os.O_RDWR, 0)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	e, err := ebda.ReadEBDA(f)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("EBDA starts at %#X, length %#X bytes", e.BaseOffset, e.Length)
	fmt.Println(hex.Dump(e.Data))
}
