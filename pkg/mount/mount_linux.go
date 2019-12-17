// Copyright 2012-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package mount implements mounting, moving, and unmounting file systems.
package mount

import (
	"errors"
	"fmt"
	"os"

	"golang.org/x/sys/unix"
)

// Most commonly used mount flags.
const (
	MS_RDONLY   = unix.MS_RDONLY
	MS_BIND     = unix.MS_BIND
	MS_LAZYTIME = unix.MS_LAZYTIME
	MS_NOEXEC   = unix.MS_NOEXEC
	MS_NOSUID   = unix.MS_NOSUID
	MS_NOUSER   = unix.MS_NOUSER
	MS_RELATIME = unix.MS_RELATIME
	MS_SYNC     = unix.MS_SYNC
	MS_NOATIME  = unix.MS_NOATIME

	ReadOnly = unix.MS_RDONLY | unix.MS_NOATIME
)

// Unmount flags.
const (
	MNT_FORCE  = unix.MNT_FORCE
	MNT_DETACH = unix.MNT_DETACH
)

// Mounter is a device that can be attached at a file system path.
type Mounter interface {
	// Mount attaches the device at path.
	Mount(path string, flags uintptr) (*MountPoint, error)
}

// MountPoint represents a mounted file system.
type MountPoint struct {
	Path   string
	Device string
	FSType string
	Flags  uintptr
	Data   string
}

// String implements fmt.Stringer.
func (mp *MountPoint) String() string {
	return fmt.Sprintf("MountPoint(path=%s, device=%s, fs=%s, flags=%#x, data=%s)", mp.Path, mp.Device, mp.FSType, mp.Flags, mp.Data)
}

// Unmount unmounts a file system that was previously mounted.
func (mp *MountPoint) Unmount(flags uintptr) error {
	if err := unix.Unmount(mp.Path, int(flags)); err != nil {
		return &os.PathError{
			Op:   "unmount",
			Path: mp.Path,
			Err:  fmt.Errorf("flags %#x: %v", flags, err),
		}
	}
	return nil
}

// Mount attaches the fsType file system at path.
//
// dev is the device to mount (this is often the path of a block device, name
// of a file, or a dummy string). data usually contains arguments for the
// specific file system.
func Mount(dev, path, fsType, data string, flags uintptr) (*MountPoint, error) {
	// Create the mount point if it doesn't already exist.
	if err := os.MkdirAll(path, 0666); err != nil {
		return nil, err
	}

	if err := unix.Mount(dev, path, fsType, flags, data); err != nil {
		return nil, &os.PathError{
			Op:   "mount",
			Path: path,
			Err:  fmt.Errorf("from device %q (fs type %s, flags %#x): %v", dev, fsType, flags, err),
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

// TryMount tries to mount a device on the given mountpoint, trying in order
// the supported block device file systems on the system.
func TryMount(device, path string, flags uintptr) (*MountPoint, error) {
	// TryMount only works on existing block devices. No weirdo devices
	// like 9P.
	if _, err := os.Stat(device); err != nil {
		return nil, err
	}

	fs, err := GetBlockFilesystems()
	if err != nil {
		return nil, fmt.Errorf("failed to mount %s on %s: %v", device, path, err)
	}
	for _, fstype := range fs {
		mp, err := Mount(device, path, fstype, "", flags)
		if err != nil {
			continue
		}
		return mp, nil
	}
	return nil, fmt.Errorf("no suitable filesystem (out of %v) found to mount %s at %v", fs, device, path)
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
	var flags = unix.UMOUNT_NOFOLLOW
	if len(path) == 0 {
		return errors.New("path cannot be empty")
	}
	if force && lazy {
		return errors.New("MNT_FORCE and MNT_DETACH (lazy unmount) cannot both be set")
	}
	if force {
		flags |= unix.MNT_FORCE
	}
	if lazy {
		flags |= unix.MNT_DETACH
	}
	if err := unix.Unmount(path, flags); err != nil {
		return fmt.Errorf("umount %q flags %x: %v", path, flags, err)
	}
	return nil
}
