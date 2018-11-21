// Copyright 2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package loop provides an interface to interacting with Linux loop devices.
//
// A loop device exposes a regular file as if it were a block device.
package loop

import (
	"github.com/u-root/u-root/pkg/mount"
	"golang.org/x/sys/unix"
)

// Loop represents a regular file exposed as a loop block device.
//
// Loop implements mount.Mount.
type Loop struct {
	// Dev is the loop device path.
	Dev string

	// Source is the regular file to use as a block device.
	Source string

	// Dir is the directory to mount the block device on.
	Dir string

	// FSType is the file system to use when mounting the block device.
	FSType string

	// Flags are flags to pass to mount(2).
	Flags uintptr

	// Data is the data to pass to mount(2).
	Data string

	// Mounted indicates whether the device has been mounted.
	Mounted bool
}

// New initializes a Loop struct and allocates a loodevice to it.
//
// source is the file to use as a loop block device. target is the directory
// the device should be mounted on.
func New(source, target, fstype string, flags uintptr, data string) (mount.Mounter, error) {
	devicename, err := FindDevice()
	if err != nil {
		return nil, err
	}
	if err := SetFile(devicename, source); err != nil {
		return nil, err
	}
	return &Loop{
		Dev:    devicename,
		Dir:    target,
		Source: source,
		FSType: fstype,
		Flags:  flags,
		Data:   data,
	}, nil
}

// Mount mounts the provided source file, with type fstype, and flags and data options
// (which are usually 0 and ""), using the allocated loop device.
func (l *Loop) Mount() error {
	if err := unix.Mount(l.Dev, l.Dir, l.FSType, l.Flags, l.Data); err != nil {
		return err
	}
	l.Mounted = true
	return nil
}
