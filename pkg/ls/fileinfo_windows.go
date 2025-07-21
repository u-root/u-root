// Copyright 2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ls

import (
	"fmt"
	"os"
	"strconv"
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
	Name  string
	Mode  os.FileMode
	UID   string
	Size  int64
	MTime time.Time
}

// FromOSFileInfo converts os.FileInfo to an ls.FileInfo.
func FromOSFileInfo(path string, fi os.FileInfo) FileInfo {
	return FileInfo{
		Name:  fi.Name(),
		Mode:  fi.Mode(),
		UID:   "bill gates", // fi.Sys().(*syscall.Dir).Uid,
		Size:  fi.Size(),
		MTime: fi.ModTime(),
	}
}

// FileString implements Stringer.FileString.
func (ls LongStringer) FileString(fi FileInfo) string {
	var size string
	if ls.Human {
		size = humanize.Bytes(uint64(fi.Size))
	} else {
		size = strconv.FormatInt(fi.Size, 10)
	}
	// Ex: -rw-rw----  myuser  1256  Feb 6 09:31  recipes.txt
	return fmt.Sprintf("%s\t%s\t%s\t%v\t%s",
		fi.Mode.String(),
		fi.UID,
		size,
		fi.MTime.Format("Jan _2 15:04"),
		ls.Name.FileString(fi))
}
