// Copyright 2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package securelaunch takes integrity measurements before launching the target system.
package securelaunch

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/u-root/u-root/pkg/storage"
)

/* used to store all block devices returned from a call to storage.GetBlockStats */
var StorageBlkDevices []storage.BlockDev

/*
 * if kernel cmd line has uroot.uinitargs=-d, debug fn is enabled.
 * kernel cmdline is checked in sluinit.
 */
var Debug = func(string, ...interface{}) {}

/*
 * WriteToFile writes a byte slice to a target file on an
 * already mounted disk and returns the target file path.
 *
 * defFileName is default dst file name, only used if user doesn't provide one.
 */
func WriteToFile(data []byte, dst, defFileName string) (string, error) {

	// make sure dst is an absolute file path
	if !filepath.IsAbs(dst) {
		return "", fmt.Errorf("dst =%s Not an absolute path ", dst)
	}

	// target is the full absolute path where []byte will be written to
	target := dst
	dstInfo, err := os.Stat(dst)
	if err == nil && dstInfo.IsDir() {
		Debug("No file name provided. Adding it now. old target=%s", target)
		target = filepath.Join(dst, defFileName)
		Debug("New target=%s", target)
	}

	Debug("target=%s", target)
	err = ioutil.WriteFile(target, data, 0644)
	if err != nil {
		return "", fmt.Errorf("failed to write date to file =%s, err=%v", target, err)
	}
	Debug("WriteToFile exit w success data written to target=%s", target)
	return target, nil
}

func getDeviceFromUUID(uuid string) (storage.BlockDev, error) {
	if e := GetBlkInfo(); e != nil {
		return storage.BlockDev{}, fmt.Errorf("GetBlkInfo err=%s", e)
	}
	devices := storage.PartitionsByFsUUID(StorageBlkDevices, uuid) // []BlockDev
	Debug("%d device(s) matched with UUID=%s", len(devices), uuid)
	for i, d := range devices {
		Debug("No#%d ,device=%s with fsUUID=%s", i, d.Name, d.FsUUID)
		return d, nil // return first device found
	}
	return storage.BlockDev{}, fmt.Errorf("no block device exists with UUID=%s", uuid)
}

func getDeviceFromName(name string) (storage.BlockDev, error) {
	if e := GetBlkInfo(); e != nil {
		return storage.BlockDev{}, fmt.Errorf("GetBlkInfo err=%s", e)
	}
	devices := storage.PartitionsByName(StorageBlkDevices, name) // []BlockDev
	Debug("%d device(s) matched with Name=%s", len(devices), name)
	for i, d := range devices {
		Debug("No#%d ,device=%s with fsUUID=%s", i, d.Name, d.FsUUID)
		return d, nil // return first device found
	}
	return storage.BlockDev{}, fmt.Errorf("no block device exists with name=%s", name)
}

func GetStorageDevice(input string) (storage.BlockDev, error) {
	device, e := getDeviceFromUUID(input)
	if e != nil {
		d2, e2 := getDeviceFromName(input)
		if e2 != nil {
			return storage.BlockDev{}, fmt.Errorf("getDeviceFromUUID: err=%v, getDeviceFromName: err=%v", e, e2)
		}
		device = d2
	}
	return device, nil
}

/*
 * GetMountedFilePath returns a file path corresponding to a <device_identifier>:<path> user input format.
 * <device_identifier> may be a Linux block device identifier like sda or a FS UUID.
 *
 * NOTE: Caller's responsbility to unmount this..use return var mountPath to unmount in caller.
 */
func GetMountedFilePath(inputVal string, flags uintptr) (string, string, error) {
	s := strings.Split(inputVal, ":")
	if len(s) != 2 {
		return "", "", fmt.Errorf("%s: Usage: <block device identifier>:<path>", inputVal)
	}

	// s[0] can be sda or UUID.
	device, err := GetStorageDevice(s[0])
	if err != nil {
		return "", "", fmt.Errorf("fn GetStorageDevice: err = %v", err)
	}

	devPath := filepath.Join("/dev", device.Name)
	Debug("Attempting to mount %s", devPath)
	mountPath, err := ioutil.TempDir("/tmp", "slaunch-")
	if err != nil {
		return "", "", fmt.Errorf("failed to create tmp mount directory: %v", err)
	}

	if _, err := device.Mount(mountPath, flags); err != nil {
		return "", "", fmt.Errorf("failed to mount %s , flags=%v, err=%v", devPath, flags, err)
	}

	Debug("Mounted %s", devPath)
	fPath := filepath.Join(mountPath, s[1]) // mountPath=/tmp/path/to/target/file if /dev/sda mounted on /tmp
	return fPath, mountPath, nil
}

/*
 * GetBlkInfo calls storage package to get information on all block devices.
 * The information is stored in a global variable 'StorageBlkDevices'
 * If the global variable is already non-zero, we skip the call to storage package.
 *
 * In debug mode, it also prints names and UUIDs for all devices.
 */
func GetBlkInfo() error {
	if len(StorageBlkDevices) == 0 {
		var err error
		Debug("getBlkInfo: expensive function call to get block stats from storage pkg")
		StorageBlkDevices, err = storage.GetBlockStats()
		if err != nil {
			return fmt.Errorf("getBlkInfo: storage.GetBlockStats err=%v. Exiting", err)
		}
		// no block devices exist on the system.
		if len(StorageBlkDevices) == 0 {
			return fmt.Errorf("getBlkInfo: no block devices found")
		}
		// print the debug info only when expensive call to storage is made
		for k, d := range StorageBlkDevices {
			Debug("block device #%d, Name=%s, FSType=%s, FsUUID=%s", k, d.Name, d.FSType, d.FsUUID)
		}
		return nil
	}
	Debug("getBlkInfo: noop")
	return nil
}
