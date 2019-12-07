// Copyright 2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package securelaunch takes integrity measurements before launching the target system.
package securelaunch

import (
	"fmt"
	"log"
	"path/filepath"
	"strings"

	"github.com/u-root/u-root/pkg/boot/diskboot"
	"github.com/u-root/u-root/pkg/storage"
)

/* used to store all block devices returned from a call to storage.GetBlockStats */
var storageBlkDevices []storage.BlockDev

/*
 * if kernel cmd line has uroot.uinitargs=-d, debug fn is enabled.
 * kernel cmdline is checked in sluinit.
 */
var Debug = func(string, ...interface{}) {}

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
		if e := getBlkInfo(); e != nil {
			return "", "", fmt.Errorf("getBlkInfo err=%s", e)
		}
		devices := storage.PartitionsByFsUUID(storageBlkDevices, s[0]) // []BlockDev
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
 * getBlkInfo calls storage package to get information on all block devices.
 * The information is stored in a global variable 'storageBlkDevices'
 * If the global variable is already non-zero, we skip the call to storage package.
 *
 * In debug mode, it also prints names and UUIDs for all devices.
 */
func getBlkInfo() error {
	if len(storageBlkDevices) == 0 {
		var err error
		storageBlkDevices, err = storage.GetBlockStats()
		if err != nil {
			log.Printf("getBlkInfo: storage.GetBlockStats err=%v. Exiting", err)
			return err
		}
	}

	for k, d := range storageBlkDevices {
		Debug("block device #%d, Name=%s, FsUUID=%s", k, d.Name, d.FsUUID)
	}
	return nil
}
