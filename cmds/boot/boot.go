// Copyright 2012-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// boot allows to handover a system running linuxboot/u-root
// to a legacy preinstalled operating system by replacing the traditional
// bootloader path

//
// Synopsis:
//	boot [-dev][-v][-dryrun]
//
// Description:
//	If returns to u-root shell, the code didn't found a local bootable option
//      -dev glob to use; default is /sys/class/block/*
//      -v prints messages
//      -dryrun doesn't really boot
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
	"flag"
	"fmt"
	"github.com/u-root/u-root/pkg/kexec"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
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
	devGlob     = flag.String("dev", "/sys/block/*", "Glob for devices")
	v           = flag.Bool("v", false, "Print debug messages")
	verbose     = func(string, ...interface{}) {}
	dryrun      = flag.Bool("dryrun", true, "Boot")
	defaultBoot = flag.String("boot", "default", "Default entry to boot")
	uroot       string
)

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
	verbose("Try to mount %v", d)

	// find or create the mountpoint.
	m := filepath.Join(uroot, d)
	if _, err = os.Stat(m); err != nil && os.IsNotExist(err) {
		err = os.MkdirAll(m, 0777)
	}
	if err != nil {
		verbose("Can't make %v", m)
		return err
	}
	for _, filesystem := range supportedFilesystem {
		var flags = uintptr(syscall.MS_RDONLY)
		verbose("\twith %v", filesystem)
		if err := syscall.Mount(d, m, filesystem, flags, ""); err == nil {
			return err
		}
	}
	verbose("No mount succeeded")
	return fmt.Errorf("Unable to mount any partition on %v", d)
}

func umountEntry(n string) error {
	return syscall.Unmount(n, syscall.MNT_DETACH)
}

// loadISOLinux reads an isolinux.cfg file. It further needs to
// process include directives.
// The include files are correctly placed in line at the place they
// were included. Include files can include other files, i.e. this
// function can recurse.
// It returns a []string of all the files, the path to the directory
// contain the boot images (bzImage, initrd, etc.) and the mount point.
// For now, these are the same, but in, e.g., grub, they are different,
// and they might in future change here too.
func loadISOLinux(dir, base string) ([]string, string, string, error) {
	sol := filepath.Join(dir, base)
	isolinux, err := ioutil.ReadFile(sol)
	if err != nil {
		return nil, "", "", err
	}

	// it's easier we think to do include processing here.
	lines := strings.Split(string(isolinux), "\n")
	var result []string
	for _, l := range lines {
		f := strings.Fields(l)
		if len(f) == 0 {
			continue
		}
		switch f[0] {
		default:
			result = append(result, l)
		case "include":
			i, _, _, err := loadISOLinux(dir, f[1])
			if err != nil {
				return nil, "", "", err
			}
			result = append(result, i...)
			// If the string contain include
			// we must call back as to read the content of the file
		}
	}
	return result, dir, dir, nil
}

// checkBootEntry is looking for grub.cfg file
// and return absolute path to it. It returns a []string with all the commands,
// the path to the direction containing files to load for kexec, and the mountpoint.
func checkBootEntry(mountPoint string) ([]string, string, string, error) {
	grub, err := ioutil.ReadFile(filepath.Join(mountPoint, "boot/grub/grub.cfg"))
	if err == nil {
		return strings.Split(string(grub), "\n"), filepath.Join(mountPoint, "/boot"), mountPoint, nil
	}
	return loadISOLinux(filepath.Join(mountPoint, "isolinux"), "isolinux.cfg")
}

// getFileMenuContent is parsing a grub.cfg file
// output: bootEntries
// grub parsing is a good deal messier than it looks.
// This is a simple parser will fail on anything tricky.
func getBootEntries(lines []string) (map[string]*bootEntry, map[string]string, error) {
	// There are two ways to reference a bootEntry: its order in the file and its
	// name.
	var (
		lineno      int
		line        string
		err         error
		curEntry    string
		numEntry    int
		bootEntries = make(map[string]*bootEntry)
		bootVars    = make(map[string]string)
		f           []string
		be          *bootEntry
	)
	defer func() {
		switch err := recover().(type) {
		case nil:
		case error:
			log.Fatalf("Bummer: %v, line #%d, line %q, fields %q", err, lineno, line, f)
		default:
			log.Fatalf("unexpected panic value: %T(%v)", err, err)
		}
	}()

	verbose("getBootEntries: %s", lines)
	curLabel := ""
	for _, line = range lines {
		lineno++
		f = strings.Fields(line)
		verbose("%d: %q, %q", lineno, line, f)
		if len(f) == 0 {
			continue
		}
		switch f[0] {
		default:
		case "#":
		case "default":
			bootVars["default"] = f[1]
		case "set":
			vals := strings.SplitN(f[1], "=", 2)
			if len(vals) == 2 {
				bootVars[vals[0]] = vals[1]
			}
		case "menuentry", "label":
			verbose("Menuentry %v", line)
			// nasty, but the alternatives are not much better.
			// grub config language is kind of arbitrary
			curEntry = fmt.Sprintf("\"%s\"", strconv.Itoa(numEntry))
			be = &bootEntry{}
			curLabel = f[1]
			bootEntries[f[1]] = be
			bootEntries[curEntry] = be
			numEntry++
		case "menu":
			// keyword default marks this as default
			if f[1] != "default" {
				continue
			}
			bootVars["default"] = curLabel
		case "linux", "kernel":
			verbose("linux %v", line)
			be.kernel = f[1]
			be.cmdline = strings.Join(f[2:], " ")
		case "initrd":
			verbose("initrd %v", line)
			be.initrd = f[1]
		case "append":
			// Format is a little bit strange
			//   append   MENU=/bin/cdrom-checker-menu vga=788 initrd=/install/initrd.gz quiet --
			var current_parameter int
			for _, parameter := range f {
				if current_parameter > 0 {
					if strings.HasPrefix(parameter, "initrd") {
						initrd := strings.Split(parameter, "=")
						bootEntries[curEntry].initrd = bootEntries[curEntry].initrd + initrd[1]
					} else {
						if parameter != "--" {
							bootEntries[curEntry].cmdline = bootEntries[curEntry].cmdline + " " + parameter
						}
					}
				}
				current_parameter++
			}
		}
	}
	verbose("grub config menu decoded to [%q, %q, %v]", bootEntries, bootVars, err)
	return bootEntries, bootVars, err

}

