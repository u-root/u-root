// Copyright 2026 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build js && tinygo

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
	ATime         time.Time
	MTime         time.Time
	CTime         time.Time
	SymlinkTarget string
	Dev           uint64
	Ino           uint64
	NLink         uint64
}

// FromOSFileInfo converts os.FileInfo to an ls.FileInfo (TinyGo js/wasm
// variant).
//
// It differs from the js variant only in how the timestamps are spelled:
// TinyGo's syscall.Stat_t carries Atim and Ctim as a Timespec, where Go's
// js implementation has whole-second and nanosecond fields. Every other
// field is read the same way.
func FromOSFileInfo(path string, fi os.FileInfo) FileInfo {
	var link string

	uID, gID, rdev := uint32(math.MaxUint32), uint32(math.MaxUint32), uint64(math.MaxUint64)
	var aTime, cTime time.Time
	var dev, ino, nLink uint64
	var blkSize, blocks int64
	if s, ok := fi.Sys().(*syscall.Stat_t); ok {
		uID, gID, rdev = uint32(s.Uid), uint32(s.Gid), uint64(s.Rdev)
		aTime = time.Unix(int64(s.Atim.Sec), s.Atim.Nsec)
		cTime = time.Unix(int64(s.Ctim.Sec), s.Ctim.Nsec)
		dev = uint64(s.Dev)
		ino = s.Ino
		nLink = uint64(s.Nlink)
		blkSize = int64(s.Blksize)
		blocks = int64(s.Blocks)
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
		ATime:         aTime,
		CTime:         cTime,
		SymlinkTarget: link,
		Dev:           dev,
		Ino:           ino,
		NLink:         nLink,
	}
}
