// Copyright 2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package securelaunch takes integrity measurements before launching the target system.
package securelaunch

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/u-root/u-root/pkg/mount"
	"github.com/u-root/u-root/pkg/mount/block"
)

type persistDataItem struct {
	desc        string // Description
	data        []byte
	location    string // of form sda:/path/to/file
	defaultFile string // if location turns out to be dir only
}

var persistData []persistDataItem

type mountCacheData struct {
	flags     uintptr
	mountPath string
}

type mountCacheType struct {
	m  map[string]mountCacheData
	mu sync.RWMutex
}

// ErrUsage indicates a usage error.
var ErrUsage = errors.New("incorrect usage")

// mountCache is used by sluinit to reduce number of mount/unmount operations
var mountCache = mountCacheType{m: make(map[string]mountCacheData)}

// StorageBlkDevices helps securelaunch pkg mount devices.
var StorageBlkDevices block.BlockDevices

// Debug enables verbose logs if kernel cmd line has uroot.uinitargs=-d flag set.
// kernel cmdline is checked in sluinit.
var Debug = func(string, ...interface{}) {}

// ReadFile reads a file into a byte slice. It mounts the disk if necessary.
//
// policyLocation is formatted as `<block device id>:<path>`
//
//	e.g., sda1:/boot/securelaunch.policy
//	e.g., 4qccd342-12zr-4e99-9ze7-1234cb1234c4:/foo.txt
func ReadFile(fileLocation string) ([]byte, error) {
	mountedFilePath, err := GetMountedFilePath(fileLocation, mount.MS_RDONLY)
	if err != nil {
		return nil, fmt.Errorf("writing to %q:%w", fileLocation, err)
	}

	Debug("ReadFile: reading %q (mounted at %q):%w", fileLocation, mountedFilePath, err)

	fileBytes, err := os.ReadFile(mountedFilePath)
	if err != nil {
		return nil, fmt.Errorf("%q (mounted at %q):%w", fileLocation, mountedFilePath, err)
	}

	return fileBytes, nil
}

func WriteFile(data []byte, fileLocation string) error {
	mountedFilePath, err := GetMountedFilePath(fileLocation, 0) // 0 means RW
	if err != nil {
		return err
	}

	Debug("WriteFile: writing file %q (mounted at %q)", fileLocation, mountedFilePath)

	err = os.WriteFile(mountedFilePath, data, 0o644)
	if err != nil {
		return fmt.Errorf("failed to write file %q: %w", mountedFilePath, err)
	}

	return nil
}

// persist writes data to the given targetPath. If targetPath is a directory,
// then defaultFileName is used as the actual file name.
//
// targetPath is formatted as `<block device id>:<path>`
func persist(data []byte, targetPath string, defaultFileName string) error {
	mountedFilePath, err := GetMountedFilePath(targetPath, 0) // 0 is flag for rw mount option
	if err != nil {
		return fmt.Errorf("failed to locate file %q for writing: %w", targetPath, err)
	}

	// Check if a file name was provided or just the path. If just the path,
	// add the provided default file name.
	mountedFilePathInfo, err := os.Stat(mountedFilePath)
	if err != nil {
		return fmt.Errorf("failed to stat file %q (mounted at %q): %w", targetPath, mountedFilePath, err)
	}
	if mountedFilePathInfo.IsDir() {
		Debug("persist: No file name provided, adding default name %q", defaultFileName)
		targetPath = filepath.Join(targetPath, defaultFileName)
		Debug("persist: New file path: %q", targetPath)
	}

	// Write the file.
	if err := WriteFile(data, targetPath); err != nil {
		return err
	}

	return nil
}

// AddToPersistQueue enqueues an action item to persistData slice
// so that it can be deferred to the last step of sluinit.
func AddToPersistQueue(desc string, data []byte, location string, defFile string) error {
	persistData = append(persistData, persistDataItem{desc, data, location, defFile})
	return nil
}

