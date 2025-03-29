// Copyright 2025 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ls

import (
	"math"
	"os"
	"syscall"
	"time"
)

// FileInfo holds file metadata.
//
// Since `os.FileInfo` is an interface, it is difficult to tweak some of its
// internal values. For example, replacing the starting directory with a dot.
// `extractImportantParts` populates our own struct which we can modify at will
// before printing.
type FileInfo struct {
	Name          string
	Mode          os.FileMode
	Rdev          uint64
	UID, GID      uint32
	Size          int64
	BlkSize       int64
	Blocks        int64
	MTime         time.Time
	SymlinkTarget string
	Dev           uint64
	Ino           uint64
	NLink         uint64
}

// FromOSFileInfo converts os.FileInfo to an ls.FileInfo.
func FromOSFileInfo(path string, fi os.FileInfo) FileInfo {
	var link string

	// A filesystem with a bug will result
	// in sys not being the right type.
	// This turns out to be surprisingly messy to test.
	uID, gID, rdev := uint32(math.MaxUint32), uint32(math.MaxUint32), uint64(math.MaxUint64)
	var dev, ino, nLink uint64
	var blkSize, blocks int64
	if s, ok := fi.Sys().(*syscall.Stat_t); ok {
		uID, gID, rdev = s.Uid, s.Gid, uint64(s.Rdev)
		dev = uint64(s.Dev)
		ino = s.Ino
		nLink = uint64(s.Nlink)
		blkSize = int64(s.Blksize)
		blocks = int64(blocks)
	}

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
		Rdev:          rdev,
		UID:           uID,
		GID:           gID,
		Size:          fi.Size(),
		BlkSize:       blkSize,
		Blocks:        blocks,
		MTime:         fi.ModTime(),
		SymlinkTarget: link,
		Dev:           dev,
		Ino:           ino,
		NLink:         nLink,
	}
}
