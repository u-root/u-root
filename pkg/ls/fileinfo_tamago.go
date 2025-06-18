// Copyright 2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build tamago

package ls

import (
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
	"syscall"
	"time"

	humanize "github.com/dustin/go-humanize"
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
	MTime         time.Time
	SymlinkTarget string
}

// FromOSFileInfo converts os.FileInfo to an ls.FileInfo.
func FromOSFileInfo(path string, fi os.FileInfo) FileInfo {
	var link string

	// A filesystem with a bug will result
	// in sys not being the right type.
	// This turns out to be surprisingly messy to test.
	uID, gID, rdev := uint32(math.MaxUint32), uint32(math.MaxUint32), uint64(math.MaxUint64)
	if s, ok := fi.Sys().(*syscall.Stat_t); ok {
		uID, gID, rdev = s.Uid, s.Gid, uint64(s.Rdev)
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
		MTime:         fi.ModTime(),
		SymlinkTarget: link,
	}
}

// FileString implements Stringer.FileString.
func (ls LongStringer) FileString(fi FileInfo) string {
	// Golang's FileMode.String() is almost sufficient, except we would
	// rather use b and c for devices.
	replacer := strings.NewReplacer("Dc", "c", "D", "b")

	// Ex: crw-rw-rw-  root  root  1, 3  Feb 6 09:31  null
	pattern := "%[1]s\t%[2]s\t%[3]s\t%[4]d, %[5]d\t%[7]v\t%[8]s"
	if fi.Mode&os.ModeDevice == 0 && fi.Mode&os.ModeCharDevice == 0 {
		// Ex: -rw-rw----  myuser  myuser  1256  Feb 6 09:31  recipes.txt
		pattern = "%[1]s\t%[2]s\t%[3]s\t%[6]s\t%[7]v\t%[8]s"
	}

	var size string
	if ls.Human {
		size = humanize.Bytes(uint64(fi.Size))
	} else {
		size = strconv.FormatInt(fi.Size, 10)
	}

	s := fmt.Sprintf(pattern,
		replacer.Replace(fi.Mode.String()),
		lookupUserName(fi.UID),
		lookupGroupName(fi.GID),
		0, // unix.Major(fi.Rdev),
		0, // unix.Minor(fi.Rdev),
		size,
		fi.MTime.Format("Jan _2 15:04"),
		ls.Name.FileString(fi))

	if fi.Mode&os.ModeType == os.ModeSymlink {
		s += fmt.Sprintf(" -> %v", fi.SymlinkTarget)
	}
	return s
}
