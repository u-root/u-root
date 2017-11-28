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
	"syscall"
)

var verbose bool

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func blk_devices_list(blkpath string, devpath string) []string {
	var blk_devices []string
	files, err := ioutil.ReadDir(blkpath)
	if err != nil {
		log.Fatal(err)
	}
	for _, file := range files {
		device_entry, err := ioutil.ReadDir(blkpath + file.Name() + devpath)
		if err == nil {
			blk_devices = append(blk_devices, device_entry[0].Name())
		}
	}
	return (blk_devices)
}

// Current support is limited to Hard disk devices and USB devices

func check_for_bootable_mbr(path string) int {
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

// Return all devices attached to a specific name like /dev/sdaX where X can move from 0 to 127
// FIXME no support for devices which are included into subdirectory within /dev

func get_device_part_list(path string) []string {
	var return_value []string
	files, err := ioutil.ReadDir("/dev/")
	if err != nil {
		log.Fatal(err)
	}
	for _, file := range files {
		if strings.HasPrefix(file.Name(), path) {
			// We shall not return full device name
			if file.Name() != path {
				// We shall check that the remaining part is a number
				return_value = append(return_value, file.Name())
			}
		}
	}
	return return_value
}

// Return all block file system supported by the linuxboot kernel

func get_supported_filesystem() []string {
	var return_value []string
	file, err := os.Open("/proc/filesystems")
	check(err)
	defer file.Close()
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
				return_value = append(return_value, string1)
			}
		}
	}
	return return_value

}

// try to mount a specific block device

func mount_entry(path string, supported_filesystem []string) bool {
	var return_value bool
	syscall.Mkdir("/u-root", 0777)
	var flags uintptr
	// Was supposed to be unecessary for kernel 4.x.x
	if verbose {
		println("/dev/" + path)
	}
	for _, filesystem := range supported_filesystem {
		flags = syscall.MS_MGC_VAL
		// Need to load the filesystem kind supported
		syscall.Mkdir("/u-root/"+path, 0777)
		err := syscall.Mount("/dev/"+path, "/u-root/"+path, filesystem, flags, "")
		if err == nil {
			return true
		}
	}
	return_value = false
	return return_value
}

func umount_entry(path string) bool {
	var return_value bool
	var flags int
	// Was supposed to be unecessary for kernel 4.x.x
	flags = syscall.MNT_DETACH
	err := syscall.Unmount(path, flags)
	if err == nil {
		return true
	}
	return_value = false
	return return_value
}

// This function is looking for grub.cfg file
// and return absolute path to it

func check_boot_entry(mount_point string) string {
	_, err := os.Stat(mount_point + "/boot")
	if err == nil {
		// The boot directory is there
		_, err2 := os.Stat(mount_point + "/boot" + "/grub")
		if err2 == nil {
			_, err3 := os.Stat(mount_point + "/boot" + "/grub" + "/grub.cfg")
			if err3 == nil {
				// found
				return (mount_point + "/boot" + "/grub")
			}
		}
	}
	return ""

}

// This function is parsing a grub.cfg file
// input: absolute directory path to grub.cfg
// output: Return a list of strings with the following format
//	 line[3*x] - menuconfig
//	 line[3*x+1] - linux kernel + boot options
// 	 line[3*x+2] - initrd
// and the default boot entry configured into grub.cfg

func get_file_menu_content(path string) ([]string, int) {
	var return_value []string
	file, _ := os.Open(path + "/grub.cfg")
	defer file.Close()
	scanner := bufio.NewScanner(file)
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	var status int
	var int_return int
	int_return = 0
	status = 0
	// When status = 0 we are looking for a menu entry
	// When status = 1 we are looking for a linux entry
	// When status = 2 we are looking for a initrd entry
	var trimmed_line string
	for scanner.Scan() {
		trimmed_line = strings.TrimSpace(scanner.Text())
		trimmed_line = strings.Join(strings.Fields(trimmed_line), " ")
		if strings.HasPrefix(trimmed_line, "#") == false {
			if (strings.HasPrefix(trimmed_line, "set default=") == true) && (status == 0) {
				fmt.Sscanf(trimmed_line, "set default=\"%d\"", &int_return)
			}
			if (strings.HasPrefix(trimmed_line, "menuentry ") == true) && (status == 0) {
				status = 1
				return_value = append(return_value, trimmed_line)
			}
			if (strings.HasPrefix(trimmed_line, "linux ") == true) && (status == 1) {
				status = 2
				return_value = append(return_value, trimmed_line)
			}
			if (strings.HasPrefix(trimmed_line, "initrd ") == true) && (status == 2) {
				status = 0
				return_value = append(return_value, trimmed_line)
			}
		}
	}
	return return_value, int_return

}

