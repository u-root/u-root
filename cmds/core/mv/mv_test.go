// Copyright 2015-2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"errors"
	"io"
	"os"
	"path/filepath"
	"syscall"
	"testing"
)

func setup(t *testing.T) string {
	t.Helper()
	d := t.TempDir()
	for _, tt := range []struct {
		name    string
		content []byte
		mode    os.FileMode
	}{
		{
			name:    "hi1.txt",
			mode:    0o666,
			content: []byte("hi"),
		},
		{
			name:    "hi2.txt",
			mode:    0o777,
			content: []byte("hi"),
		},
		{
			name:    "old.txt",
			mode:    0o777,
			content: []byte("old"),
		},
		{
			name:    "new.txt",
			mode:    0o777,
			content: []byte("new"),
		},
	} {
		if err := os.WriteFile(filepath.Join(d, tt.name), tt.content, tt.mode); err != nil {
			t.Fatalf("setup failed: %v", err)
		}
	}
	return d
}

func TestMove(t *testing.T) {
	d := setup(t)

	for _, tt := range []struct {
		err   error
		name  string
		args  []string
		files []string
	}{
		{
			name:  "Is a directory",
			files: []string{filepath.Join(d, "hi1.txt"), filepath.Join(d, "hi1.txt"), filepath.Join(d, "hi1.txt")},
			args:  []string{"-u"},
			err:   syscall.ENOTDIR,
		},
		{
			name:  "Is not a directory",
			files: []string{filepath.Join(d, "hi1.txt"), filepath.Join(d, "hi3.txt"), "d"},
			args:  []string{"-u"},
			err:   syscall.ENOTDIR,
		},
		{
			name:  "mv logFatalf err",
			files: []string{filepath.Join(d, "hi1.txt"), filepath.Join(d, "hi3.txt")},
			args:  []string{"-u"},
			err:   os.ErrNotExist,
		},
		{
			name:  "no flags no error",
			files: []string{filepath.Join(d, "hi1.txt"), filepath.Join(d, "hi1_new.txt")},
		},
		{
			name:  "not enough args",
			files: []string{filepath.Join(d, "hi1.txt")},
			args:  []string{"-un"},
			err:   errUsage,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			tt.args = append(tt.args, tt.files...)
			err := move(io.Discard, tt.args)
			if !errors.Is(err, tt.err) {
				t.Errorf("move() = '%v', want: '%v'", err, tt.err)
			}
		})
	}

}

func TestMv(t *testing.T) {
	d := setup(t)

	for _, tt := range []struct {
		err       error
		name      string
		files     []string
		update    bool
		noClobber bool
		todir     bool
	}{
		{
			name:   "len(files) > 2",
			files:  []string{filepath.Join(d, "hi1.txt"), filepath.Join(d, "hi2.txt"), d},
			update: true,
		},
		{
			name:   "len(files) > 2 && d does not exist",
			files:  []string{filepath.Join(d, "hi1.txt"), filepath.Join(d, "hi2.txt"), "d"},
			err:    os.ErrNotExist,
			update: true,
		},
		{
			name:   "len(files) = 2",
			files:  []string{filepath.Join(d, "hi1.txt"), filepath.Join(d, "hi2.txt")},
			update: true,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			err := mv(tt.files, tt.update, tt.noClobber, tt.todir)
			if !errors.Is(err, tt.err) {
				t.Errorf("mv() = '%v', want: '%v'", err, tt.err)
			}
		})
	}
}

func TestMoveFile(t *testing.T) {
	d := setup(t)

	var testTable = []struct {
		err  error
		name string
		src  string
		dst  string
	}{
		{
			name: "first file in update path does not exist",
			src:  filepath.Join(d, "hi3.txt"),
			dst:  filepath.Join(d, "hi2.txt"),
			err:  os.ErrNotExist,
		},
		{
			name: "second file in update path does not exist",
			src:  filepath.Join(d, "hi2.txt"),
			dst:  filepath.Join(d, "hi3.txt"),
			err:  os.ErrNotExist,
		},
	}

	for _, tt := range testTable {
		t.Run(tt.name, func(t *testing.T) {
			err := moveFile(tt.src, tt.dst, true, false)
			if !errors.Is(err, tt.err) {
				t.Errorf("moveFile() = '%v', want: '%v'", err, tt.err)
			}
		})
	}

	t.Run("test for noClobber", func(t *testing.T) {
		if err := moveFile(testTable[0].src, testTable[0].dst, false, true); err != nil {
			t.Errorf("Expected err: %v, got: %v", err, testTable[0].err)
		}
	})
}
