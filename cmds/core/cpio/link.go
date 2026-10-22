// Copyright 2026 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build !tinygo

package main

import "os"

// linkInRoot creates newname as a hard link to oldname, with both names
// confined to dir the same way cpio.CreateFileInRoot confines records.
func linkInRoot(dir, oldname, newname string) error {
	root, err := os.OpenRoot(dir)
	if err != nil {
		return err
	}
	defer root.Close()
	return root.Link(oldname, newname)
}
