// Copyright 2026 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build tinygo

package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/u-root/u-root/pkg/upath"
)

// linkInRoot creates newname as a hard link to oldname, with both names
// confined to dir. TinyGo has no os.Root, so after the lexical check the
// parent of each name is resolved on disk, which catches a symlink unpacked
// from an earlier record (e.g. "escape" -> "..") that the name alone hides.
func linkInRoot(dir, oldname, newname string) error {
	root, err := filepath.Abs(dir)
	if err != nil {
		return err
	}
	if root, err = filepath.EvalSymlinks(root); err != nil {
		return err
	}
	o, err := inRoot(root, oldname)
	if err != nil {
		return err
	}
	n, err := inRoot(root, newname)
	if err != nil {
		return err
	}
	return os.Link(o, n)
}

// inRoot joins name to root and checks that its parent directory, with
// symlinks resolved, is still inside root.
func inRoot(root, name string) (string, error) {
	p, err := upath.SafeFilepathJoin(root, name)
	if err != nil {
		return "", err
	}
	parent, err := filepath.EvalSymlinks(filepath.Dir(p))
	if err != nil {
		return "", err
	}
	rel, err := filepath.Rel(root, parent)
	if err != nil || rel == ".." || strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
		return "", fmt.Errorf("%q: path escapes from parent", name)
	}
	return p, nil
}
