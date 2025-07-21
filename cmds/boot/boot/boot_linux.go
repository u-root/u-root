// Copyright 2012-2020 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build !tinygo || tinygo.enable

// Command boot allows to handover a system running linuxboot/u-root
// to a legacy preinstalled operating system by replacing the traditional
// bootloader path
//
// Synopsis:
//
//	boot [-v][-no-load][-no-exec]
//
// Description:
//
//	If returns to u-root shell, the code didn't found a local bootable option
//
//	-v prints messages
//	-no-load prints the boot image paths it was going to load, but doesn't load + exec them
//	-no-exec loads the boot image, but doesn't exec it
//
// Notes:
//
//	The code is looking for boot/grub/grub.cfg file as to identify the
//	boot option.
//	The first bootable device found in the block device tree is the one used
//	Windows is not supported (that is a work in progress)
//
// Example:
//
//	boot -v - Start the script in verbose mode for debugging purpose
package main

import (
	"flag"
	"log"
	"strings"

	"github.com/u-root/u-root/pkg/boot"
	"github.com/u-root/u-root/pkg/boot/bootcmd"
	"github.com/u-root/u-root/pkg/boot/localboot"
	"github.com/u-root/u-root/pkg/boot/menu"
	"github.com/u-root/u-root/pkg/cmdline"
	"github.com/u-root/u-root/pkg/mount"
	"github.com/u-root/u-root/pkg/mount/block"
	"github.com/u-root/u-root/pkg/ulog"
)

var (
	verbose = flag.Bool("v", false, "Print debug messages")
	noLoad  = flag.Bool("no-load", false, "print chosen boot configuration, but do not load + exec it")
	noExec  = flag.Bool("no-exec", false, "load boot configuration, but do not exec it")

	removeCmdlineItem = flag.String("remove", "console", "comma separated list of kernel params value to remove from parsed kernel configuration (default to console)")
	reuseCmdlineItem  = flag.String("reuse", "console", "comma separated list of kernel params value to reuse from current kernel (default to console)")
	appendCmdline     = flag.String("append", "", "Additional kernel params")
	blockList         = flag.String("block", "", "comma separated list of pci vendor and device ids to ignore (format vendor:device). E.g. 0x8086:0x1234,0x8086:0xabcd")
)

// updateBootCmdline get the kernel command line parameters and filter it:
// it removes parameters listed in 'remove' and append extra parameters from
// the 'append' and 'reuse' flags
func cmdlineModifier(li *boot.LinuxImage) {
	f := cmdline.NewUpdateFilter(*appendCmdline, strings.Split(*removeCmdlineItem, ","), strings.Split(*reuseCmdlineItem, ","))
	li.Cmdline = f.Update(cmdline.NewCmdLine(), li.Cmdline)
}

func main() {
	flag.Parse()

	if *verbose {
		block.Debug = log.Printf
	}
	blockDevs, err := block.GetBlockDevices()
	if err != nil {
		log.Fatal("No available block devices to boot from")
	}

	// Try to only boot from "good" block devices.
	blockDevs = blockDevs.FilterZeroSize()

	// Parse and filter blocklist
	if *blockList != "" {
		blockDevs, err = blockDevs.FilterBlockPCIString(*blockList)
		if err != nil {
			log.Fatal(err)
		}
	}

	log.Printf("Booting from the following block devices: %v", blockDevs)

	l := ulog.Null
	if *verbose {
		l = ulog.Log
	}
	mountPool := &mount.Pool{}
	images, err := localboot.Localboot(l, blockDevs, mountPool)
	if err != nil {
		log.Fatal(err)
	}
	// Make changes to the kernel command line based on our cmdline.
	boot.ApplyLinuxModifiers(images, cmdlineModifier)

	menuEntries := menu.OSImages(*verbose, images...)
	menuEntries = append(menuEntries, menu.Reboot{})
	menuEntries = append(menuEntries, menu.StartShell{})

	// Boot does not return.
	bootcmd.ShowMenuAndBoot(menuEntries, mountPool, *noLoad, *noExec)
}
