// Copyright 2018-2020 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package loop provides an interface to interacting with Linux loop devices.
//
// A loop device exposes a regular file as if it were a block device.
package loop

import (
	"github.com/u-root/u-root/pkg/mount"
)

// Loop represents a regular file exposed as a loop block device.
//
// Loop implements mount.Mounter.
type Loop struct {
	// Dev is the loop device path.
	Dev string

	// Source is the regular file to use as a block device.
	Source string

	// FSType is the file system to use when mounting the block device.
	FSType string

	// Data is the data to pass to mount(2).
	Data string
}

var _ mount.Mounter = &Loop{}

// New initializes a Loop struct and allocates a loop device to it.
//
// source is the file to use as a loop block device. fstype the file system
// name. data is the data argument to the mount(2) syscall.
func New(source, fstype string, data string) (*Loop, error) {
	devicename, err := FindDevice()
	if err != nil {
		return nil, err
	}
	if err := SetFile(devicename, source); err != nil {
		return nil, err
	}
	return &Loop{
		Dev:    devicename,
		Source: source,
		FSType: fstype,
		Data:   data,
	}, nil
}

// DevName implements mount.Mounter.
func (l *Loop) DevName() string {
	return l.Dev
}

// Mount mounts the provided source file, with type fstype, and flags and data options
// (which are usually 0 and ""), using the allocated loop device.
func (l *Loop) Mount(path string, flags uintptr, opts ...func() error) (*mount.MountPoint, error) {
	return mount.Mount(l.Dev, path, l.FSType, l.Data, flags, opts...)
}

// Free frees the loop device.
//
// All mount points must have been unmounted prior to calling this.
func (l *Loop) Free() error {
	return ClearFile(l.Dev)
}
