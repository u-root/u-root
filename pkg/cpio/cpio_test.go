// Copyright 2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cpio

import (
	"os"
	"strings"
	"testing"
)

func TestRecord(t *testing.T) {
	r := StaticFile("file", "hello", 0o644)
	toFileMode(r)
}

func TestFormatInit(t *testing.T) {
	tests := []struct {
		format string
	}{
		{
			format: "newc",
		},
	}

	for _, test := range tests {
		_, err := Format(test.format)
		if err != nil {
			t.Errorf("expected %q in init format map, got %v", test.format, err)
		}
	}
}

func TestFormatInfo(t *testing.T) {
	f := StaticFile("file", "hello", 0o644)
	s := f.Info.String()

	if !strings.HasPrefix(s, "file") {
		t.Error("missing file name in info")
	}

	if !strings.Contains(s, "0644") {
		t.Error("missing permissions in info")
	}
}

func TestModeFromLinux(t *testing.T) {
	tests := []struct {
		mode     uint64
		fileMode os.FileMode
	}{
		{
			mode:     S_IFREG | 0o644,
			fileMode: 0o644,
		},
		{
			mode:     S_IFREG | 0o1755,
			fileMode: os.ModeSticky | 0o755,
		},
		{
			mode:     S_IFBLK | 0o644,
			fileMode: os.ModeDevice | 0o644,
		},
		{
			mode:     S_IFCHR | 0o644,
			fileMode: os.ModeDevice | os.ModeCharDevice | 0o644,
		},
		{
			mode:     S_IFDIR | 0o755,
			fileMode: os.ModeDir | 0o755,
		},
		{
			mode:     S_IFDIR | 0o2755,
			fileMode: os.ModeDir | os.ModeSetgid | 0o755,
		},
		{
			mode:     S_IFIFO | 0o600,
			fileMode: os.ModeNamedPipe | 0o600,
		},
		{
			mode:     S_IFLNK | 0o755,
			fileMode: os.ModeSymlink | 0o755,
		},
		{
			mode:     S_IFLNK | 0o1777,
			fileMode: os.ModeSymlink | os.ModeSticky | 0o777,
		},
		{
			mode:     S_IFSOCK | 0o755,
			fileMode: os.ModeSocket | 0o755,
		},
		{
			mode:     S_IFREG | 0o4755,
			fileMode: os.ModeSetuid | 0o755,
		},
		{
			mode:     S_IFREG | 0o2755,
			fileMode: os.ModeSetgid | 0o755,
		},
	}

	for _, test := range tests {
		r := modeFromLinux(test.mode)
		if r != test.fileMode {
			t.Errorf("expected %x, got %x", test.fileMode, r)
		}
	}
}
