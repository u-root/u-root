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

	// Mounted indicates whether the device has been mounted.
	Mounted bool

	// dir is the directory the block device was mounted on.
	dir string
}

// New initializes a Loop struct and allocates a loop device to it.
//
// source is the file to use as a loop block device. fstype the file system
// name. data is the data argument to the mount(2) syscall.
func New(source, fstype string, data string) (mount.Mounter, error) {
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

// Mount mounts the provided source file, with type fstype, and flags and data options
// (which are usually 0 and ""), using the allocated loop device.
func (l *Loop) Mount(path string, flags uintptr) error {
	l.dir = path
	if err := unix.Mount(l.Dev, path, l.FSType, flags, l.Data); err != nil {
		return err
	}
	l.Mounted = true
	return nil
}

const forceUnmount = unix.MNT_FORCE | unix.MNT_DETACH

// Unmount unmounts and frees a loop. If it is mounted, it will try to unmount it.
// If the unmount fails, we try to free it anyway, after trying a more
// forceful unmount.
func (l *Loop) Unmount(flags int) error {
	if l.Mounted {
		if err := unix.Unmount(l.dir, flags); err != nil {
			unix.Unmount(l.dir, flags|forceUnmount)
		}
	}
	return ClearFile(l.Dev)
}
