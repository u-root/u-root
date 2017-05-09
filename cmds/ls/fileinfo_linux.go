// Copyright 2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"os"
	"os/user"
	"regexp"
	"strings"
	"syscall"
	"time"
)

// From Linux header: /include/uapi/linux/kdev_t.h
const (
	minorBits = 8
	minorMask = (1 << minorBits) - 1
)

// Matches characters which would interfere with ls's formatting.
var unprintableRe = regexp.MustCompile("[[:cntrl:]\n]")

// Since `os.FileInfo` is an interface, it is difficult to tweak some of its
// internal values. For example, replacing the starting directory with a dot.
// `extractImportantParts` populates our own struct which we can modify at will
// before printing.
type fileInfo struct {
	name         string
	mode         os.FileMode
	major, minor uint64
	uid, gid     uint32
	size         int64
	modTime      time.Time
	symlink      string
}

func extractImportantParts(n string, fi os.FileInfo) fileInfo {
	var link string

	s := fi.Sys().(*syscall.Stat_t)
	if fi.Mode()&os.ModeType == os.ModeSymlink {
		if l, err := os.Readlink(n); err != nil {
			link = err.Error()
		} else {
			link = l
		}
	}

	return fileInfo{
		name:    fi.Name(),
		mode:    fi.Mode(),
		major:   s.Rdev >> minorBits,
		minor:   s.Rdev & minorMask,
		uid:     s.Uid,
		gid:     s.Gid,
		size:    fi.Size(),
		modTime: fi.ModTime(),
		symlink: link,
	}
}

// Without this cache, `ls -l` is orders of magnitude slower.
var (
	uidCache = map[uint32]string{}
	gidCache = map[uint32]string{}
)

// Convert uid to username, or return uid on error.
func lookupUserName(id uint32) string {
	if s, ok := uidCache[id]; ok {
		return s
	}
	s := fmt.Sprint(id)
	if u, err := user.LookupId(s); err == nil {
		s = u.Username
	}
	uidCache[id] = s
	return s
}

// Convert gid to group name, or return gid on error.
func lookupGroupName(id uint32) string {
	if s, ok := gidCache[id]; ok {
		return s
	}
	s := fmt.Sprint(id)
	if g, err := user.LookupGroupId(s); err == nil {
		s = g.Name
	}
	gidCache[id] = s
	return s
}

// The default stringer. Return only the filename. Unprintable characters are
// replaced with '?'.
func (fi fileInfo) String() string {
	return unprintableRe.ReplaceAllLiteralString(fi.name, "?")
}

// Two alternative stringers
type quotedStringer struct {
	fileInfo
}
type longStringer struct {
	fileInfo
	comp fmt.Stringer // decorator pattern
}

// Return the name surrounded by quotes with escaped control characters.
func (fi quotedStringer) String() string {
	return fmt.Sprintf("%#v", fi.name)
}

// The long and quoted stringers can be combined like so:
//     longStringer{fi, quotedStringer{fi}}
func (fi longStringer) String() string {
	// Golang's FileMode.String() is almost sufficient, except we would
	// rather use b and c for devices.
	replacer := strings.NewReplacer("Dc", "c", "D", "b")

	// Ex: crw-rw-rw-  root  root  1, 3  Feb 6 09:31  null
	pattern := "%[1]s\t%[2]s\t%[3]s\t%[4]d, %[5]d\t%[7]v\t%[8]s"
	if fi.major == 0 && fi.minor == 0 {
		// Ex: -rw-rw----  myuser  myuser  1256  Feb 6 09:31  recipes.txt
		pattern = "%[1]s\t%[2]s\t%[3]s\t%[6]d\t%[7]v\t%[8]s"
	}

	s := fmt.Sprintf(pattern,
		replacer.Replace(fi.mode.String()),
		lookupUserName(fi.uid),
		lookupGroupName(fi.gid),
		fi.major,
		fi.minor,
		fi.size,
		fi.modTime.Format("Jan _2 15:04"),
		fi.comp.String())

	if fi.mode&os.ModeType == os.ModeSymlink {
		s += fmt.Sprintf(" -> %v", fi.symlink)
	}
	return s
}
