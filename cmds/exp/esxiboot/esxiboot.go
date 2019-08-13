// Copyright 2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// esxiboot executes ESXi kernel over the running kernel.
//
// Synopsis:
//     esxiboot [-d --device -p --partition] [-c --config] [-r --cdrom]
//
// Description:
//     Loads and executes ESXi kernel.
//
// Options:
//     --config=FILE or -c=FILE: set the ESXi config
//     --device=FILE or -d=FILE: set an ESXi disk to boot from
//     --cdrom=FILE or -r=FILE: set an ESXI CDROM to boot from
//     --partition=NUM or -p=NUM: which partition to boot ESXi from (either 5 or 6), only used with --device
//
// --device is required to kexec installed ESXi instance.
// You don't need it if you kexec ESXi installer.
//
// The config file has the following syntax:
//
// kernel=PATH
// kernelopt=OPTS
// modules=MOD1 [ARGS] --- MOD2 [ARGS] --- ...
//
// Lines starting with '#' are ignored.

package main

import (
	"io/ioutil"
	"log"
	"os"

	flag "github.com/spf13/pflag"

	"github.com/u-root/u-root/pkg/bootld/boot"
	"github.com/u-root/u-root/pkg/bootld/esxi"
)

var (
	cfg       = flag.StringP("config", "c", "", "ESXi config file")
	cdrom     = flag.StringP("cdrom", "r", "", "ESXi CDROM boot device")
	diskDev   = flag.StringP("device", "d", "", "ESXi disk boot device")
	partition = flag.IntP("partition", "p", 5, "ESXi boot partition")
)

func main() {
	flag.Parse()
	if *diskDev == "" && *cfg == "" && *cdrom == "" {
		log.Printf("Either --config, --device, or --cdrom must not be empty")
		flag.PrintDefaults()
		os.Exit(1)
	}

	var mi *boot.MultibootImage
	var err error
	if len(*cfg) > 0 {
		mi, err = esxi.LoadConfig(*cfg)
	} else if len(*cdrom) > 0 {
		// This is where the ESXi disk will be mounted.
		mountPoint, xerr := ioutil.TempDir("", "esxicdrom")
		if xerr != nil {
			log.Fatal(xerr)
		}

		mi, err = esxi.LoadCDROM(mountPoint, *cdrom)
	} else {
		// This is where the ESXi disk will be mounted.
		mountPoint, xerr := ioutil.TempDir("", "esxidisk")
		if xerr != nil {
			log.Fatal(xerr)
		}

		mi, err = esxi.LoadOS(mountPoint, *diskDev, *partition)
	}
	if err != nil {
		log.Fatalf("Failed to find ESXi: %v", err)
	}

	if err := mi.Load(false); err != nil {
		log.Fatalf("Failed to load ESXi into memory: %v", err)
	}
	if err := boot.Execute(); err != nil {
		log.Fatalf("Failed to boot image: %v", err)
	}
}
