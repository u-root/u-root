// Copyright 2012-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// boot allows to handover a system running linuxboot/u-root
// to a legacy preinstalled operating system by replacing the traditional
// bootloader path

//
// Synopsis:
//	boot [-v][-dry-run]
//
// Description:
//	If returns to u-root shell, the code didn't found a local bootable option
//
//      -v prints messages
//      -dry-run doesn't really boot
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
	"github.com/u-root/u-root/pkg/boot/syslinux"
	"github.com/u-root/u-root/pkg/cmdline"
	"github.com/u-root/u-root/pkg/mount"
	"github.com/u-root/u-root/pkg/mount/block"
	"github.com/u-root/u-root/pkg/ulog"
)

const (
	bootableMBR     = 0xaa55
	signatureOffset = 510
)

type bootEntry struct {
	kernel  string
	initrd  string
	cmdline string
}

var (
	verbose     = flag.Bool("v", false, "Print debug messages")
	debug       = func(string, ...interface{}) {}
	dryRun      = flag.Bool("dry-run", false, "load kernel, but don't kexec it")
	defaultBoot = flag.String("boot", "", "entry to boot (default to the configuration file default)")
	list        = flag.Bool("list", false, "list found configurations")

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

func mountAndBoot(device *block.BlockDev, mountDir string) {
	os.MkdirAll(mountDir, 0777)

	mp, err := device.Mount(mountDir, mount.ReadOnly)
	if err != nil {
		return
	}
	defer mp.Unmount(mount.MNT_DETACH)

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

	if len(imgs) == 0 {
		return
	}

	// Boot just the first image.
	img := imgs[0]

	// Make changes to the kernel command line based on our cmdline.
	if li, ok := img.(*boot.LinuxImage); ok {
		li.Cmdline = updateBootCmdline(li.Cmdline)
	}

	log.Printf("BootImage: %s", img)
	if err := img.Load(*verbose); err != nil {
		log.Printf("kexec load of %v failed: %v", img, err)
		return
	}

	if err := mp.Unmount(mount.MNT_DETACH); err != nil {
		log.Printf("Can't unmount %v: %v", mp, err)
	}
	if *dryRun {
		return
	}

	if err := boot.Execute(); err != nil {
		log.Printf("boot.Execute of %v failed: %v", img, err)
	}

	// kexec was successful. kexec should have taken over. What happened?
	log.Fatalf("kexec boot returned success, but new kernel is not running...")
}

// Localboot tries to boot from any local filesystem by parsing grub configuration
func Localboot() error {
	blockDevs, err := block.GetBlockDevices()
	if err != nil {
		return errors.New("no available block devices to boot from")
	}

	// Try to only boot from "good" block devices.
	blockDevs = blockDevs.FilterZeroSize()
	debug("Booting from the following block devices: %v", blockDevs)

	mountPoints, err := ioutil.TempDir("", "u-root-boot")
	if err != nil {
		return fmt.Errorf("Can't create tmpdir: %v", err)
	}
	defer os.RemoveAll(mountPoints)

	for _, device := range blockDevs {
		dir := filepath.Join(mountPoints, device.Name)
		mountAndBoot(device, dir)
	}
	return fmt.Errorf("Sorry no bootable device found")
}

func main() {
	flag.Parse()

	if *verbose {
		debug = log.Printf
	}

	if err := Localboot(); err != nil {
		log.Fatal(err)
	}
}
