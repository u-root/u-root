// Copyright 2013-2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cpio

import (
	"syscall"
)

func sysInfo(n string, sys *syscall.Stat_t) Info {
	return Info{
		Ino:      sys.Ino,
		Mode:     uint64(sys.Mode),
		UID:      uint64(sys.Uid),
		GID:      uint64(sys.Gid),
		NLink:    sys.Nlink,
		FileSize: uint64(sys.Size),
		Major:    sys.Dev >> 8,
		Minor:    sys.Dev & 0xff,
		Rmajor:   sys.Rdev >> 8,
		Rminor:   sys.Rdev & 0xff,
		Name:     n,
	}
}
