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
	"errors"
	"flag"
	"fmt"
	"github.com/u-root/u-root/pkg/kexec"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"syscall"
)

const (
	bootableMBR     = 0xaa55
	signatureOffset = 510
)

var (
	devGlob = flag.String("dev", "/sys/block/*", "Glob for devices")
	v       = flag.Bool("v", false, "Print debug messages")
	verbose = func(string, ...interface{}) {}
	dryrun  = flag.Bool("dryrun", true, "Boot")
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
		err := errors.New("Not a bootable device")
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
	m := filepath.Join("/u-root", d)
	if _, err = os.Stat(m); err != nil && os.IsNotExist(err) {
		err = os.MkdirAll(m, 0777)
	}
	if err != nil {
		return err
	}
	for _, filesystem := range supportedFilesystem {
		// Was supposed to be unecessary for kernel 4.x.x
		var flags = uintptr(syscall.MS_MGC_VAL)
		if err := syscall.Mount(d, m, filesystem, flags, ""); err == nil {
			return err
		}
	}
	return errors.New("Unable to mount any partition on " + d)
}

func umountEntry(path string) bool {
	var returnValue bool
	var flags int
	// Was supposed to be unecessary for kernel 4.x.x
	flags = syscall.MNT_DETACH
	err := syscall.Unmount(path, flags)
	if err == nil {
		return true
	}
	returnValue = false
	return returnValue
}

// checkBootEntry is looking for grub.cfg file
// and return absolute path to it
func checkBootEntry(mountPoint string) ([]byte, string) {
	grub, err := ioutil.ReadFile(filepath.Join(mountPoint, "/boot/grub/grub.cfg"))
	if err == nil {
		return grub, filepath.Join(mountPoint, "/boot/grub")
	}
	return grub, ""

}

// getFileMenuContent is parsing a grub.cfg file
// input: absolute directory path to grub.cfg
// output: Return a list of strings with the following format
//	 line[3*x] - menuconfig
//	 line[3*x+1] - linux kernel + boot options
// 	 line[3*x+2] - initrd
// and the default boot entry configured into grub.cfg
func getFileMenuContent(file []byte) ([]string, int, error) {
	var returnValue []string
	var err error
	var status int
	var intReturn int
	intReturn = 0
	status = 0
	// When status = 0 we are looking for a menu entry
	// When status = 1 we are looking for a linux entry
	// When status = 2 we are looking for a initrd entry
	var trimmedLine string
	s := string(file)
	for _, line := range strings.Split(s, "\n") {
		trimmedLine = strings.TrimSpace(line)
		trimmedLine = strings.Join(strings.Fields(trimmedLine), " ")
		if strings.HasPrefix(trimmedLine, "#") {
			continue
		}
		if (strings.HasPrefix(trimmedLine, "set default=")) && (status == 0) {
			fmt.Sscanf(trimmedLine, "set default=\"%d\"", &intReturn)
		}
		if (strings.HasPrefix(trimmedLine, "menuentry ")) && (status == 0) {
			status = 1
			returnValue = append(returnValue, trimmedLine)
		}
		if (strings.HasPrefix(trimmedLine, "linux ")) && (status == 1) {
			status = 2
			returnValue = append(returnValue, trimmedLine)
		}
		if (strings.HasPrefix(trimmedLine, "initrd ")) && (status == 2) {
			status = 0
			returnValue = append(returnValue, trimmedLine)
		}
	}
	return returnValue, intReturn, err

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

// kexecEntry is booting new kernel based on the content of grub.cfg
func kexecEntry(grubConfPath string, grub []byte, mountPoint string) error {
	var fileMenuContent []string
	var entry int
	var localKernelPath string
	var localInitrdPath string
	verbose(grubConfPath)
	fileMenuContent, entry, err := getFileMenuContent(grub)
	if err != nil {
		return err
	}
	var kernel string
	var kernelParameter string
	var initrd string
	var kernelInfos []string
	kernelInfos = strings.Fields(fileMenuContent[3*entry+1])
	kernel = kernelInfos[1]
	var count int
	count = 0
	for _, field := range kernelInfos {
		if count > 1 {
			kernelParameter = kernelParameter + " " + field
		}
		count = count + 1
	}
	fmt.Sscanf(fileMenuContent[3*entry+2], "initrd %s", &initrd)
	if *v {
		log.Printf("************** boot parameters  ********************")
		log.Printf(kernel)
		log.Printf(kernelParameter)
		log.Printf(initrd)
		log.Printf("****************************************************")
	}
	localKernelPath, err = copyLocal(mountPoint + kernel)
	if err != nil {
		return err
	}
	localInitrdPath, err = copyLocal(mountPoint + initrd)
	if err != nil {
		return err
	}
	verbose(localKernelPath)

	umountEntry(mountPoint)
	// We can kexec the kernel with localKernelPath as kernel entry, kernelParameter as parameter and initrd as initrd !
	log.Printf("Loading %s for kernel\n", localKernelPath)

	kernelDesc, err := os.OpenFile(localKernelPath, os.O_RDONLY, 0)
	if err != nil {
		return err
	}
	// defer kernelDesc.Close()

	var ramfs *os.File
	ramfs, err = os.OpenFile(localInitrdPath, os.O_RDONLY, 0)
	if err != nil {
		return err
	}
	// defer ramfs.Close()

	if err := kexec.FileLoad(kernelDesc, ramfs, kernelParameter); err != nil {
		return err
	}
	err = ramfs.Close()
	if err != nil {
		return err
	}
	err = kernelDesc.Close()
	if err != nil {
		return err
	}
	if *dryrun {
		return nil
	}
	if err := kexec.Reboot(); err != nil {
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
	verbose("Trying to boot from %v", allparts)
	for _, d := range allparts {
		if err := mountEntry(d, fs); err != nil {
			continue
		}
		verbose("mount succeed")
		u := filepath.Join("/u-root", d)
		var grubContent, grubConfPath = checkBootEntry(u)
		if grubConfPath != "" {
			verbose("calling basic kexec")
			if err = kexecEntry(grubConfPath, grubContent, u); err != nil {
				log.Fatalf("kexec failed: %v", err)
			}
		}

		umountEntry(u)
	}
	log.Fatalf("Sorry no bootable device found")
}
