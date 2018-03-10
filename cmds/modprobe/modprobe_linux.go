// Copyright 2013-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// modprobe - Add and remove modules from the Linux Kernel
//
// Synopsis:
//     modprobe [-n] modulename [parameters...]
//
// Author:
//     Roland Kammerer <dev.rck@gmail.com>
package main

import (
	"flag"
	"log"
	"os"
	"strings"

	"github.com/u-root/u-root/pkg/kmodule"
)

const cmd = "modprobe [-n] modulename [parameters...]"

var (
	dryRun = flag.Bool("n", false, "Try run")
)

func init() {
	defUsage := flag.Usage
	flag.Usage = func() {
		os.Args[0] = cmd
		defUsage()
	}
}

func main() {
	flag.Parse()

	if flag.NArg() == 0 {
		log.Println("Usage: ERROR: one module and optional module options.")
		flag.Usage()
		os.Exit(1)
	}

	modName := flag.Args()[0]
	modOptions := strings.Join(flag.Args()[1:], " ")

	if err := kmodule.ProbeOptions(modName, modOptions, kmodule.ProbeOpts{DryRun: *dryRun}); err != nil {
		log.Fatalf("modprobe: Could not load module %q: %v", modName, err)
	}
}
