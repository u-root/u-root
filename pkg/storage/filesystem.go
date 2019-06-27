// Copyright 2017-2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package storage

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"syscall"
)

// Mountpoint holds mount point information for a given device
type Mountpoint struct {
	DeviceName string
	Path       string
	FsType     string
}

// GetSupportedFilesystems returns the supported file systems for block devices
func GetSupportedFilesystems() (fstypes []string, err error) {
	return internalGetFilesystems("/proc/filesystems")
}

func internalGetFilesystems(file string) (fstypes []string, err error) {
	var bytes []byte
	if bytes, err = ioutil.ReadFile(file); err != nil {
		return nil, fmt.Errorf("Failed to read %s: %v", file, err)
	}
	for _, line := range strings.Split(string(bytes), "\n") {
		//len(fields)==1, 2 possibilites for fs: "nodev" fs and
		// fs's. "nodev" fs cannot be mounted through devices.
		// len(fields)==1 prevents this from occurring.
		if fields := strings.Fields(line); len(fields) == 1 {
			fstypes = append(fstypes, fields[0])
		}
	}
	return fstypes, nil
}

// Mount tries to mount a block device on the given mountpoint, trying in order
// the provided file system types. It returns a Mountpoint structure, or an error
// if the device could not be mounted. If the mount point does not exist, it will
// be created.
func Mount(devname, mountpath string, filesystems []string) (*Mountpoint, error) {
	if err := os.MkdirAll(mountpath, 0744); err != nil {
		return nil, err
	}
	for _, fstype := range filesystems {
		log.Printf(" * trying %s on %s", fstype, devname)
		// MS_RDONLY should be enough. See mount(2)
		flags := uintptr(syscall.MS_RDONLY)
		// no options
		data := ""
		if err := syscall.Mount(devname, mountpath, fstype, flags, data); err != nil {
			log.Printf("    failed with %v", err)
			continue
		}
		log.Printf(" * mounted %s on %s with filesystem type %s", devname, mountpath, fstype)
		return &Mountpoint{DeviceName: devname, Path: mountpath, FsType: fstype}, nil
	}
	return nil, fmt.Errorf("no suitable filesystem type found to mount %s", devname)
}
