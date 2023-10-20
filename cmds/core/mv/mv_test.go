// Copyright 2015-2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"os"
	"path/filepath"
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
		want  error
		name  string
		files []string
	}{
		{
			name:  "Is a directory",
			files: []string{filepath.Join(d, "hi1.txt"), filepath.Join(d, "hi1.txt"), filepath.Join(d, "hi1.txt")},
			want:  fmt.Errorf("not a directory: %s", filepath.Join(d, "hi1.txt")),
		},
		{
			name:  "Is not a directory",
			files: []string{filepath.Join(d, "hi1.txt"), filepath.Join(d, "hi3.txt"), "d"},
			want:  fmt.Errorf("not a directory: %s", "d"),
		},
		{
			name:  "mv logFatalf err",
			files: []string{filepath.Join(d, "hi1.txt"), filepath.Join(d, "hi3.txt")},
			want:  fmt.Errorf("lstat %s: no such file or directory", filepath.Join(d, "hi3.txt")),
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			if got := move(tt.files, true, false); got != nil {
				if got.Error() != tt.want.Error() {
					t.Errorf("move() = '%v', want: '%v'", got, tt.want)
				}
			}
		})
	}

}

func TestMv(t *testing.T) {
	d := setup(t)

	for _, tt := range []struct {
		want      error
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
			want:   fmt.Errorf("lstat %s: no such file or directory", filepath.Join("d", "hi1.txt")),
			update: true,
		},
		{
			name:   "len(files) = 2",
			files:  []string{filepath.Join(d, "hi1.txt"), filepath.Join(d, "hi2.txt")},
			update: true,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			if got := mv(tt.files, tt.update, tt.noClobber, tt.todir); got != nil {
				if got.Error() != tt.want.Error() {
					t.Errorf("mv() = '%v', want: '%v'", got, tt.want)
				}
			}

		})
	}
}

func TestMoveFile(t *testing.T) {
	d := setup(t)

	var testTable = []struct {
		want error
		name string
		src  string
		dst  string
	}{
		{
			name: "first file in update path does not exist",
			src:  filepath.Join(d, "hi3.txt"),
			dst:  filepath.Join(d, "hi2.txt"),
			want: fmt.Errorf("lstat %s: no such file or directory", filepath.Join(d, "hi3.txt")),
		},
		{
			name: "second file in update path does not exist",
			src:  filepath.Join(d, "hi2.txt"),
			dst:  filepath.Join(d, "hi3.txt"),
			want: fmt.Errorf("lstat %s: no such file or directory", filepath.Join(d, "hi3.txt")),
		},
	}

	for _, tt := range testTable {
		t.Run(tt.name, func(t *testing.T) {
			if got := moveFile(tt.src, tt.dst, true, false); got != nil {
				if got.Error() != tt.want.Error() {
					t.Errorf("moveFile() = '%v', want: '%v'", got, tt.want)
				}
			}
		})
	}

	t.Run("test for noClobber", func(t *testing.T) {
		if err := moveFile(testTable[0].src, testTable[0].dst, false, true); err != nil {
			t.Errorf("Expected err: %v, got: %v", err, testTable[0].want)
		}
	})
}
