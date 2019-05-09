// Copyright 2012-2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build linux,amd64 linux,386

package memio

import (
	"fmt"
	"log"
)

func ExampleIn() {
	var data Uint8
	if err := In(0x3f8, &data); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%v\n", data)
}

func ExampleOut() {
	data := Uint8('A')
	if err := Out(0x3f8, &data); err != nil {
		log.Fatal(err)
	}
}