func copy_local(path string) string {
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

// This function is booting new kernel based on the content of grub.cfg

func kexec_entry(grub_conf_path string, mount_point string) {
	var file_menu_content []string
	var entry int
	var local_kernel_path string
	var local_initrd_path string
	if verbose {
		println(grub_conf_path)
	}
	file_menu_content, entry = get_file_menu_content(grub_conf_path)
	var kernel string
	var kernel_parameter string
	var initrd string
	var kernel_infos []string
	kernel_infos = strings.Fields(file_menu_content[3*entry+1])
	kernel = kernel_infos[1]
	var count int
	count = 0
	for _, field := range kernel_infos {
		if count > 1 {
			kernel_parameter = kernel_parameter + " " + field
		}
		count = count + 1
	}
	fmt.Sscanf(file_menu_content[3*entry+2], "initrd %s", &initrd)
	if verbose {
		println("************** boot parameters  ********************")
		println(kernel)
		println(kernel_parameter)
		println(initrd)
		println("****************************************************")
	}
	local_kernel_path = copy_local(mount_point + kernel)
	local_initrd_path = copy_local(mount_point + initrd)
	if verbose {
		println(local_kernel_path)
	}
	umount_entry(mount_point)
	// We can kexec the kernel with local_kernel_path as kernel entry, kernel_parameter as parameter and initrd as initrd !
	log.Printf("Loading %s for kernel\n", local_kernel_path)

	kernel_desc, err := os.OpenFile(local_kernel_path, os.O_RDONLY, 0)
	if err != nil {
		log.Fatalf("open(%q): %v", local_kernel_path, err)
	}
	defer kernel_desc.Close()

	var ramfs *os.File
	ramfs, err = os.OpenFile(local_initrd_path, os.O_RDONLY, 0)
	if err != nil {
		log.Fatalf("open(%q): %v", local_initrd_path, err)
	}
	defer ramfs.Close()

	if err := kexec.FileLoad(kernel_desc, ramfs, kernel_parameter); err != nil {
		log.Fatalf("%v", err)
	}
	if err := kexec.Reboot(); err != nil {
		log.Fatalf("%v", err)
	}

}

func usage() {
	log.Printf("Usage: %s [-v]\n", os.Args[1])
	os.Exit(2)
}

func main() {

	// Device list is in /sys/dev/block/
	var blk_list []string
	var supported_filesystem []string
	var parameter string
	verbose = false
	if len(os.Args) != 1 {
		if len(os.Args) != 2 {
			usage()
		}
		parameter = os.Args[1]
		if parameter == "-v" {
			verbose = true
			println("verbose mode activated")
		} else {
			usage()
		}
	}
	supported_filesystem = get_supported_filesystem()
	if verbose {
		println("************** Supported Filesystem by current linuxboot ********************")
		for _, filesystem := range supported_filesystem {
			println(filesystem)
		}
		println("*****************************************************************************")
	}
	blk_list = blk_devices_list("/sys/dev/block/", "/device/block/")
	// We must validate if the MBR is bootable or not and keep the
	// devices which do have such support
	// drive are easy to detect
	for _, entry := range blk_list {
		if check_for_bootable_mbr("/dev/"+entry) == 1 {
			fmt.Println("Bootable device found")
			// We need to loop on the device entries which are into /dev/<device>X
			// and mount each partitions as to find /boot entry if it is available somewhere
			var device_part_list []string
			device_part_list = get_device_part_list(entry)
			for _, device_list := range device_part_list {
				if mount_entry(device_list, supported_filesystem) {
					if verbose {
						println("mount succeed")
					}
					var grub_conf_path = check_boot_entry("/u-root/" + device_list)
					if grub_conf_path != "" {
						if verbose {
							println("calling basic kexec")
						}
						kexec_entry(grub_conf_path, "/u-root/"+device_list)
					}
				}
				umount_entry("/u-root/" + device_list)
			}
		}
	}
	println("Sorry no bootable device found")
}
