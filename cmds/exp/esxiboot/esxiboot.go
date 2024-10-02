// Copyright 2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//go:build !tinygo || tinygo.enable

// esxiboot executes ESXi kernel over the running kernel.
//
// Synopsis:
//     esxiboot [-d --device] [-c --config] [-r --cdrom]
//
// Description:
//     Loads and executes ESXi kernel.
//
// Options:
//     --config=FILE or -c=FILE: set the ESXi config
//     --device=FILE or -d=FILE: set an ESXi disk to boot from
//     --cdrom=FILE or -r=FILE: set an ESXI CDROM to boot from
//     --append: append kernel cmdline arguments
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
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/u-root/u-root/pkg/boot"
	"github.com/u-root/u-root/pkg/boot/esxi"
	"github.com/u-root/u-root/pkg/mount"
	"github.com/u-root/u-root/pkg/uroot/unixflag"
)

type cmd struct {
	cfg           string
	cdrom         string
	diskDev       string
	appendCmdline []string
	dryRun        bool
}

func (c *cmd) run() error {
	if len(c.diskDev) > 0 {
		imgs, mps, err := esxi.LoadDisk(c.diskDev)
		if err != nil {
			return fmt.Errorf("failed to load ESXi configuration: %w", err)
		}

		loaded := false
		for _, img := range imgs {
			if len(c.appendCmdline) > 0 {
				img.Cmdline = img.Cmdline + " " + strings.Join(c.appendCmdline, " ")
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
			return fmt.Errorf("failed to load all ESXi images found")
		}
	} else {
		var err error
		var img *boot.MultibootImage
		var mp *mount.MountPoint
		if len(c.cfg) > 0 {
			img, err = esxi.LoadConfig(c.cfg)
		} else if len(c.cdrom) > 0 {
			img, mp, err = esxi.LoadCDROM(c.cdrom)
		}
		if err != nil {
			return fmt.Errorf("failed to load ESXi configuration: %w", err)
		}
		if len(c.appendCmdline) > 0 {
			img.Cmdline = img.Cmdline + " " + strings.Join(c.appendCmdline, " ")
		}
		if err := img.Load(); err != nil {
			return fmt.Errorf("failed to load ESXi image (%v) into memory: %w", img, err)
		}
		log.Printf("Loaded image: %v", img)
		if mp != nil {
			if err := mp.Unmount(mount.MNT_DETACH); err != nil {
				log.Printf("Failed to unmount %s: %v", mp, err)
			}
		}
	}

	if c.dryRun {
		log.Printf("Dry run: not booting kernel.")
		os.Exit(0)
	}
	if err := boot.Execute(); err != nil {
		return fmt.Errorf("failed to boot image: %w", err)
	}
	return nil
}

func command(args []string) *cmd {
	c := &cmd{}
	f := flag.NewFlagSet(args[0], flag.ExitOnError)
	f.StringVar(&c.cfg, "config", "", "ESXi config file")
	f.StringVar(&c.cfg, "c", "", "ESXi config file (shorthand)")

	f.StringVar(&c.cdrom, "cdrom", "", "ESXi CDROM boot device")
	f.StringVar(&c.cdrom, "r", "", "ESXi CDROM boot device (shorthand)")

	f.StringVar(&c.diskDev, "device", "", "ESXi disk boot device")
	f.StringVar(&c.diskDev, "d", "", "ESXi disk boot device (shorthand)")

	f.Var((*unixflag.StringArray)(&c.appendCmdline), "append", "Arguments to append to kernel cmdline")

	f.BoolVar(&c.dryRun, "dry-run", false, "dry run (just mount + load the kernel, don't kexec)")

	f.Parse(unixflag.ArgsToGoArgs(args[1:]))

	if c.diskDev == "" && c.cfg == "" && c.cdrom == "" {
		log.Printf("Either --config, --device, or --cdrom must be specified")
		f.PrintDefaults()
		os.Exit(1)
	}

	return c
}

func main() {
	if err := command(os.Args).run(); err != nil {
		log.Fatal(err)
	}
}
