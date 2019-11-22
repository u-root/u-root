// Copyright 2012-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// boot allows to handover a system running linuxboot/u-root
// to a legacy preinstalled operating system by replacing the traditional
// bootloader path

//
// Synopsis:
//	boot [-dev][-v][-dry-run]
//
// Description:
//	If returns to u-root shell, the code didn't found a local bootable option
//      -dev glob to use; default is /sys/class/block/*
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
	"encoding/binary"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/u-root/u-root/pkg/boot"
	"github.com/u-root/u-root/pkg/boot/grub"
	"github.com/u-root/u-root/pkg/boot/syslinux"
	"github.com/u-root/u-root/pkg/cmdline"
	"github.com/u-root/u-root/pkg/curl"
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
	devGlob           = flag.String("dev", "/sys/block/*", "Glob for devices")
	verbose           = flag.Bool("v", false, "Print debug messages")
	debug             = func(string, ...interface{}) {}
	dryRun            = flag.Bool("dry-run", false, "load kernel, but don't kexec it")
	defaultBoot       = flag.String("boot", "", "entry to boot (default to the configuration file default)")
	list              = flag.Bool("list", false, "list found configurations")
	uroot             string
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

// checkForBootableMBR is looking for bootable MBR signature
// Current support is limited to Hard disk devices and USB devices
func checkForBootableMBR(path string) error {
	var sig uint16
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	if err := binary.Read(io.NewSectionReader(f, signatureOffset, 2), binary.LittleEndian, &sig); err != nil {
		return err
	}
	if sig != bootableMBR {
		err := fmt.Errorf("%v is not a bootable device", path)
		return err
	}
	return nil
}

// getSupportedFilesystem returns all block file system supported by the linuxboot kernel
func getSupportedFilesystem() ([]string, error) {
	var err error
	fs, err := ioutil.ReadFile("/proc/filesystems")
	if err != nil {
		return nil, err
	}
	var returnValue []string
	for _, f := range strings.Split(string(fs), "\n") {
		n := strings.Fields(f)
		if len(n) != 1 {
			continue
		}
		returnValue = append(returnValue, n[0])
	}
	return returnValue, err

}

// mountEntry tries to mount a specific block device using a list of
// supported file systems. We have to try to mount the device
// itself, since devices can be filesystem formatted but not
// partitioned; and all its partitions.
func mountEntry(d string, supportedFilesystem []string) error {
	var err error
	debug("Try to mount %v", d)

	// find or create the mountpoint.
	m := filepath.Join(uroot, d)
	if _, err = os.Stat(m); err != nil && os.IsNotExist(err) {
		err = os.MkdirAll(m, 0777)
	}
	if err != nil {
		debug("Can't make %v", m)
		return err
	}
	for _, filesystem := range supportedFilesystem {
		var flags = uintptr(syscall.MS_RDONLY)
		debug("\twith %v", filesystem)
		if err := syscall.Mount(d, m, filesystem, flags, ""); err == nil {
			return err
		}
	}
	debug("No mount succeeded")
	return fmt.Errorf("Unable to mount any partition on %v", d)
}

func umountEntry(n string) error {
	return syscall.Unmount(n, syscall.MNT_DETACH)
}

