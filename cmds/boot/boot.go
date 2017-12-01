// Copyright 2012-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// boot allows to handover a system running linuxboot/u-root
// to a legacy preinstalled operating system by replacing the traditional
// bootloader path

//
// Synopsis:
//	boot
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

var verbose bool

type options struct {
	verbose bool
}

func blkDevicesList(blkpath string, devpath string) ([]string, error) {
	var blkDevices []string
	files, err := ioutil.ReadDir(blkpath)
	if err != nil {
		return blkDevices, err
	}
	for _, file := range files {
		check, err := os.Stat(blkpath + file.Name() + devpath)
		if check == nil {
			continue
		}
		if err != nil {
			continue
		}
		deviceEntry, err := ioutil.ReadDir(blkpath + file.Name() + devpath)
		if err != nil {
			if verbose {
				log.Printf("can t read directory")
			}
			continue
		}
		blkDevices = append(blkDevices, deviceEntry[0].Name())
	}
	return blkDevices, nil
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
		err := errors.New("Not a bootable device")
		return err
	}
	return nil
}

// getDevicePartList returns all devices attached to a specific name like /dev/sdaX where X can move from 0 to 127
// FIXME no support for devices which are included into subdirectory within /dev
func getDevicePartList(path string) ([]string, error) {
	var returnValue []string
	files, err := ioutil.ReadDir("/dev/")
	if err != nil {
		return returnValue, err
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
	return returnValue, nil
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

// mountEntry tries to mount a specific block device
func mountEntry(path string, supportedFilesystem []string) (bool, error) {
	var returnValue bool
	var err error
	exist, err := os.Stat("/u-root")
	if exist == nil {
		err = syscall.Mkdir("/u-root", 0777)
		if err != nil {
			return false, err
		}
	}
	var flags uintptr
	// Was supposed to be unecessary for kernel 4.x.x
	if verbose {
		log.Printf("/dev/" + path)
	}
	for _, filesystem := range supportedFilesystem {
		flags = syscall.MS_MGC_VAL
		// Need to load the filesystem kind supported
		exist, err = os.Stat("/u-root/" + path)
		if exist == nil {
			err = syscall.Mkdir("/u-root/"+path, 0777)
			if err != nil {
				return false, err
			}
		}
		err := syscall.Mount("/dev/"+path, "/u-root/"+path, filesystem, flags, "")
		if err == nil {
			return true, nil
		}
	}
	returnValue = false
	return returnValue, nil
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
	if verbose {
		log.Printf(grubConfPath)
	}
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
	if verbose {
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
	if verbose {
		log.Printf(localKernelPath)
	}
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
	if err := kexec.Reboot(); err != nil {
		return err
	}
	return err

}

// init parse input parameters
func init() {
	flag.CommandLine.BoolVar(&verbose, "v", false, "Set verbose output")
}

func main() {
	flag.Parse()

	supportedFilesystem, err := getSupportedFilesystem()
	if err != nil {
		log.Panic("No filesystem support found")
	}
	if verbose {
		log.Printf("************** Supported Filesystem by current linuxboot ********************")
		for _, filesystem := range supportedFilesystem {
			log.Printf(filesystem)
		}
		log.Printf("*****************************************************************************")
	}
	blkList, err := blkDevicesList("/sys/dev/block/", "/device/block/")
	if err != nil {
		log.Panic("No available block devices to boot from")
	}
	// We must validate if the MBR is bootable or not and keep the
	// devices which do have such support
	// drive are easy to detect
	for _, entry := range blkList {
		dev := filepath.Join("/dev", entry)
		err := checkForBootableMBR(dev)
		if err != nil {
			// Not sure it matters; there can be many bogus entries?
			log.Printf("MBR for %s failed: %v", dev, err)
			continue
		}
		fmt.Println("Bootable device found")
		// We need to loop on the device entries which are into /dev/<device>X
		// and mount each partitions as to find /boot entry if it is available somewhere
		var devicePartList []string
		devicePartList, err = getDevicePartList(entry)
		if err != nil {
			continue
		}
		for _, deviceList := range devicePartList {
			mount, err := mountEntry(deviceList, supportedFilesystem)
			if err != nil {
				continue
			}
			if mount {
				if verbose {
					log.Printf("mount succeed")
				}
				var grubContent, grubConfPath = checkBootEntry("/u-root/" + deviceList)
				if grubConfPath != "" {
					if verbose {
						log.Printf("calling basic kexec")
					}
					err = kexecEntry(grubConfPath, grubContent, "/u-root/"+deviceList)
					if err != nil {
						log.Fatal("kexec failed")
					}
				}
			}
			umountEntry("/u-root/" + deviceList)
		}
	}
	log.Printf("Sorry no bootable device found")
}
