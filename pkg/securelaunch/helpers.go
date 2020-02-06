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

	"github.com/u-root/u-root/pkg/boot/diskboot"
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

	// s[0] can be sda or UUID. if UUID, then we need to find its name
	deviceId := s[0]
	if !strings.HasPrefix(deviceId, "sd") {
		if e := GetBlkInfo(); e != nil {
			return "", "", fmt.Errorf("GetBlkInfo err=%s", e)
		}
		devices := storage.PartitionsByFsUUID(StorageBlkDevices, s[0]) // []BlockDev
		for _, device := range devices {
			Debug("device =%s with fsuuid=%s", device.Name, s[0])
			deviceId = device.Name
		}
	}

	devicePath := filepath.Join("/dev", deviceId) // assumes deviceId is sda, devicePath=/dev/sda
	Debug("Attempting to mount %s", devicePath)
	dev, err := diskboot.FindDevice(devicePath, flags) // FindDevice fn mounts devicePath=/dev/sda.
	if err != nil {
		return "", "", fmt.Errorf("failed to mount %v , flags=%v, err=%v", devicePath, flags, err)
	}

	Debug("Mounted %s", devicePath)
	fPath := filepath.Join(dev.MountPoint.Path, s[1]) // mountPath=/tmp/path/to/target/file if /dev/sda mounted on /tmp
	return fPath, dev.MountPoint.Path, nil
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
