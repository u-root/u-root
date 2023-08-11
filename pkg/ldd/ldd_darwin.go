// Copyright 2021 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build darwin
// +build darwin

package ldd

import (
	"fmt"
	"os"
	"path/filepath"
)

// Follow starts at a pathname and adds it
// to a map if it is not there.
// If the pathname is a symlink, indicated by the Readlink
// succeeding, links repeats and continues
// for as long as the name is not found in the map.
func follow(l string, names map[string]*FileInfo) error {
	for {
		if names[l] != nil {
			return nil
		}
		i, err := os.Lstat(l)
		if err != nil {
			return fmt.Errorf("%v", err)
		}

		names[l] = &FileInfo{FullName: l, FileInfo: i}
		if i.Mode().IsRegular() {
			return nil
		}
		// If it's a symlink, the read works; if not, it fails.
		// we can skip testing the type, since we still have to
		// handle any error if it's a link.
		next, err := os.Readlink(l)
		if err != nil {
			return err
		}
		// It may be a relative link, so we need to
		// make it abs.
		if filepath.IsAbs(next) {
			l = next
			continue
		}
		l = filepath.Join(filepath.Dir(l), next)
	}
}

// Ldd returns the list of files passed to it, and resolves all symbolic
// links, returning them as well.
//
// It's not an error for a file to not be an ELF.
func Ldd(names ...string) ([]*FileInfo, error) {
	var (
		list = make(map[string]*FileInfo)
		libs []*FileInfo
	)
	for _, n := range names {
		if err := follow(n, list); err != nil {
			return nil, err
		}
	}
	for i := range list {
		libs = append(libs, list[i])
	}

	return libs, nil
}
