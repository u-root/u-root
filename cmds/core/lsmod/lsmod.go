// Copyright 2012-2022 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build linux || darwin

// lsmod list currently loaded Linux kernel modules.
//
// Synopsis:
//
//	lsmod
//
// Description:
//
//	lsmod is a clone of lsmod(8)
//
// Author:
//
//	Roland Kammerer <dev.rck@gmail.com>
package main

import (
	"log"
	"os"

	"github.com/u-root/u-root/pkg/kmodule"
)

func main() {
	l, err := kmodule.New()
	if err != nil {
		log.Fatalln(err)
	}
	if err := kmodule.List(l, os.Stdout); err != nil {
		log.Fatalln(err)
	}
}