// ClearPersistQueue persists any pending data/logs to disk
func ClearPersistQueue() error {
	for _, entry := range persistData {
		if err := persist(entry.data, entry.location, entry.defaultFile); err != nil {
			return fmt.Errorf("%s: persist failed for location %s", entry.desc, entry.location)
		}
	}
	return nil
}

func getDeviceFromUUID(uuid string) (*block.BlockDev, error) {
	if e := GetBlkInfo(); e != nil {
		return nil, fmt.Errorf("fn GetBlkInfo err=%w", e)
	}
	devices := StorageBlkDevices.FilterFSUUID(uuid)
	Debug("%d device(s) matched with UUID=%s", len(devices), uuid)
	for i, d := range devices {
		Debug("No#%d ,device=%s with fsUUID=%s", i, d.Name, d.FsUUID)
		return d, nil // return first device found
	}
	return nil, fmt.Errorf("no block device exists with UUID=%s", uuid)
}

func getDeviceFromName(name string) (*block.BlockDev, error) {
	if e := GetBlkInfo(); e != nil {
		return nil, fmt.Errorf("fn GetBlkInfo err=%w", e)
	}
	devices := StorageBlkDevices.FilterName(name)
	Debug("%d device(s) matched with Name=%s", len(devices), name)
	for i, d := range devices {
		Debug("No#%d ,device=%s with fsUUID=%s", i, d.Name, d.FsUUID)
		return d, nil // return first device found
	}
	return nil, fmt.Errorf("no block device exists with name=%s", name)
}

// GetStorageDevice parses input of type UUID:/tmp/foo or sda2:/tmp/foo,
// and returns any matching devices.
func GetStorageDevice(input string) (*block.BlockDev, error) {
	device, e := getDeviceFromUUID(input)
	if e != nil {
		d2, e2 := getDeviceFromName(input)
		if e2 != nil {
			return nil, fmt.Errorf("getDeviceFromUUID: err=%w, getDeviceFromName: err=%w", e, e2)
		}
		device = d2
	}
	return device, nil
}

func deleteEntryMountCache(key string) {
	mountCache.mu.Lock()
	delete(mountCache.m, key)
	mountCache.mu.Unlock()

	Debug("mountCache: Deleted key %s", key)
}

func setMountCache(key string, val mountCacheData) {
	mountCache.mu.Lock()
	mountCache.m[key] = val
	mountCache.mu.Unlock()

	Debug("mountCache: Updated key %s, value %v", key, val)
}

// getMountCacheData looks up mountCache using devName as key
// and clears an entry in cache if result is found with different
// flags, otherwise returns the cached entry or nil.
func getMountCacheData(key string, flags uintptr) (string, error) {
	Debug("mountCache: Lookup with key %s", key)
	cachedData, ok := mountCache.m[key]
	if ok {
		cachedMountPath := cachedData.mountPath
		cachedFlags := cachedData.flags
		Debug("mountCache: Lookup succeeded: cachedMountPath %s, cachedFlags %d found for key %s", cachedMountPath, cachedFlags, key)
		if cachedFlags == flags {
			return cachedMountPath, nil
		}
		Debug("mountCache: need to mount the same device with different flags")
		Debug("mountCache: Unmounting %s first", cachedMountPath)
		if err := mount.Unmount(cachedMountPath, true, false); err != nil {
			return "", fmt.Errorf("failed to unmount %q: %w", cachedMountPath, err)
		}
		Debug("mountCache: unmount successfull. lets delete entry in map")
		deleteEntryMountCache(key)
		return "", fmt.Errorf("device was already mounted: mount again")
	}

	return "", fmt.Errorf("mountCache: lookup failed, no key exists that matches %s", key)
}

