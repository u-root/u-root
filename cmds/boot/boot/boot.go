// Copyright 2012-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// boot allows to handover a system running linuxboot/u-root
// to a legacy preinstalled operating system by replacing the traditional
// bootloader path

//
// Synopsis:
//	boot [-v][-no-load][-no-exec]
//
// Description:
//	If returns to u-root shell, the code didn't found a local bootable option
//
//      -v prints messages
//      -no-load prints the boot image paths it was going to load, but doesn't load + exec them
//      -no-exec loads the boot image, but doesn't exec it
//
// Notes:
//	The code is looking for boot/grub/grub.cfg file as to identify the
//	boot option.
//	The first bootable device found in the block device tree is the one used
//	Windows is not supported (that is a work in progress)
//
// Example:
//	boot -v 	- Start the script in verbose mode for debugging purpose

package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/u-root/u-root/pkg/boot"
	"github.com/u-root/u-root/pkg/boot/bls"
	"github.com/u-root/u-root/pkg/boot/grub"
	"github.com/u-root/u-root/pkg/boot/menu"
	"github.com/u-root/u-root/pkg/boot/syslinux"
	"github.com/u-root/u-root/pkg/cmdline"
	"github.com/u-root/u-root/pkg/mount"
	"github.com/u-root/u-root/pkg/mount/block"
	"github.com/u-root/u-root/pkg/ulog"
)

var (
	debug   = func(string, ...interface{}) {}
	verbose = flag.Bool("v", false, "Print debug messages")
	noLoad  = flag.Bool("no-load", false, "print chosen boot configuration, but do not load + exec it")
	noExec  = flag.Bool("no-exec", false, "load boot configuration, but do not exec it")

	removeCmdlineItem = flag.String("remove", "console", "comma separated list of kernel params value to remove from parsed kernel configuration (default to console)")
	reuseCmdlineItem  = flag.String("reuse", "console", "comma separated list of kernel params value to reuse from current kernel (default to console)")
	appendCmdline     = flag.String("append", "", "Additional kernel params")
)

// updateBootCmdline get the kernel command line parameters and filter it:
// it removes parameters listed in 'remove' and append extra parameters from
// the 'append' and 'reuse' flags
func updateBootCmdline(cl string) string {
	f := cmdline.NewUpdateFilter(*appendCmdline, strings.Split(*removeCmdlineItem, ","), strings.Split(*reuseCmdlineItem, ","))
	return f.Update(cl)
}

func parse(device *block.BlockDev, mountDir string) []boot.OSImage {
	imgs, err := bls.ScanBLSEntries(ulog.Log, mountDir)
	if err != nil {
		log.Printf("Failed to parse systemd-boot BootLoaderSpec configs, trying another format...: %v", err)
	}

	grubImgs, err := grub.ParseLocalConfig(context.Background(), mountDir)
	if err != nil {
		log.Printf("Failed to parse GRUB configs from %s, trying another format...: %v", device, err)
	}
	imgs = append(imgs, grubImgs...)

	syslinuxImgs, err := syslinux.ParseLocalConfig(context.Background(), mountDir)
	if err != nil {
		log.Printf("Failed to parse syslinux configs from %s: %v", device, err)
	}
	imgs = append(imgs, syslinuxImgs...)

	return imgs
}

// Localboot tries to boot from any local filesystem by parsing grub configuration
func Localboot() ([]boot.OSImage, []*mount.MountPoint, error) {
	blockDevs, err := block.GetBlockDevices()
	if err != nil {
		return nil, nil, errors.New("no available block devices to boot from")
	}

	// Try to only boot from "good" block devices.
	blockDevs = blockDevs.FilterZeroSize()
	debug("Booting from the following block devices: %v", blockDevs)

	mountPoints, err := ioutil.TempDir("", "u-root-boot")
	if err != nil {
		return nil, nil, fmt.Errorf("cannot create tmpdir: %v", err)
	}

	var images []boot.OSImage
	var mps []*mount.MountPoint
	for _, device := range blockDevs {
		dir := filepath.Join(mountPoints, device.Name)

		os.MkdirAll(dir, 0777)
		mp, err := device.Mount(dir, mount.ReadOnly)
		if err != nil {
			continue
		}

		imgs := parse(device, dir)
		images = append(images, imgs...)
		mps = append(mps, mp)
	}
	return images, mps, nil
}

func main() {
	flag.Parse()

	if *verbose {
		debug = log.Printf
	}

	images, mps, err := Localboot()
	if err != nil {
		log.Fatal(err)
	}
	for _, img := range images {
		// Make changes to the kernel command line based on our cmdline.
		if li, ok := img.(*boot.LinuxImage); ok {
			li.Cmdline = updateBootCmdline(li.Cmdline)
		}
	}

	if *noLoad {
		if len(images) > 0 {
			log.Printf("Got configuration: %s", images[0])
		} else {
			log.Fatalf("Nothing bootable found.")
		}
		return
	}
	menuEntries := menu.OSImages(*verbose, images...)
	menuEntries = append(menuEntries, menu.Reboot{})
	menuEntries = append(menuEntries, menu.StartShell{})

	chosenEntry := menu.ShowMenuAndLoad(os.Stdin, menuEntries...)

	// Clean up.
	for _, mp := range mps {
		if err := mp.Unmount(mount.MNT_DETACH); err != nil {
			debug("Failed to unmount %s: %v", mp, err)
		}
	}
	if chosenEntry == nil {
		log.Fatalf("Nothing to boot.")
	}
	if *noExec {
		log.Printf("Chosen menu entry: %s", chosenEntry)
		os.Exit(0)
	}
	// Exec should either return an error or not return at all.
	if err := chosenEntry.Exec(); err != nil {
		log.Fatalf("Failed to exec %s: %v", chosenEntry, err)
	}

	// Kexec should either return an error or not return.
	panic("unreachable")
}
