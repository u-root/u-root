// Copyright 2026 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build js && tinygo

package ls

import (
	"math"
	"os"
	"time"
)

// FileInfo holds file metadata (tinygo js variant: no syscall.Stat_t
// access — virtual filesystems return nil from Sys anyway).
type FileInfo struct {
	Name          string
	Mode          os.FileMode
	Rdev          uint64
	UID, GID      uint32
	Size          int64
	BlkSize       int64
	Blocks        int64
	ATime         time.Time
	MTime         time.Time
	CTime         time.Time
	SymlinkTarget string
	Dev           uint64
	Ino           uint64
	NLink         uint64
}

// FromOSFileInfo converts os.FileInfo to an ls.FileInfo.
func FromOSFileInfo(path string, fi os.FileInfo) FileInfo {
	var link string
	if fi.Mode()&os.ModeType == os.ModeSymlink {
		if l, err := os.Readlink(path); err != nil {
			link = err.Error()
		} else {
			link = l
		}
	}
	return FileInfo{
		Name:          fi.Name(),
		Mode:          fi.Mode(),
		Rdev:          uint64(math.MaxUint64),
		UID:           math.MaxUint32,
		GID:           math.MaxUint32,
		Size:          fi.Size(),
		MTime:         fi.ModTime(),
		SymlinkTarget: link,
	}
}
