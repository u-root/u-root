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

type command struct {
	dir string
	l   ulog.Logger
}

func (c *command) execute(args ...string) error {
	var err error

	if len(c.dir) == 0 {
		c.dir, err = os.MkdirTemp("", "cmd2pkg-")
		return err
	}

	opts := &Opts{
		Env:          golang.Default(),
		GenSrcDir:    c.dir,
		CommandPaths: args,
	}

	if err := BuildBusybox(c.l, opts); err != nil {
		return err
	}

	return nil
}

func main() {
	var dir = flag.String("dir", "", "Temp directory, if not set, os.MkdirTemp is used")
	flag.Parse()
	if err := (&command{l: ulog.Log, dir: *dir}).execute(flag.Args()...); err != nil {
		log.Fatal(err)
	}
}
