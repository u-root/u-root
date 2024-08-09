// Copyright 2017-2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build plan9

package main

import (
	"fmt"
	"log"
	"os"
)

func main() {
	if len(os.Args) != 1 {
		log.Fatalf("Usage: id")
	}
	id, err := os.ReadFile("/env/user")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(id))
}
