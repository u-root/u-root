// Copyright 2017-2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"log"
	"os"

	"github.com/u-root/u-root/pkg/boot/fit"
)

var (
	dryRun    = flag.Bool("dryrun", true, "Do not actually kexec into the boot config")
	debug     = flag.Bool("d", true, "Print debug output")
	rootfs    = flag.String("r", "", "Root file system name")
	cmdline   = flag.String("c", "earlyprintk=ttyS0,115200,keep console=ttyS0", "command line")
	initramfs = flag.String("i", "", "InitRAMFS node name -- default none")
)

var v = func(string, ...interface{}) {}

func main() {
	flag.Parse()
	if *debug {
		v = log.Printf
	}
	if len(flag.Args()) < 2 {
		log.Fatal("Usage: fitboot <file> <node name>")
	}
	f, err := fit.New(flag.Args()[0])
	if err != nil {
		log.Fatal(err)
	}
	v("Loaded uimage: %s", f)
	f.Cmdline, f.Kernel, f.InitRAMFS, f.RootFS, f.Dryrun = *cmdline, flag.Args()[1], *initramfs, *rootfs, *dryRun
	if err := f.Load(*debug); err != nil {
		log.Fatal(err)
	}
	if *dryRun {
		v("Not trying to boot since this is a dry run")
		os.Exit(0)
	}
}
