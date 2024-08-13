// Copyright 2012-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build linux

package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/u-root/u-root/pkg/mount"
)

var (
	help    = flag.Bool("h", false, "Help")
	version = flag.Bool("V", false, "Version")
)

func usage() string {
	return "switch_root [-h] [-V]\nswitch_root newroot init"
}

func main() {
	flag.Parse()

	if len(flag.Args()) == 0 {
		fmt.Println(usage())
		os.Exit(0)
	}

	if *help {
		fmt.Println(usage())
		os.Exit(0)
	}

	if *version {
		fmt.Println("Version XX")
		os.Exit(0)
	}

	newRoot := flag.Args()[0]
	init := flag.Args()[1]

	if err := mount.SwitchRoot(newRoot, init); err != nil {
		log.Fatalf("switch_root failed %v\n", err)
	}
}
