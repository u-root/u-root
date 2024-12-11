// Copyright 2014-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mount

import (
	"fmt"
	"os"
	"strings"
)

// FindFileSystem returns nil if a file system is available for use.
//
// It rereads /proc/filesystems each time as the supported file systems can change
// as modules are added and removed.
func FindFileSystem(fstype string) error {
	b, err := os.ReadFile("/proc/filesystems")
	if err != nil {
		return err
	}
	return internalFindFileSystem(string(b), fstype)
}

func internalFindFileSystem(content string, fstype string) error {
	for _, l := range strings.Split(content, "\n") {
		f := strings.Fields(l)
		if (len(f) > 1 && f[0] == "nodev" && f[1] == fstype) || (len(f) > 0 && f[0] != "nodev" && f[0] == fstype) {
			return nil
		}
	}
	return fmt.Errorf("file system type %q not found", fstype)
}

// GetBlockFilesystems returns the supported file systems for block devices.
func GetBlockFilesystems() (fstypes []string, err error) {
	return internalGetFilesystems("/proc/filesystems")
}

func internalGetFilesystems(file string) (fstypes []string, err error) {
	var bytes []byte
	if bytes, err = os.ReadFile(file); err != nil {
		return nil, fmt.Errorf("failed to read supported file systems: %w", err)
	}
	for _, line := range strings.Split(string(bytes), "\n") {
		// len(fields)==1, 2 possibilites for fs: "nodev" fs and
		// fs's. "nodev" fs cannot be mounted through devices.
		// len(fields)==1 prevents this from occurring.
		if fields := strings.Fields(line); len(fields) == 1 {
			fstypes = append(fstypes, fields[0])
		}
	}
	return fstypes, nil
}
