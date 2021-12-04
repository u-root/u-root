// Copyright 2020 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// SPDX-License-Identifier: BSD-3-Clause
//

package boot

import (
	"log"
	"os"
	fp "path/filepath"

	"github.com/u-root/u-root/pkg/mount"
	"github.com/u-root/u-root/pkg/mount/block"
	"golang.org/x/sys/unix"
)

type EfiPathSegmentResolver interface {
	// Returns description, does not require cleanup
	String() string

	// Mount fs, etc. You must call Cleanup() eventually.
	Resolve(suggestedBasePath string) (string, error)

	// For devices, returns BlockDev. Returns nil otherwise.
	BlockInfo() *block.BlockDev

	// For mounted devices, returns MountPoint. Returns nil otherwise.
	MntPoint() *mount.MountPoint

	// Unmount fs, free resources, etc
	Cleanup() error
}

// HddResolver can identify and mount a partition.
type HddResolver struct {
	*block.BlockDev
	*mount.MountPoint
}

var _ EfiPathSegmentResolver = (*HddResolver)(nil)

func (r *HddResolver) String() string { return "/dev/" + r.BlockDev.Name }

func (r *HddResolver) Resolve(basePath string) (string, error) {
	if r.MountPoint != nil {
		return r.MountPoint.Path, nil
	}
	var err error
	if len(basePath) == 0 {
		basePath, err = os.MkdirTemp("", "uefiPath")
		if err != nil {
			return "", err
		}
	} else {
		fi, err := os.Stat(basePath)
		if err != nil || !fi.IsDir() {
			err = os.RemoveAll(basePath)
			if err != nil {
				return "", err
			}
			err = os.MkdirAll(basePath, 0o755)
			if err != nil {
				return "", err
			}
		}
	}
	r.MountPoint, err = r.BlockDev.Mount(basePath, 0)
	return r.MountPoint.Path, err
}

func (r *HddResolver) BlockInfo() *block.BlockDev { return r.BlockDev }

func (r *HddResolver) MntPoint() *mount.MountPoint { return r.MountPoint }

func (r *HddResolver) Cleanup() error {
	if r.MountPoint != nil {
		err := r.MountPoint.Unmount(unix.UMOUNT_NOFOLLOW | unix.MNT_DETACH)
		if err == nil {
			r.MountPoint = nil
		}
		return err
	}
	return nil
}

// PathResolver outputs a file path.
type PathResolver string

var _ EfiPathSegmentResolver = (*PathResolver)(nil)

func (r *PathResolver) String() string { return string(*r) }

func (r *PathResolver) Resolve(basePath string) (string, error) {
	if len(basePath) == 0 {
		log.Printf("uefi.PathResolver: empty base path")
	}
	return fp.Join(basePath, string(*r)), nil
}

func (r *PathResolver) BlockInfo() *block.BlockDev { return nil }

func (r *PathResolver) MntPoint() *mount.MountPoint { return nil }

func (r *PathResolver) Cleanup() error { return nil }
