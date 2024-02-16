// Copyright 2013-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cpio

import "syscall"

func sysInfo(n string, sys *syscall.Dir) Info {
	// Similar to how the standard library converts Plan 9 Dir to os.FileInfo:
	// https://github.com/golang/go/blob/go1.16beta1/src/os/stat_plan9.go#L14
	mode := sys.Mode & 0o777
	if sys.Mode&syscall.DMDIR != 0 {
		mode |= modeDir
	} else {
		mode |= modeFile
	}
	return Info{
		Mode:     uint64(mode),
		UID:      0,
		MTime:    uint64(sys.Mtime),
		FileSize: uint64(sys.Length),
		Name:     n,
	}
}
