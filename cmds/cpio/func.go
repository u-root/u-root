// Copyright 2013-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import "fmt"

func round4(n ...uint64) (ret uint64) {
	for _, v := range n {
		ret += v
	}

	ret = ((ret + 3) / 4) * 4
	return
}

func (f *File) String() string {

	return fmt.Sprintf("%s: Ino %d Mode %#o UID %d GID %d Nlink %d Mtime %#x FileSize %d Major %d Minor RMajor %d Rminor %d %d NameSize %d",
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
