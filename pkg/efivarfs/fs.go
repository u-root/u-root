// Copyright 2022 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package efivarfs allows interaction with efivarfs of the
// linux kernel.
package efivarfs

import (
	"os"

	"golang.org/x/sys/unix"
)

// FS_IMMUTABLE_FL is an inode flag to make a file immutable.
const FS_IMMUTABLE_FL = 0x10

// getInodeFlags returns the extended attributes of a file.
func getInodeFlags(f *os.File) (int, error) {
	// If I knew how unix.Getxattr works I'd use that...
	flags, err := unix.IoctlGetInt(int(f.Fd()), unix.FS_IOC_GETFLAGS)
	if err != nil {
		return 0, &os.PathError{Op: "ioctl(FS_IOC_GETFLAGS)", Path: f.Name(), Err: err}
	}
	return flags, nil
}

// setInodeFlags sets the extended attributes of a file.
func setInodeFlags(f *os.File, flags int) error {
	// If I knew how unix.Setxattr works I'd use that...
	if err := unix.IoctlSetPointerInt(int(f.Fd()), unix.FS_IOC_SETFLAGS, flags); err != nil {
		return &os.PathError{Op: "ioctl(FS_IOC_SETFLAGS)", Path: f.Name(), Err: err}
	}
	return nil
}

// makeMutable will change a files xattrs so that
// the immutable flag is removed and return a restore
// function which can reset the flag for that filee.
func makeMutable(f *os.File) (restore func(), err error) {
	flags, err := getInodeFlags(f)
	if err != nil {
		return nil, err
	}
	if flags&FS_IMMUTABLE_FL == 0 {
		return func() {}, nil
	}

	if err := setInodeFlags(f, flags&^FS_IMMUTABLE_FL); err != nil {
		return nil, err
	}
	return func() {
		if err := setInodeFlags(f, flags); err != nil {
			// If setting the immutable did
			// not work it's alright to do nothing
			// because after a reboot the flag is
			// automatically reapplied
			return
		}
	}, nil
}
