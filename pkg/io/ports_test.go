// Copyright 2012-2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build linux,amd64 linux,386

package io

import (
	"log"
)

func ExampleIn() {
	var data uint8
	if err := In(0x3f8, &data); err != nil {
		log.Fatal(err)
	}
	log.Printf("%#02x\n", data)
}

func ExampleOut() {
	if err := Out(0x3f8, uint8('A')); err != nil {
		log.Fatal(err)
	}
}
