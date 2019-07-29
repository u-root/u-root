// Copyright 2012-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mount

import (
	"errors"
	"fmt"

	"golang.org/x/sys/unix"
)

// Mount attaches the fsType file system at path.
//
// dev is the device to mount (this is often the path of a block device, name
// of a file, or a dummy string). data usually contains arguments for the
// specific file system.
func Mount(dev, path, fsType, data string, flags uintptr) error {
	if err := unix.Mount(dev, path, fsType, flags, data); err != nil {
		return fmt.Errorf("Mount %q on %q type %q flags %x: %v",
			dev, path, fsType, flags, err)
	}
	return nil
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
		return errors.New("force and lazy unmount cannot both be set")
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
