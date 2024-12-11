// Copyright 2012-2020 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package mount implements mounting, moving, and unmounting file systems.
package mount

import (
	"errors"
	"fmt"
	"os"
	"runtime"
	"syscall"
	"unsafe"

	"golang.org/x/sys/unix"
)

// Unmount flags.
const (
	MNT_FORCE = unix.MNT_FORCE
)

// MountPoint represents a mounted file system.
type MountPoint struct {
	Path   string
	Device string
	FSType string
	Flags  uintptr
	Data   string
}

// Unmount detaches any file system mounted at path.
//
// force forces an unmount regardless of currently open or otherwise used files
// within the file system to be unmounted.
//
// lazy disallows future uses of any files below path -- i.e. it hides the file
// system mounted at path, but the file system itself is still active and any
// currently open files can continue to be used. When all references to files
// from this file system are gone, the file system will actually be unmounted.
func Unmount(path string, force, lazy bool) error {
	if lazy {
		return fmt.Errorf("lazy: not on %v:%w", runtime.GOOS, os.ErrInvalid)
	}
	if len(path) == 0 {
		return errors.New("path cannot be empty")
	}
	var flags int
	if force {
		flags |= unix.MNT_FORCE
	}
	if err := unix.Unmount(path, flags); err != nil {
		return fmt.Errorf("umount %q flags %x: %w", path, flags, err)
	}
	return nil
}

// iov returns an iovec for a string.
// there is no official package, and it is simple
// enough, that we just create it here.
func iov(val string) syscall.Iovec {
	s := val + "\x00"
	vec := syscall.Iovec{Base: (*byte)(unsafe.Pointer(&[]byte(s)[0]))}
	vec.SetLen(len(s))
	return vec
}

// Mount attaches the fsType file system at path.
//
// dev is the device to mount (this is often the path of a block device, name
// of a file, or a placeholder string). data usually contains arguments for the
// specific file system.
//
// opts is usually empty, but if you want, e.g., to pre-create the mountpoint,
// you can call Mount with a mkdirall, e.g.
// mount.Mount("none", dst, fstype, "", 0,
//
//	func() error { return os.MkdirAll(dst, 0o666)})
func Mount(dev, path, fsType, data string, flags uintptr, opts ...func() error) (*MountPoint, error) {
	if len(path) == 0 {
		return nil, fmt.Errorf("len(path) can not be 0:%w", os.ErrInvalid)
	}

	if len(fsType) == 0 {
		return nil, fmt.Errorf("len(fstype) can not be 0:%w", os.ErrInvalid)
	}

	for _, f := range opts {
		if err := f(); err != nil {
			return nil, err
		}
	}

	// Create an array of iovec structures
	vec := []syscall.Iovec{
		iov("fstype"), iov(fsType),
		iov("fspath"), iov(path),
	}
	if len(dev) > 0 {
		vec = append(vec, iov("from"), iov(dev))
	}

	// Convert the slice of iovec to a pointer
	iovPtr := unsafe.Pointer(&vec[0])

	// Call nmount
	if _, _, errno := syscall.Syscall(syscall.SYS_NMOUNT, uintptr(iovPtr), uintptr(len(vec)), flags); errno != 0 {
		return nil, &os.PathError{
			Op:   "mount",
			Path: path,
			Err:  fmt.Errorf("from device %q (fs type %s, flags %#x): %w", dev, fsType, flags, errno),
		}
	}
	return &MountPoint{
		Path:   path,
		Device: dev,
		FSType: fsType,
		Data:   data,
		Flags:  flags,
	}, nil
}
