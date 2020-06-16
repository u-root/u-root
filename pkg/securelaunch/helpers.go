// Copyright 2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package securelaunch takes integrity measurements before launching the target system.
package securelaunch

import (
	"fmt"
	"io/ioutil"
	"log"
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

// mountCache is used by sluinit to reduce number of mount/unmount operations
var mountCache = mountCacheType{m: make(map[string]mountCacheData)}

// StorageBlkDevices helps securelaunch pkg mount devices.
var StorageBlkDevices block.BlockDevices

// Debug enables verbose logs if kernel cmd line has uroot.uinitargs=-d flag set.
// kernel cmdline is checked in sluinit.
var Debug = func(string, ...interface{}) {}

// WriteToFile writes a byte slice to a target file on an
// already mounted disk and returns the target file path.
//
// defFileName is default dst file name, only used if user doesn't provide one.
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

	Debug("WriteToFile: target=%s", target)
	err = ioutil.WriteFile(target, data, 0644)
	if err != nil {
		return "", fmt.Errorf("failed to write date to file =%s, err=%v", target, err)
	}
	Debug("WriteToFile: exit w success data written to target=%s", target)
	return target, nil
}

// persist writes data to targetPath.
// targetPath is of form sda:/boot/cpuid.txt
func persist(data []byte, targetPath string, defaultFile string) error {

	filePath, r := GetMountedFilePath(targetPath, 0) // 0 is flag for rw mount option
	if r != nil {
		return fmt.Errorf("persist: err: input %s could NOT be located, err=%v", targetPath, r)
	}

	dst := filePath // /tmp/boot-733276578/cpuid

	target, err := WriteToFile(data, dst, defaultFile)
	if err != nil {
		log.Printf("persist: err=%s", err)
		return err
	}

	Debug("persist: Target File%s", target)
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
		return nil, fmt.Errorf("fn GetBlkInfo err=%s", e)
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
		return nil, fmt.Errorf("fn GetBlkInfo err=%s", e)
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
			return nil, fmt.Errorf("getDeviceFromUUID: err=%v, getDeviceFromName: err=%v", e, e2)
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
		if e := mount.Unmount(cachedMountPath, true, false); e != nil {
			log.Printf("Unmount failed for %s. PANIC", cachedMountPath)
			panic(e)
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
		log.Printf("getMountCacheData succeeded for %s", devName)
		return cachedMountPath, nil
	}
	Debug("MountDevice: cache lookup failed for %s", devName)

	Debug("MountDevice: Attempting to mount %s with flags %d", devName, flags)
	mountPath, err := ioutil.TempDir("/tmp", "slaunch-")
	if err != nil {
		return "", fmt.Errorf("failed to create tmp mount directory: %v", err)
	}

	if _, err := device.Mount(mountPath, flags); err != nil {
		return "", fmt.Errorf("failed to mount %s, flags %d, err=%v", devName, flags, err)
	}

	Debug("MountDevice: Mounted %s with flags %d", devName, flags)
	setMountCache(devName, mountCacheData{flags: flags, mountPath: mountPath}) // update cache
	return mountPath, nil
}

// GetMountedFilePath returns a file path corresponding to a <device_identifier>:<path> user input format.
// <device_identifier> may be a Linux block device identifier like sda or a FS UUID.
func GetMountedFilePath(inputVal string, flags uintptr) (string, error) {
	s := strings.Split(inputVal, ":")
	if len(s) != 2 {
		return "", fmt.Errorf("%s: Usage: <block device identifier>:<path>", inputVal)
	}

	// s[0] can be sda or UUID.
	device, err := GetStorageDevice(s[0])
	if err != nil {
		return "", fmt.Errorf("fn GetStorageDevice: err = %v", err)
	}

	devName := device.Name
	mountPath, err := MountDevice(device, flags)
	if err != nil {
		return "", fmt.Errorf("failed to mount %s , flags=%v, err=%v", devName, flags, err)
	}

	fPath := filepath.Join(mountPath, s[1]) // mountPath=/tmp/path/to/target/file if /dev/sda mounted on /tmp
	return fPath, nil
}

// UnmountAll loops detaches any mounted device from the file heirarchy.
func UnmountAll() {
	Debug("UnmountAll: %d devices need to be unmounted", len(mountCache.m))
	for key, mountCacheData := range mountCache.m {
		cachedMountPath := mountCacheData.mountPath
		Debug("UnmountAll: Unmounting %s", cachedMountPath)
		if e := mount.Unmount(cachedMountPath, true, false); e != nil {
			log.Printf("Unmount failed for %s. PANIC", cachedMountPath)
			panic(e)
		}
		Debug("UnmountAll: Unmounted %s", cachedMountPath)
		deleteEntryMountCache(key)
		Debug("UnmountAll: Deleted key %s from cache", key)
	}
}

// GetBlkInfo calls storage package to get information on all block devices.
// The information is stored in a global variable 'StorageBlkDevices'
// If the global variable is already non-zero, we skip the call to storage package.
//
// In debug mode, it also prints names and UUIDs for all devices.
func GetBlkInfo() error {
	if len(StorageBlkDevices) == 0 {
		var err error
		Debug("getBlkInfo: expensive function call to get block stats from storage pkg")
		StorageBlkDevices, err = block.GetBlockDevices()
		if err != nil {
			return fmt.Errorf("getBlkInfo: storage.GetBlockDevices err=%v. Exiting", err)
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
