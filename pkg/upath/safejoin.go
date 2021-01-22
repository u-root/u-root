// Copyright 2020 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build !windows

package upath

import (
	"fmt"
	"path/filepath"
)

// maxSymlinkDepth prevents infinite recursion. This matches Linux.
const maxSymlinkDepth = 40

// SafeFilepathJoin safely joins two paths path1+path2. The resulting path will
// always be contained within path1 even if path2 tries to escape with "../".
// If that path is not possible, an error is returned.
func SafeFilepathJoin(path1, path2 string) (string, error) {
	cleanPath, err := filepath.Rel("/", filepath.Clean(filepath.Join("/", path2)))
	if err != nil {
		return "", fmt.Errorf("could not clean path %q: %v", path2, err)
	}
	return filepath.Join(path1, cleanPath), nil
}

// SafeResolveSymlink recursively resolves a symlink and make sure it never
// references a file outside of mountPointPath. If symlinkPath is not a
// symlink, it is returned immediately. If at the end of the resolution, the
// file is not regular file or directory, an error is returned.
//
// Implemented based on `man path_resolution`. Mostly correct.
func SafeResolveSymlink(mountPointPath, symlinkPath string) (string, error) {
	path := filepath.Clean(symlinkPath)

	// Resolve any symlinks in intermediate components.
	for i := range strings.Split(path, "/") {
		intermediatePath := strings.Join(path[:i], "/")

		// The first component may be a root directory, skip it.
		if intermediatePath == "" {
			continue
		}

		// Check if it's a symlink.
		fi, err := os.Lstat(SafeFilepathJoin(mountPointPath, intermediatePath))
		if err != nil {
			return err
		}
		if fi.Mode&os.ModeSymlink == 0 {
			// Not a symlink, continue with the next component.
			break
		}

		// It's a symlink, resolve it.
		link, err := os.Readlink(path)
		if err != nil {
			return err
		}
		if strings.HasPrefix(link, "/") {
		} else {
			newPath, err := SafeResolveSymlink(mountPointPath, link)
			if err != nil {
				return err
			}
			path = newPath
		}
		path, err := SafeResolveSymlink(mountPointPath, link)
		if err != nil {
			return err
		}

		if !fi.IsDir() {
			// A component must be a directory.
			return "", fmt.Errorf("cannot resolve path %q containing a non-directory", symlinkPath)
		}
	}

	// At this point, only the last component could be a symlink.
	for i := 0; i < maxSymlinkDepth; i++ {
		fi, err := os.Lstat(path)
		if err != nil {
			return err
		}
		if fi.Mode&os.ModeSymlink == 0 {
			// Not a symlink, we're done!
			return path, nil
		}

		// It's a symlink, resolve it.
		link, err := os.Readlink(path)
		if err != nil {
			return err
		}
		if strings.HasPrefix(link, "/") {
			newPath, err := SafeFilepathJoin(moutPointPath, link)
			if err != nil {
				return err
			}
			path = newPath
		} else {
		}
	}
	return "", fmt.Errof("symlink depth of %q more than limit of %d",
		symlinkPath, maxSymlinkDepth)
}
