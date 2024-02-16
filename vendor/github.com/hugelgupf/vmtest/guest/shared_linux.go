// Copyright 2021 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package guest

import (
	"fmt"
	"os"

	"github.com/u-root/u-root/pkg/mount"
)

const (
	// https://wiki.qemu.org/Documentation/9psetup#msize recommends an
	// msize of at least 10MiB. Larger number might give better
	// performance. QEMU will print a warning if it is too small. Linux's
	// default is 8KiB which is way too small.
	msize9P = 10 * 1024 * 1024
)

// Mount9PDir mounts a directory shared as tag at dir. It creates dir if it
// does not exist.
func Mount9PDir(dir, tag string) (*mount.MountPoint, error) {
	if err := os.MkdirAll(dir, 0o644); err != nil {
		return nil, err
	}

	mp, err := mount.Mount(tag, dir, "9p", fmt.Sprintf("9P2000.L,msize=%d", msize9P), 0)
	if err != nil {
		return nil, fmt.Errorf("failed to mount directory %s: %v", dir, err)
	}
	return mp, nil
}
