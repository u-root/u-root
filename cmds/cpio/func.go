// Copyright 2013-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"os"
)

const (
	typeMask    = 0170000 //This masks the file type bits.
	typeSocket  = 0140000 //File type value for sockets.
	typeSymLink = 0120000 //File type value for symbolic links.  For symbolic links, the link body is stored as file data.
	typeFile    = 0100000 //File type value for regular files.
	typeBlock   = 0060000 // File type value for block special devices.
	typeDir     = 0040000 // File type value for directories.
	typeChar    = 0020000 // File type value for character special devices.
	typeFIFO    = 0010000 // File type value for named pipes or FIFOs.
	SUID        = 0004000 // SUID bit.
	SGID        = 0002000 // SGID bit.
	sticky      = 0001000 // Sticky bit.  On some systems, this modifies the behavior
	// of executables and/or directories.
	mode = 0000777 //The lower 9 bits specify read/write/execute permissions
	//for world, group, and user following standard POSIX con-
	//ventions.
)

var (
	ModeMap = map[uint64]os.FileMode{
		typeSocket:  os.ModeSocket,
		typeSymLink: os.ModeNamedPipe,
		typeFile:    os.FileMode(0),
		typeBlock:   os.ModeDevice,
		typeDir:     os.ModeDir,
		typeChar:    os.ModeCharDevice,
		typeFIFO:    os.ModeNamedPipe,
	}
)

func perm(f *File) uint32 {
	return uint32(f.Mode) & mode
}

func dev(f *File) int {
	return int(f.Rmajor<<8 | f.Rminor)
}

func cpioModetoMode(m uint64) (os.FileMode, error) {
	if t, ok := ModeMap[m&typeMask]; ok {
		return t, nil
	}
	return os.FileMode(0), fmt.Errorf("Invalid file type %#x", m&typeMask)
}

// round4 is intended to be use with offsets to WriteAt and ReadAt
func round4(n ...int64) (ret int64) {
	for _, v := range n {
		ret += v
	}

	ret = ((ret + 3) / 4) * 4
	return
}

func (f *File) String() string {

	return fmt.Sprintf("%s: Ino %d Mode %#o UID %d GID %d Nlink %d Mtime %#x FileSize %d Major %d Minor %d RMajor %d Rminor %d NameSize %d",
		f.Name,
		f.Ino,
		f.Mode,
		f.UID,
		f.GID,
		f.Nlink,
		// what a mess. This fails on travis.
		//time.Unix(int64(f.Mtime), 0).String(),
		f.Mtime,
		f.FileSize,
		f.Major,
		f.Minor,
		f.Rmajor,
		f.Rminor,
		f.NameSize)
}
