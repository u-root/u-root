// Copyright 2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// esxiboot executes ESXi kernel over the running kernel.
//
// Synopsis:
//
//	esxiboot [-d --device] [-c --config] [-r --cdrom]
//
// Description:
//
//	Loads and executes ESXi kernel.
//
// Options:
//
//	--config=FILE or -c=FILE: set the ESXi config
//	--device=FILE or -d=FILE: set an ESXi disk to boot from
//	--cdrom=FILE or -r=FILE: set an ESXI CDROM to boot from
//	--append: append kernel cmdline arguments
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
	"log"
	"os"
	"strings"

	flag "github.com/spf13/pflag"
	"github.com/u-root/u-root/pkg/boot"
	"github.com/u-root/u-root/pkg/boot/esxi"
	"github.com/u-root/u-root/pkg/boot/multiboot"
	"github.com/u-root/u-root/pkg/mount"
)

var (
	cfg           = flag.StringP("config", "c", "", "ESXi config file")
	cdrom         = flag.StringP("cdrom", "r", "", "ESXi CDROM boot device")
	diskDev       = flag.StringP("device", "d", "", "ESXi disk boot device")
	appendCmdline = flag.StringArray("append", nil, "Arguments to append to kernel cmdline")
	dryRun        = flag.Bool("dry-run", false, "dry run (just mount + load the kernel, don't kexec)")
)

func main() {
	flag.Parse()
	if *diskDev == "" && *cfg == "" && *cdrom == "" {
		log.Printf("Either --config, --device, or --cdrom must be specified")
		flag.PrintDefaults()
		os.Exit(1)
	}

	if len(*diskDev) > 0 {
		imgs, mps, err := esxi.LoadDisk(*diskDev)
		if err != nil {
			log.Fatalf("Failed to load ESXi configuration: %v", err)
		}

		loaded := false
		for _, img := range imgs {
			if len(*appendCmdline) > 0 {
				img.Cmdline = img.Cmdline + " " + strings.Join(*appendCmdline, " ")
			}
			if err := img.Load(); err != nil {
				log.Printf("Failed to load ESXi image (%v) into memory: %v", img, err)
			} else {
				log.Printf("Loaded image: %v", img)
				// We loaded one, that's it.
				loaded = true
				break
			}
		}
		for _, mp := range mps {
			if err := mp.Unmount(mount.MNT_DETACH); err != nil {
				log.Printf("Failed to unmount %s: %v", mp, err)
			}
		}
		if !loaded {
			log.Fatalf("Failed to load all ESXi images found.")
		}
	} else {
		var err error
		var img *multiboot.Image
		var mp *mount.MountPoint
		if len(*cfg) > 0 {
			img, err = esxi.LoadConfig(*cfg)
		} else if len(*cdrom) > 0 {
			img, mp, err = esxi.LoadCDROM(*cdrom)
		}
		if err != nil {
			log.Fatalf("Failed to load ESXi configuration: %v", err)
		}
		if len(*appendCmdline) > 0 {
			img.Cmdline = img.Cmdline + " " + strings.Join(*appendCmdline, " ")
		}
		if err := img.Load(); err != nil {
			log.Fatalf("Failed to load ESXi image (%v) into memory: %v", img, err)
		}
		log.Printf("Loaded image: %v", img)
		if mp != nil {
			if err := mp.Unmount(mount.MNT_DETACH); err != nil {
				log.Printf("Failed to unmount %s: %v", mp, err)
			}
		}
	}

	if *dryRun {
		log.Printf("Dry run: not booting kernel.")
		os.Exit(0)
	}
	if err := boot.Execute(); err != nil {
		log.Fatalf("Failed to boot image: %v", err)
	}
}
