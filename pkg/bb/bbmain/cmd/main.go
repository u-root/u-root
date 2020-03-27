// Copyright 2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package main is the busybox main.go template.
package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/u-root/u-root/pkg/bb/bbmain"
	"github.com/u-root/u-root/pkg/upath"
)

func run() {
	name := filepath.Base(os.Args[0])
	if err := bbmain.Run(name); err != nil {
		log.Fatalf("%s: %v", name, err)
	}
}

func main() {
	os.Args[0] = upath.ResolveUntilLastSymlink(os.Args[0])

	run()
}

func init() {
	m := func() {
		if len(os.Args) == 1 {
			log.Fatalf("Invalid busybox command: %q", os.Args)
		}
		// Use argv[1] as the name.
		os.Args = os.Args[1:]
		run()
	}
	bbmain.Register("bb", bbmain.Noop, bbmain.ListCmds)
	bbmain.RegisterDefault(bbmain.Noop, m)
}
