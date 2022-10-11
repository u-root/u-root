// Copyright 2012-2020 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Unmount unmounts new from old, or everything mounted on old if new is omitted.
//
// Synopsis:
//
//	unmount [ new ] old
package main

import (
	"log"
	"os"

	"github.com/u-root/u-root/pkg/namespace"
)

func main() {
	mod, err := namespace.ParseArgs(os.Args)
	if err != nil {
		log.Fatal(err)
	}
	if err := mod.Modify(namespace.DefaultNamespace, &namespace.Builder{}); err != nil {
		log.Fatal(err)
	}
}
