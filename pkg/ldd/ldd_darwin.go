// Copyright 2021 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build darwin
// +build darwin

// ldd returns none of the library dependencies of an executable.
//
// On many Unix kernels, the kernel ABI is stable. On OSX, the stability
// is held in the library interface; the kernel ABI is explicitly not
// stable. The ldd package on OSX will only return the files passed to it.
// It will continue to resolve symbolic links.
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

// Follow starts at a pathname and adds it
// to a map if it is not there.
// If the pathname is a symlink, indicated by the Readlink
// succeeding, links repeats and continues
// for as long as the name is not found in the map.
func followAllowReachable(l string, names map[string]*FileInfo, unreachableNames map[string]*FileInfo) error {
	for {
		if names[l] != nil {
			return nil
		}
		i, err := os.Lstat(l)
		if err != nil && !os.IsNotExist(err){
			unreachableNames[l] = &FileInfo{FullName: l, FileInfo: i}
			return nil
		}else if err != nil {
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
func Ldd(names []string) ([]*FileInfo, error) {
	var (
		list = make(map[string]*FileInfo)
		libs []*FileInfo
	)
	for _, n := range names {
		if err := follow(n, list); err != nil {
			return nil, err
		}else if err != nil {

		}
	}
	for i := range list {
		libs = append(libs, list[i])
	}

	return libs, nil
}

// LddAllowUnreachable returns the list of files passed to it, and resolves all symbolic
// links, returning them as well.
//
// It's not an error for a file to not be an ELF.
func LddAllowUnreachable(names []string) ([]*FileInfo, []*FileInfo, error) {
	var (
		list = make(map[string]*FileInfo)
		unreachable = make(map[string]*FileInfo)
		libs []*FileInfo
		unreachableLibs = []*FileInfo
	)
	for _, n := range names {
		if err := followAllowReachable(n, list, unreachable); err != nil {
			return nil, nil, err
		}
	}
	for i := range list {
		libs = append(libs, list[i])
	}
	for i := range unreachable {
		unreachableLibs = append(unreachable, unreachable[i])
	}

	return libs, unreachableLibs, nil
}

type FileInfo struct {
	FullName string
	os.FileInfo
}

// List returns the dependency file paths of files in names.
func List(names []string) ([]string, error) {
	var list []string
	l, err := Ldd(names)
	if err != nil {
		return nil, err
	}
	for i := range l {
		list = append(list, l[i].FullName)
	}
	return list, nil
}
