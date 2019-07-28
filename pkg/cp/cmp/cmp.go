// Copyright 2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package cmp compares trees of files.
//
// cmp is an internal package for pkg/cp's and cmds/core/cp's tests.
package cmp

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"

	"github.com/u-root/u-root/pkg/cp"
	"github.com/u-root/u-root/pkg/uio"
)

// isEqualFile compare two files by checksum
func isEqualFile(fpath1, fpath2 string) error {
	file1, err := os.Open(fpath1)
	if err != nil {
		return err
	}
	defer file1.Close()
	file2, err := os.Open(fpath2)
	if err != nil {
		return err
	}
	defer file2.Close()

	if !uio.ReaderAtEqual(file1, file2) {
		return fmt.Errorf("%q and %q do not have equal content", fpath1, fpath2)
	}
	return nil
}

func readDirNames(path string) ([]string, error) {
	entries, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, err
	}
	var basenames []string
	for _, entry := range entries {
		basenames = append(basenames, entry.Name())
	}
	return basenames, nil
}

func stat(o cp.Options, path string) (os.FileInfo, error) {
	if o.NoFollowSymlinks {
		return os.Lstat(path)
	}
	return os.Stat(path)
}

// IsEqualTree compare the content in the file trees in src and dst paths
func IsEqualTree(o cp.Options, src, dst string) error {
	srcInfo, err := stat(o, src)
	if err != nil {
		return err
	}
	dstInfo, err := stat(o, dst)
	if err != nil {
		return err
	}
	if sm, dm := srcInfo.Mode()&os.ModeType, dstInfo.Mode()&os.ModeType; sm != dm {
		return fmt.Errorf("mismatched mode: %q has mode %s while %q has mode %s", src, sm, dst, dm)
	}

	switch {
	case srcInfo.Mode().IsDir():
		srcEntries, err := readDirNames(src)
		if err != nil {
			return err
		}
		dstEntries, err := readDirNames(dst)
		if err != nil {
			return err
		}
		// ioutil.ReadDir guarantees these are sorted.
		if !reflect.DeepEqual(srcEntries, dstEntries) {
			return fmt.Errorf("directory contents did not match:\n%q had %v\n%q had %v", src, srcEntries, dst, dstEntries)
		}
		for _, basename := range srcEntries {
			if err := IsEqualTree(o, filepath.Join(src, basename), filepath.Join(dst, basename)); err != nil {
				return err
			}
		}
		return nil

	case srcInfo.Mode().IsRegular():
		return isEqualFile(src, dst)

	case srcInfo.Mode()&os.ModeSymlink == os.ModeSymlink:
		srcTarget, err := os.Readlink(src)
		if err != nil {
			return err
		}
		dstTarget, err := os.Readlink(dst)
		if err != nil {
			return err
		}
		if srcTarget != dstTarget {
			return fmt.Errorf("target mismatch: symlink %q had target %q, while %q had target %q", src, srcTarget, dst, dstTarget)
		}
		return nil

	default:
		return fmt.Errorf("unsupported mode: %s", srcInfo.Mode())
	}
}
