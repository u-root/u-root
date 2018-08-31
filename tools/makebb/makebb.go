// Copyright 2015-2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// makebb compiles many Go commands into one bb-style binary.
package main

import (
	"flag"
	"log"
	"path/filepath"

	"github.com/u-root/u-root/pkg/golang"
	"github.com/u-root/u-root/pkg/uroot"
	"github.com/u-root/u-root/pkg/uroot/builder"
)

var outputPath = flag.String("o", "bb", "Path to busybox binary")

func main() {
	flag.Parse()

	env := golang.Default()
	if env.CgoEnabled {
		log.Printf("Disabling CGO for u-root...")
		env.CgoEnabled = false
	}
	log.Printf("Build environment: %s", env)

	pkgs := flag.Args()
	if len(pkgs) == 0 {
		pkgs = []string{"github.com/u-root/u-root/cmds/*"}
	}
	pkgs, err := uroot.ResolvePackagePaths(env, pkgs)
	if err != nil {
		log.Fatal(err)
	}

	o, err := filepath.Abs(*outputPath)
	if err != nil {
		log.Fatal(err)
	}

	if err := builder.BuildBusybox(env, pkgs, o); err != nil {
		log.Fatal(err)
	}
}
