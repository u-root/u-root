// Copyright 2013-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cpio

import "os"

func sysInfo(n string, fi os.FileInfo) Info {
	return Info{
		Mode:     uint64(fi.Mode()),
		UID:      0,
		MTime:    uint64(fi.ModTime().Second()),
		FileSize: uint64(fi.Size()),
		Name:     n,
	}
}