// MountDevice looks up mountCache map. if no entry is found, it
// mounts a device and updates cache, otherwise returns mountPath.
func MountDevice(device *block.BlockDev, flags uintptr) (string, error) {
	devName := device.Name

	Debug("MountDevice: Checking cache first for %s", devName)
	cachedMountPath, err := getMountCacheData(devName, flags)
	if err == nil {
		return cachedMountPath, nil
	}
	Debug("MountDevice: cache lookup failed for %q", devName)

	Debug("MountDevice: Attempting to mount %q with flags %#x", devName, flags)
	mountPath, err := os.MkdirTemp("/tmp", "slaunch-")
	if err != nil {
		return "", fmt.Errorf("create tmp mount directory: %w", err)
	}

	if _, err := device.Mount(mountPath, flags); err != nil {
		return "", fmt.Errorf("mount %q, flags %#x:%w", devName, flags, err)
	}

	Debug("MountDevice: Mounted %q with flags %#x", devName, flags)
	setMountCache(devName, mountCacheData{flags: flags, mountPath: mountPath}) // update cache
	return mountPath, nil
}

// GetMountedFilePath returns the file path corresponding to the given
// <device_identifier>:<path>.
// <device_identifier> is a Linux block device identifier (e.g, sda or UUID).
func GetMountedFilePath(inputVal string, flags uintptr) (string, error) {
	s := strings.Split(inputVal, ":")
	if len(s) != 2 {
		return "", fmt.Errorf("%s: Usage: <block device identifier>:<path>", inputVal)
	}

	// s[0] can be sda or UUID.
	device, err := GetStorageDevice(s[0])
	if err != nil {
		return "", fmt.Errorf("GetStorageDevice:%w", err)
	}

	devName := device.Name
	mountPath, err := MountDevice(device, flags)
	if err != nil {
		return "", fmt.Errorf("failed to mount %s , flags=%v, err=%w", devName, flags, err)
	}

	fPath := filepath.Join(mountPath, s[1])
	return fPath, nil
}

// UnmountAll unmounts all mounted devices from the file heirarchy.
func UnmountAll() error {
	Debug("UnmountAll: %d devices need to be unmounted", len(mountCache.m))
	for key, mountCacheData := range mountCache.m {
		cachedMountPath := mountCacheData.mountPath
		Debug("UnmountAll: Unmounting %s", cachedMountPath)
		if err := mount.Unmount(cachedMountPath, true, false); err != nil {
			return fmt.Errorf("failed to unmount %q: %w", cachedMountPath, err)
		}
		Debug("UnmountAll: Unmounted %s", cachedMountPath)
		deleteEntryMountCache(key)
		Debug("UnmountAll: Deleted key %s from cache", key)
	}

	return nil
}

// GetBlkInfo gets information on all block devices and stores it in the
// global variable 'StorageBlkDevices'. If it is called more than once, the
// subsequent calls just return.
//
// In debug mode, it also prints names and UUIDs for all devices.
func GetBlkInfo() error {
	if len(StorageBlkDevices) == 0 {
		var err error
		Debug("getBlkInfo: expensive function call to get block stats from storage pkg")
		StorageBlkDevices, err = block.GetBlockDevices()
		if err != nil {
			return fmt.Errorf("getBlkInfo: storage.GetBlockDevices err=%w. Exiting", err)
		}
		// no block devices exist on the system.
		if len(StorageBlkDevices) == 0 {
			return fmt.Errorf("getBlkInfo: no block devices found")
		}
		// print the debug info only when expensive call to storage is made
		for k, d := range StorageBlkDevices {
			Debug("block device #%d: %s", k, d)
		}
		return nil
	}
	Debug("getBlkInfo: noop")
	return nil
}

// GetFileBytes reads the given file and returns the contents as a byte slice.
func GetFileBytes(fileName string) ([]byte, error) {
	filePath, err := GetMountedFilePath(fileName, mount.MS_RDONLY)
	if err != nil {
		return nil, fmt.Errorf("could not get mounted file path %q: %w", fileName, err)
	}
	Debug("GetFileBytes: file path = %q", filePath)

	fileBytes, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("could not read file %q: %w", filePath, err)
	}

	return fileBytes, nil
}
