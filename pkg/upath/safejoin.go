// Copyright 2020 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package upath

import (
	"fmt"
	"path/filepath"
)

// SafeFilepathJoin safely joins two paths path1+path2. The resulting path will
// always be contained within path1 even if path2 tries to escape with "../".
//
// Note: This only works on Unix-like file systems.
// Second note: Additional care has to be taken to make sure a symlink does not
// traverse.
func SafeFilepathJoin(path1, path2 string) (string, error) {
	cleanPath, err := filepath.Rel("/", filepath.Clean(filepath.Join("/", path2)))
	if err != nil {
		return "", fmt.Errorf("could not clean path %q: %v", path2, err)
	}
	return filepath.Join(path1, cleanPath), nil
}

// SafeResolveSymlink recursively resolves a symlink and makes sure it never
// references a file outside of mountPointPath. If symlinkPath is not a
// symlink, it is returned immediately. If at the end of the resolution, the
// file is not regular file or directory, an error is returned.
/*
func SafeResolveSymlink(mountPointPath, symlinkPath string) (string, error) {
	// Fixed upper limit to prevent an infinite loop.
	for i := 0; i < 10; i++ {
		fi, err := os.Lstat(symlinkPath)
		if err != nil {
			return err
		}
		if fi.Mode&os.ModeSymlink != 0 {
			continue
		}
		if fi.Mode.IsRegular() || fi.Mode.IsDir() {
			return symlinkPath, nil
		}
	}
}
*/
