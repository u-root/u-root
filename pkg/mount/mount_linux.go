// Copyright 2012-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// The mount package implements functions for mounting and unmounting
// file systems and defines the mount interface.
package mount

import (
	"errors"
	"fmt"

	"golang.org/x/sys/unix"
)

func Mount(dev, path, fsType, data string, flags uintptr) error {
	if err := unix.Mount(dev, path, fsType, flags, data); err != nil {
		return fmt.Errorf("Mount %q on %q type %q flags %x: %v",
			dev, path, fsType, flags, err)
	}
	return nil
}

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
