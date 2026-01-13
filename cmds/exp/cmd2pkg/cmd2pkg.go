// Copyright 2026 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// cmd2pkg converts one or more commands into a set of packages.
package main

import (
	"flag"
	"log"
	"os"

	"github.com/u-root/gobusybox/src/pkg/golang"
	"github.com/u-root/u-root/pkg/ulog"
)

func main() {
	var (
		debug = flag.Bool("d", false, "Turn on forth package debugging using log.Printf")
		dir   = flag.String("dir", "", "Temp directory, if not set, os.MkdirTemp is used")
		l     = ulog.Log
		err   error
		v     = func(string, ...any) {}
	)
	flag.Parse()
	if *debug {
		v = log.Printf
	}
	v("Starting ...")

	if len(*dir) == 0 {
		*dir, err = os.MkdirTemp("", "cmd2pkg-")
		if err != nil {
			log.Fatal(err)
		}
	}

	opts := &Opts{
		Env:          golang.Default(),
		GenSrcDir:    *dir,
		CommandPaths: flag.Args(),
	}

	if err := BuildBusybox(l, opts); err != nil {
		log.Print(err)
	}

	log.Printf("Check your output (if any ...) in %s", *dir)
}
