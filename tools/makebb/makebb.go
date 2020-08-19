// Copyright 2015-2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// makebb compiles many Go commands into one bb-style binary.
package main

import (
	"flag"
	"log"
	"os"
	"path/filepath"

	"github.com/u-root/u-root/pkg/bb"
	"github.com/u-root/u-root/pkg/golang"
	"github.com/u-root/u-root/pkg/uroot"
)

var outputPath = flag.String("o", "bb", "Path to busybox binary")

func main() {
	flag.Parse()

	// Why doesn't the log package export this as a default?
	l := log.New(os.Stdout, "", log.LstdFlags)
	env := golang.Default()
	if env.CgoEnabled {
		l.Printf("Disabling CGO for u-root...")
		env.CgoEnabled = false
	}
	l.Printf("Build environment: %s", env)

	pkgs := flag.Args()
	if len(pkgs) == 0 {
		pkgs = []string{"github.com/u-root/u-root/cmds/*/*"}
	}
	pkgs, err := uroot.ResolvePackagePaths(l, env, pkgs)
	if err != nil {
		l.Fatal(err)
	}

	o, err := filepath.Abs(*outputPath)
	if err != nil {
		l.Fatal(err)
	}

	if err := bb.BuildBusybox(env, pkgs, false /* noStrip */, o); err != nil {
		l.Fatal(err)
	}
}
