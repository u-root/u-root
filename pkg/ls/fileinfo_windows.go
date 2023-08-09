// Copyright 2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ls

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"time"

	humanize "github.com/dustin/go-humanize"
)

// Matches characters which would interfere with ls's formatting.
var unprintableRe = regexp.MustCompile("[[:cntrl:]\n]")

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
		UID:   "bill gates", //fi.Sys().(*syscall.Dir).Uid,
		Size:  fi.Size(),
		MTime: fi.ModTime(),
	}
}

// PrintableName returns a printable file name.
func (fi FileInfo) PrintableName() string {
	return unprintableRe.ReplaceAllLiteralString(fi.Name, "?")
}

// Stringer provides a consistent way to format FileInfo.
type Stringer interface {
	// FileString formats a FileInfo.
	FileString(fi FileInfo) string
}

// NameStringer is a Stringer implementation that just prints the name.
type NameStringer struct{}

// FileString implements Stringer.FileString and just returns fi's name.
func (ns NameStringer) FileString(fi FileInfo) string {
	return fi.PrintableName()
}

// QuotedStringer is a Stringer that returns the file name surrounded by qutoes
// with escaped control characters.
type QuotedStringer struct{}

// FileString returns the name surrounded by quotes with escaped control characters.
func (qs QuotedStringer) FileString(fi FileInfo) string {
	return fmt.Sprintf("%#v", fi.Name)
}

// LongStringer is a Stringer that returns the file info formatted in `ls -l`
// long format.
type LongStringer struct {
	Human bool
	Name  Stringer
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
