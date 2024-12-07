// Copyright 2012-2020 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package mount implements mounting, moving, and unmounting file systems.
package mount

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

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
	// DevName returns the name of the device.
	DevName() string
	// Mount attaches the device at path.
	Mount(path string, flags uintptr, opts ...func() error) (*MountPoint, error)
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
			Err:  fmt.Errorf("flags %#x: %w", flags, err),
		}
	}
	return nil
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
	for _, f := range opts {
		if err := f(); err != nil {
			return nil, err
		}
	}

	if err := unix.Mount(dev, path, fsType, flags, data); err != nil {
		return nil, &os.PathError{
			Op:   "mount",
			Path: path,
			Err:  fmt.Errorf("from device %q (fs type %s, flags %#x): %w", dev, fsType, flags, err),
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
func TryMount(device, path, data string, flags uintptr, opts ...func() error) (*MountPoint, error) {
	fstype, extraflags, err := FSFromBlock(device)
	if err != nil {
		// try statfs
		var statErr error
		fstype, extraflags, statErr = FromStatFS(device)
		if statErr == nil {
			return Mount(device, path, fstype, data, flags|extraflags, opts...)
		}

		return nil, err
	}

	return Mount(device, path, fstype, data, flags|extraflags, opts...)
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
	flags := unix.UMOUNT_NOFOLLOW
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
		return fmt.Errorf("umount %q flags %x: %w", path, flags, err)
	}
	return nil
}

// Pool keeps track of multiple MountPoint.
type Pool struct {
	// List of items mounted by this pool.
	MountPoints []*MountPoint
	// Temporary directory which contains sub-directories for mounts.
	tmpDir string
}

// Mount mounts a file system using Mounter and returns the MountPoint. If the
// device has already been mounted, it is not mounted again.
//
// Note the pool is keyed on Mounter.DevName() alone meaning DevName is used to
// determine whether it has already been mounted.
func (p *Pool) Mount(mounter Mounter, flags uintptr) (*MountPoint, error) {
	for _, m := range p.MountPoints {
		if m.Device == filepath.Join("/dev", mounter.DevName()) {
			return m, nil
		}
	}

	// Create temporary directory if one does not already exist.
	if p.tmpDir == "" {
		tmpDir, err := os.MkdirTemp("", "u-root-mounts")
		if err != nil {
			return nil, fmt.Errorf("cannot create tmpdir: %w", err)
		}
		p.tmpDir = tmpDir
	}

	path := filepath.Join(p.tmpDir, mounter.DevName())
	os.MkdirAll(path, 0o777)
	m, err := mounter.Mount(path, flags)
	if err != nil {
		// unix.Rmdir is used (instead of os.RemoveAll) because it
		// fails when the directory is non-empty. It would be a bit
		// dangerous to use os.RemoveAll because it could accidentally
		// delete everything in a mount.
		unix.Rmdir(path)
		return nil, err
	}
	p.MountPoints = append(p.MountPoints, m)
	return m, err
}

// Add adds MountPoints to the pool.
func (p *Pool) Add(m ...*MountPoint) {
	p.MountPoints = append(p.MountPoints, m...)
}

// UnmountAll umounts all the mountpoints from the pool. This makes a
// best-effort attempt to unmount everything and cleanup temporary directories.
// If this function fails, it can be re-tried.
func (p *Pool) UnmountAll(flags uintptr) error {
	// Errors get concatenated together here.
	var returnErr error

	for _, m := range p.MountPoints {
		if err := m.Unmount(flags); err != nil {
			if returnErr == nil {
				returnErr = fmt.Errorf("(Unmount) %s", err.Error())
			} else {
				returnErr = fmt.Errorf("%w; (Unmount) %s", returnErr, err.Error())
			}
		}

		// unix.Rmdir is used (instead of os.RemoveAll) because it
		// fails when the directory is non-empty. It would be a bit
		// dangerous to use os.RemoveAll because it could accidentally
		// delete everything in a mount.
		unix.Rmdir(m.Path)
	}

	if returnErr == nil && p.tmpDir != "" {
		if err := unix.Rmdir(p.tmpDir); err != nil {
			if returnErr == nil {
				returnErr = fmt.Errorf("(Rmdir) %s", err.Error())
			} else {
				returnErr = fmt.Errorf("%w; (Rmdir) %s", returnErr, err.Error())
			}
		}
		p.tmpDir = ""
	}

	return returnErr
}