// Localboot tries to boot from any local filesystem by parsing grub configuration
func Localboot() error {
	fs, err := getSupportedFilesystem()
	if err != nil {
		return errors.New("No filesystem support found")
	}
	debug("Supported filesystems: %v", fs)
	sysList, err := filepath.Glob(*devGlob)
	if err != nil {
		return errors.New("No available block devices to boot from")
	}
	// The Linux /sys file system is a bit, er, awkward. You can't find
	// the device special in there; just everything else.
	var blkList []string
	for _, b := range sysList {
		blkList = append(blkList, filepath.Join("/dev", filepath.Base(b)))
	}

	// We must validate if the MBR is bootable or not and keep the
	// devices which do have such support drive are easy to
	// detect.  This whole loop is pretty bogus at present, it
	// assumes the first partiton we find with grub.cfg is the one
	// we want. It works for now but ...
	var allparts []string
	for _, d := range blkList {
		err := checkForBootableMBR(d)
		if err != nil {
			// Not sure it matters; there can be many bogus entries?
			debug("MBR for %s failed: %v", d, err)
			continue
		}
		debug("Bootable device %v found", d)
		// You can't just look for numbers to match. Consider names like
		// mmcblk0, where has parts like mmcblk0p1. Just glob.
		g := d + "*"
		all, err := filepath.Glob(g)
		if err != nil {
			log.Printf("Glob for all partitions of %s failed: %v", g, err)
		}
		allparts = append(allparts, all...)
	}
	uroot, err = ioutil.TempDir("", "u-root-boot")
	if err != nil {
		return fmt.Errorf("Can't create tmpdir: %v", err)
	}
	debug("Trying to boot from %v", allparts)
	for _, d := range allparts {
		if err := mountEntry(d, fs); err != nil {
			continue
		}
		debug("mount succeed")
		u := filepath.Join(uroot, d)
		wd := &url.URL{
			Scheme: "file",
			Path:   u,
		}

		img, err := GrubBootImage(curl.DefaultSchemes, wd, *defaultBoot, *list)
		if err != nil {
			log.Printf("GrubBootImage failed: %v", err)
			// not grub config found, try isolinux
			img, err = IsolinuxBootImage(curl.DefaultSchemes, wd)
		}
		if err != nil {
			log.Printf("IsolinuxBootImage failed: %v", err)
			if err := umountEntry(u); err != nil {
				log.Printf("Can't unmount %v: %v", u, err)
			}
			continue
		}
		if li, ok := img.(*boot.LinuxImage); ok {
			// Filter the kernel command line
			li.Cmdline = updateBootCmdline(li.Cmdline)
		}
		log.Printf("BootImage: %s", img)
		if err := img.Load(*verbose); err != nil {
			return fmt.Errorf("kexec load of %v failed: %v", img, err)
		}

		if err := umountEntry(u); err != nil {
			log.Printf("Can't unmount %v: %v", u, err)
		}
		if *dryRun {
			continue
		}

		if err := boot.Execute(); err != nil {
			log.Printf("boot.Execute of %v failed: %v", u, err)
		}
		// TODO: We should probably return a real error.
		return nil
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

// use syslinux parser

func probeIsolinuxFiles() []string {
	files := make([]string, 0, 10)
	// search order from the syslinux wiki
	// http://wiki.syslinux.org/wiki/index.php?title=Config
	// TODO: do we want to handle extlinux too ?
	dirs := []string{
		"boot/isolinux",
		"isolinux",
		"boot/syslinux",
		"syslinux",
		"",
	}
	confs := []string{
		"isolinux.cfg",
		"syslinux.cfg",
	}
	for _, dir := range dirs {
		for _, conf := range confs {
			if dir == "" {
				files = append(files, conf)
			} else {
				files = append(files, path.Join(dir, conf))
			}
		}
	}
	return files
}

func IsolinuxParseConfigWithSchemes(workingDir *url.URL, s curl.Schemes) (*syslinux.Config, error) {
	for _, relname := range probeIsolinuxFiles() {
		c, err := syslinux.ParseConfigFileWithSchemes(s, relname, workingDir)
		if curl.IsURLError(err) {
			continue
		}
		return c, err
	}
	return nil, fmt.Errorf("no valid syslinux config found")
}

// call IsolinuxBootImage(curl.DefaultSchemes, dir)
func IsolinuxBootImage(schemes curl.Schemes, workingDir *url.URL) (*boot.LinuxImage, error) {
	pc, err := IsolinuxParseConfigWithSchemes(workingDir, schemes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse pxelinux config: %v", err)
	}

	label := pc.Entries[pc.DefaultEntry]
	return label, nil
}

// grub parser

var probeGrubFiles = []string{
	"boot/grub/grub.cfg",
	"grub/grub.cfg",
	"grub2/grub.cfg",
}

func GrubParseConfigWithSchemes(workingDir *url.URL, s curl.Schemes) (*grub.Config, error) {
	for _, relname := range probeGrubFiles {
		c, err := grub.ParseConfigFileWithSchemes(s, relname, workingDir)
		if curl.IsURLError(err) {
			continue
		}
		return c, err
	}
	return nil, fmt.Errorf("no valid grub config found")
}

// GrubBootImage
func GrubBootImage(schemes curl.Schemes, workingDir *url.URL, entryID string, list bool) (boot.OSImage, error) {
	pc, err := GrubParseConfigWithSchemes(workingDir, schemes)
	if err != nil && err != grub.ErrDefaultEntryNotFound {
		return nil, fmt.Errorf("failed to parse grub config: %v", err)
	}
	if list {
		fmt.Printf("%v\n", pc.Entries)
	}
	var entry boot.OSImage
	if entryID == "" {
		if err == grub.ErrDefaultEntryNotFound {
			return nil, err
		}
		entry = pc.Entries[pc.DefaultEntry]
	} else {
		var ok bool
		entry, ok = pc.Entries[entryID]
		if !ok {
			return nil, fmt.Errorf("entry not found in grub config: %v", entryID)
		}
	}
	return entry, nil
}
