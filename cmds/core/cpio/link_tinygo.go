// Copyright 2026 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build tinygo

package main

import (
	"os"

	"github.com/u-root/u-root/pkg/upath"
)

// linkInRoot creates newname as a hard link to oldname, with both names
// confined to dir. TinyGo has no os.Root, so this uses the same lexical
// check as cpio.CreateFileInRoot does under TinyGo.
func linkInRoot(dir, oldname, newname string) error {
	o, err := upath.SafeFilepathJoin(dir, oldname)
	if err != nil {
		return err
	}
	n, err := upath.SafeFilepathJoin(dir, newname)
	if err != nil {
		return err
	}
	return os.Link(o, n)
}
