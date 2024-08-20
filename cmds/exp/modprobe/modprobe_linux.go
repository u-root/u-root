// Copyright 2013-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// modprobe - Add and remove modules from the Linux Kernel
//
// Synopsis:
//
//	modprobe [-n] modulename [parameters...]
//	modprobe [-n] -a modulename...
//
// Author:
//
//	Roland Kammerer <dev.rck@gmail.com>
package main

import (
	"flag"
	"log"
	"os"
	"strings"

	"github.com/u-root/u-root/pkg/kmodule"
)

const cmd = "modprobe [-an] modulename[s] [parameters...]"

var (
	dryRun     = flag.Bool("n", false, "Dry run")
	all        = flag.Bool("a", false, "Insert all module names on the command line.")
	verboseAll = flag.Bool("va", false, "Insert all module names on the command line.")
	rootDir    = flag.String("d", "/", "Root directory for modules")
	kernelVer  = flag.String("S", "", "Set kernel version instead of using uname")
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

	opts := kmodule.ProbeOpts{
		RootDir: *rootDir,
		KVer:    *kernelVer,
	}
	if *dryRun {
		log.Println("Unique dependencies in load order, already loaded ones get skipped:")
		opts.DryRunCB = func(modPath string) {
			log.Println(modPath)
		}
	}

	// -va is just an alias for -a
	*all = *all || *verboseAll
	if *all {
		modNames := flag.Args()
		for _, modName := range modNames {
			if err := kmodule.ProbeOptions(modName, "", opts); err != nil {
				log.Printf("modprobe: Could not load module %q: %v", modName, err)
			}
		}
		os.Exit(0)
	}

	modName := flag.Args()[0]
	modOptions := strings.Join(flag.Args()[1:], " ")

	if err := kmodule.ProbeOptions(modName, modOptions, opts); err != nil {
		log.Fatalf("modprobe: Could not load module %q: %v", modName, err)
	}
}