func copyLocal(path string) (string, error) {
	var dest string
	var err error
	result := strings.Split(path, "/")
	for _, entry := range result {
		dest = entry
	}
	dest = "/tmp/" + dest
	srcFile, err := os.Open(path)
	if err != nil {
		return dest, err
	}

	destFile, err := os.Create(dest) // creates if file doesn't exist
	if err != nil {
		return dest, err
	}

	_, err = io.Copy(destFile, srcFile) // check first var for number of bytes copied
	if err != nil {
		return dest, err
	}

	err = destFile.Sync()
	if err != nil {
		return dest, err
	}
	err = destFile.Close()
	if err != nil {
		return dest, err
	}
	err = srcFile.Close()
	if err != nil {
		return dest, err
	}
	return dest, nil
}

// kexecLoad Loads a new kernel and initrd.
func kexecLoad(grubConfPath string, grub []string, mountPoint string) error {
	verbose("kexecEntry: boot from %v", grubConfPath)
	b, v, err := getBootEntries(grub)
	if err != nil {
		return err
	}

	verbose("Boot Entries: %q", b)
	entry, ok := v[*defaultBoot]
	if !ok {
		return fmt.Errorf("Entry %v not found in config file", *defaultBoot)
	}
	be, ok := b[entry]
	if !ok {
		return fmt.Errorf("Entry %v not found in boot entries file", entry)
	}

	verbose("Boot params: %q", be)
	localKernelPath, err := copyLocal(filepath.Join(mountPoint, be.kernel))
	if err != nil {
		verbose("copyLocal(%v, %v): %v", filepath.Join(mountPoint, be.kernel), err)
		return err
	}
	localInitrdPath, err := copyLocal(filepath.Join(mountPoint, be.initrd))
	if err != nil {
		verbose("copyLocal(%v, %v): %v", filepath.Join(mountPoint, be.initrd), err)
		return err
	}
	verbose(localKernelPath)

	// We can kexec the kernel with localKernelPath as kernel entry, kernelParameter as parameter and initrd as initrd !
	log.Printf("Loading %s for kernel\n", localKernelPath)

	kernelDesc, err := os.OpenFile(localKernelPath, os.O_RDONLY, 0)
	if err != nil {
		verbose("%v", err)
		return err
	}
	// defer kernelDesc.Close()

	log.Printf("Loading %s for initramfs", localInitrdPath)
	ramfs, err := os.OpenFile(localInitrdPath, os.O_RDONLY, 0)
	if err != nil {
		verbose("%v", err)
		return err
	}
	// defer ramfs.Close()

	if err := kexec.FileLoad(kernelDesc, ramfs, be.cmdline); err != nil {
		verbose("%v", err)
		return err
	}
	if err = ramfs.Close(); err != nil {
		verbose("%v", err)
		return err
	}
	if err = kernelDesc.Close(); err != nil {
		verbose("%v", err)
		return err
	}
	return err

}

func main() {
	flag.Parse()

	if *v {
		verbose = log.Printf
	}
	fs, err := getSupportedFilesystem()
	if err != nil {
		log.Panic("No filesystem support found")
	}
	verbose("Supported filesystems: %v", fs)
	sysList, err := filepath.Glob(*devGlob)
	if err != nil {
		log.Panic("No available block devices to boot from")
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
			log.Printf("MBR for %s failed: %v", d, err)
			continue
		}
		verbose("Bootable device %v found", d)
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
		log.Fatalf("Can't create tmpdir: %v", err)
	}
	verbose("Trying to boot from %v", allparts)
	for _, d := range allparts {
		if err := mountEntry(d, fs); err != nil {
			continue
		}
		verbose("mount succeed")
		u := filepath.Join(uroot, d)
		config, fileDir, root, err := checkBootEntry(u)
		if err != nil {
			verbose("d: %v", d, err)
			if err := umountEntry(u); err != nil {
				log.Printf("Can't unmount %v: %v", u, err)
			}
			continue
		}
		verbose("calling basic kexec: content %v, path %v", config, fileDir)
		if err = kexecLoad(fileDir, config, root); err != nil {
			log.Printf("kexec on %v failed: %v", u, err)
		}
		verbose("kexecLoad succeeded")

		if err := umountEntry(u); err != nil {
			log.Printf("Can't unmount %v: %v", u, err)
		}
		if *dryrun {
			continue
		}
		if err := kexec.Reboot(); err != nil {
			log.Printf("Kexec Reboot %v failed, %v. Sorry", u, err)
		}
	}
	log.Fatalf("Sorry no bootable device found")
}
