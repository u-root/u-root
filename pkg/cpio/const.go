// Copyright 2013-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cpio

// These Unix constants are needed everywhere cpio is used, Unix or not.
// But we are unable to import the unix package when plan 9 is enabled,
// so lucky us, the numbers have been the same for half a century.
// It is ok to just define them.
const (
	S_IEXEC  = 0x40
	S_IFBLK  = 0x6000
	S_IFCHR  = 0x2000
	S_IFDIR  = 0x4000
	S_IFIFO  = 0x1000
	S_IFLNK  = 0xa000
	S_IFMT   = 0xf000
	S_IFREG  = 0x8000
	S_IFSOCK = 0xc000
	S_IFWHT  = 0xe000
	S_IREAD  = 0x100
	S_IRGRP  = 0x20
	S_IROTH  = 0x4
	S_IRUSR  = 0x100
	S_IRWXG  = 0x38
	S_IRWXO  = 0x7
	S_IRWXU  = 0x1c0
	S_ISGID  = 0x400
	S_ISTXT  = 0x200
	S_ISUID  = 0x800
	S_ISVTX  = 0x200
)

// Unix mode_t bits.
const (
	modeTypeMask    = 0170000
	modeSocket      = 0140000
	modeSymlink     = 0120000
	modeFile        = 0100000
	modeBlock       = 0060000
	modeDir         = 0040000
	modeChar        = 0020000
	modeFIFO        = 0010000
	modeSUID        = 0004000
	modeSGID        = 0002000
	modeSticky      = 0001000
	modePermissions = 0000777
)
