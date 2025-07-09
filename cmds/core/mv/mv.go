// Copyright 2012-2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// mv renames files and directories.
//
// Synopsis:
//
//	mv SOURCE [-u] TARGET
//	mv SOURCE... [-u] DIRECTORY
//
// Author:
//
//	Beletti (rhiguita@gmail.com)
package main

import (
	"log"
	"os"

	"github.com/u-root/u-root/pkg/core/mv"
)

func init() {
	log.SetFlags(0)
}

func main() {
	cmd := mv.New()
	err := cmd.Run(os.Args[1:]...)
	if err != nil {
		log.Fatal("mv: ", err)
	}
}
