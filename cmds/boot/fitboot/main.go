// Copyright 2017-2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//go:build !tinygo || tinygo.enable

package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/u-root/u-root/pkg/acpi"
	"github.com/u-root/u-root/pkg/boot"
	"github.com/u-root/u-root/pkg/boot/fit"
	"github.com/u-root/u-root/pkg/vfile"
)

var (
	dryRun     = flag.Bool("dryrun", false, "Do not actually kexec into the boot config")
	debug      = flag.Bool("d", false, "Print debug output")
	cmdline    = flag.String("c", "earlyprintk=ttyS0,115200,keep console=ttyS0", "command line")
	config     = flag.String("config", "", "FIT configuration to use")
	kernel     = flag.String("k", "", "Kernel image node name.")
	initramfs  = flag.String("i", "", "InitRAMFS node name -- default none")
	ringPath   = flag.String("r", "", "Path to PGP keyring. Enforces signature if non-empty path")
	rsdpLookup = flag.Bool("rsdp", false, "Derrive RSDP table pointer from environment")
)

var v = func(string, ...interface{}) {}

func main() {
	flag.Parse()

	if *debug {
		v = log.Printf
	}

	if len(flag.Args()) != 1 {
		log.Fatal("Usage: fitboot <file>")
	}
	f, err := fit.New(flag.Args()[0])
	if err != nil {
		log.Fatal(err)
	}

	f.Cmdline, f.Kernel, f.InitRAMFS, f.ConfigOverride = *cmdline, *kernel, *initramfs, *config

	kn, in, err := f.LoadConfig()
	if err == nil {
		f.Kernel, f.InitRAMFS = kn, in
	} else {
		v("Configuration is not available: %v", err)
	}

	if f.Kernel == "" {
		log.Fatal("kernel name is not found in fit configuration or pass through -k.")
	}

	v("Kernel name=%s, initramfs=%s", f.Kernel, f.InitRAMFS)

	kernelCmd := *cmdline
	if *rsdpLookup {
		r, err := acpi.GetRSDP()
		if err != nil {
			log.Fatal("Unable to find acpi table in the environment.")
		}
		v("Found an RSDP at %#x", r.RSDPAddr())
		kernelCmd = fmt.Sprintf("acpi_rsdp=%x %s", r.RSDPAddr(), kernelCmd)
	}

	f.Cmdline = kernelCmd

	if *ringPath != "" {
		ring, err := vfile.GetKeyRing(*ringPath)
		if err != nil {
			log.Fatal(err)
		}
		f.KeyRing = ring
	}

	if err := f.Load(boot.WithVerbose(*debug)); err != nil {
		log.Fatal(err)
	}

	if *dryRun {
		v("Not trying to boot since this is a dry run")
		os.Exit(0)
	}

	if err := boot.Execute(); err != nil {
		log.Fatal(err)
	}
}
