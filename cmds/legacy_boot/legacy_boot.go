// Copyright 2012-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// legacy_boot allows to handover a system running linuxboot/u-root
// to a legacy preinstalled operating system by replacing the traditional
// bootloader path

//
// Synopsis:
//	legacy_boot
//
// Description:
//	If returns to u-root shell, the code didn't found a local bootable option
//
// Notes:
//	The code is looking for boot/grub/grub.cfg file as to identify the
//	boot option.
//	The first bootable device found in the block device tree is the one used
//	Windows is not supported (that is a work in progress)
//
// Example:
//	legacy_boot -v 	- Start the script in verbose mode for debugging purpose

package main

import (
	"bufio"
	"fmt"
	"github.com/u-root/u-root/pkg/kexec"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"flag"
	"syscall"
)

var verbose bool
type options struct {
        verbose bool
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func blkDevicesList(blkpath string, devpath string) []string {
	var blkDevices []string
	files, err := ioutil.ReadDir(blkpath)
	if err != nil {
		log.Fatal(err)
	}
	for _, file := range files {
		deviceEntry, err := ioutil.ReadDir(blkpath + file.Name() + devpath)
		if err != nil {
			println("can t read directory")
		}
		blkDevices = append(blkDevices, deviceEntry[0].Name())
	}
	return (blkDevices)
}

// checkForBootableMbr is looking for bootable MBR signature 
// Current support is limited to Hard disk devices and USB devices
func checkForBootableMbr(path string) int {
	var b511, b510 byte
	f, err := os.Open(path)
	check(err)
	b1 := make([]byte, 512)
	f.Read(b1)
	b511 = b1[511]
	b510 = b1[510]
	f.Close()
	if ((b511 == 0xaa) && (b510 == 0x55)) == true {
		return 1
	}
	return 0
}

// getDevicePartList returns all devices attached to a specific name like /dev/sdaX where X can move from 0 to 127
// FIXME no support for devices which are included into subdirectory within /dev
func getDevicePartList(path string) []string {
	var returnValue []string
	files, err := ioutil.ReadDir("/dev/")
	if err != nil {
		log.Fatal(err)
	}
	for _, file := range files {
		if strings.HasPrefix(file.Name(), path) {
			// We shall not return full device name
			if file.Name() != path {
				// We shall check that the remaining part is a number
				returnValue = append(returnValue, file.Name())
			}
		}
	}
	return returnValue
}

// getSupportedFilesystem returns all block file system supported by the linuxboot kernel
func getSupportedFilesystem() []string {
	var returnValue []string
	file, err := os.Open("/proc/filesystems")
	check(err)
	scanner := bufio.NewScanner(file)
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	var string1, string2 string
	for scanner.Scan() {
		string1 = ""
		string2 = ""
		fmt.Sscanf(scanner.Text(), "%s %s", &string1, &string2)
		if string2 == "" {
			if string1 != "" {
				returnValue = append(returnValue, string1)
			}
		}
	}
	err=file.Close()
	check(err)
	return returnValue

}

// mountEntry tries to mount a specific block device
func mountEntry(path string, supportedFilesystem []string) bool {
	var returnValue bool
	syscall.Mkdir("/u-root", 0777)
	var flags uintptr
	// Was supposed to be unecessary for kernel 4.x.x
	if verbose {
		println("/dev/" + path)
	}
	for _, filesystem := range supportedFilesystem {
		flags = syscall.MS_MGC_VAL
		// Need to load the filesystem kind supported
		syscall.Mkdir("/u-root/"+path, 0777)
		err := syscall.Mount("/dev/"+path, "/u-root/"+path, filesystem, flags, "")
		if err == nil {
			return true
		}
	}
	returnValue = false
	return returnValue
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
func checkBootEntry(mountPoint string) string {
	_, err := os.Stat(mountPoint + "/boot")
	if err == nil {
		// The boot directory is there
		_, err2 := os.Stat(mountPoint + "/boot" + "/grub")
		if err2 == nil {
			_, err3 := os.Stat(mountPoint + "/boot" + "/grub" + "/grub.cfg")
			if err3 == nil {
				// found
				return (mountPoint + "/boot" + "/grub")
			}
		}
	}
	return ""

}

// getFileMenuContent is parsing a grub.cfg file
// input: absolute directory path to grub.cfg
// output: Return a list of strings with the following format
//	 line[3*x] - menuconfig
//	 line[3*x+1] - linux kernel + boot options
// 	 line[3*x+2] - initrd
// and the default boot entry configured into grub.cfg
func getFileMenuContent(path string) ([]string, int) {
	var returnValue []string
	file, _ := os.Open(path + "/grub.cfg")
	defer file.Close()
	scanner := bufio.NewScanner(file)
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	var status int
	var intReturn int
	intReturn = 0
	status = 0
	// When status = 0 we are looking for a menu entry
	// When status = 1 we are looking for a linux entry
	// When status = 2 we are looking for a initrd entry
	var trimmedLine string
	for scanner.Scan() {
		trimmedLine = strings.TrimSpace(scanner.Text())
		trimmedLine = strings.Join(strings.Fields(trimmedLine), " ")
		if !strings.HasPrefix(trimmedLine, "#") {
			if (strings.HasPrefix(trimmedLine, "set default=") ) && (status == 0) {
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
	}
	return returnValue, intReturn

}

func copyLocal(path string) string {
	var dest string
	result := strings.Split(path, "/")
	for _, entry := range result {
		dest = entry
	}
	dest = "/tmp/" + dest
	srcFile, err := os.Open(path)
	check(err)
	defer srcFile.Close()

	destFile, err := os.Create(dest) // creates if file doesn't exist
	check(err)
	defer destFile.Close()

	_, err = io.Copy(destFile, srcFile) // check first var for number of bytes copied
	check(err)

	err = destFile.Sync()
	check(err)
	return dest
}

// kexecEntry is booting new kernel based on the content of grub.cfg
func kexecEntry(grubConfPath string, mountPoint string) {
	var fileMenuContent []string
	var entry int
	var localKernelPath string
	var localInitrdPath string
	if verbose {
		println(grubConfPath)
	}
	fileMenuContent, entry = getFileMenuContent(grubConfPath)
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
	if verbose {
		println("************** boot parameters  ********************")
		println(kernel)
		println(kernelParameter)
		println(initrd)
		println("****************************************************")
	}
	localKernelPath = copyLocal(mountPoint + kernel)
	localInitrdPath = copyLocal(mountPoint + initrd)
	if verbose {
		println(localKernelPath)
	}
	umountEntry(mountPoint)
	// We can kexec the kernel with localKernelPath as kernel entry, kernelParameter as parameter and initrd as initrd !
	log.Printf("Loading %s for kernel\n", localKernelPath)

	kernelDesc, err := os.OpenFile(localKernelPath, os.O_RDONLY, 0)
	if err != nil {
		log.Fatalf("open(%q): %v", localKernelPath, err)
	}
	defer kernelDesc.Close()

	var ramfs *os.File
	ramfs, err = os.OpenFile(localInitrdPath, os.O_RDONLY, 0)
	if err != nil {
		log.Fatalf("open(%q): %v", localInitrdPath, err)
	}
	defer ramfs.Close()

	if err := kexec.FileLoad(kernelDesc, ramfs, kernelParameter); err != nil {
		log.Fatalf("%v", err)
	}
	if err := kexec.Reboot(); err != nil {
		log.Fatalf("%v", err)
	}

}

func registerFlags(f *flag.FlagSet) *options {
        o := &options{}
        f.BoolVar(&o.verbose, "v", false, "Set verbose output")
        return o
}


func main() {
	var blkList []string
	var supportedFilesystem []string
	verbose = false
        opts := registerFlags(flag.CommandLine)
	flag.Parse()

        if opts.verbose != false {
		verbose = true
        }

	supportedFilesystem = getSupportedFilesystem()
	if verbose {
		println("************** Supported Filesystem by current linuxboot ********************")
		for _, filesystem := range supportedFilesystem {
			println(filesystem)
		}
		println("*****************************************************************************")
	}
	blkList = blkDevicesList("/sys/dev/block/", "/device/block/")
	// We must validate if the MBR is bootable or not and keep the
	// devices which do have such support
	// drive are easy to detect
	for _, entry := range blkList {
		if checkForBootableMbr("/dev/"+entry) == 1 {
			fmt.Println("Bootable device found")
			// We need to loop on the device entries which are into /dev/<device>X
			// and mount each partitions as to find /boot entry if it is available somewhere
			var devicePartList []string
			devicePartList = getDevicePartList(entry)
			for _, deviceList := range devicePartList {
				if mountEntry(deviceList, supportedFilesystem) {
					if verbose {
						println("mount succeed")
					}
					var grubConfPath = checkBootEntry("/u-root/" + deviceList)
					if grubConfPath != "" {
						if verbose {
							println("calling basic kexec")
						}
						kexecEntry(grubConfPath, "/u-root/"+deviceList)
					}
				}
				umountEntry("/u-root/" + deviceList)
			}
		}
	}
	println("Sorry no bootable device found")
}
