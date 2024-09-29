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

	"golang.org/x/sys/unix"
)

// Unmount flags.
const (
	MNT_FORCE = unix.MNT_FORCE
)

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
