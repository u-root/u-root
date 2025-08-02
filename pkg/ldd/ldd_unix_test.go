// Copyright 2009-2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build freebsd || linux

package ldd

import (
	"os"
	"path/filepath"
	"slices"
	"sort"
	"testing"
)

var cases = []struct {
	name   string
	input  string
	output []string
}{
	{
		name:   "single vdso entry",
		input:  `	linux-vdso.so.1`,
		output: []string{},
	},
	{
		name:   "duplicate vdso symlink",
		input:  `	linux-vdso.so.1 => linux-vdso.so.1`,
		output: []string{},
	},
	{
		name: "multiple entries",
		input: `	linux-vdso.so.1 => linux-vdso.so.1
	libc.so.6 => /usr/lib/libc.so.6
	/lib64/ld-linux-x86-64.so.2 => /usr/lib64/ld-linux-x86-64.so.2`,
		output: []string{"/usr/lib/libc.so.6", "/usr/lib64/ld-linux-x86-64.so.2"},
	},
	{
		name:   "entry with memory address",
		input:  `linux-vdso.so.1 => (0x00007ffe4972d000)`,
		output: []string{},
	},
}

func cmp(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}

func TestParseInterp(t *testing.T) {
	for _, c := range cases {
		out, _ := parseinterp(c.input)
		if !cmp(out, c.output) {
			t.Fatalf("'%s' expected %v, but got %v", c.name, c.output, out)
		}
	}
}

func TestFollow(t *testing.T) {
	dir := t.TempDir()
	f, err := os.CreateTemp(dir, "")
	if err != nil {
		t.Fatalf("can't create tempdir: %v", err)
	}

	sPath := filepath.Join(dir, "symlink")

	err = os.Symlink(f.Name(), sPath)
	if err != nil {
		t.Fatalf("can't create symlink: %v", err)
	}

	xs, err := follow(f.Name(), sPath, f.Name())
	if err != nil {
		t.Fatalf("expected nil got %v", err)
	}

	expected := []string{f.Name(), sPath}
	sort.Strings(expected)
	sort.Strings(xs)

	if !slices.Equal(expected, xs) {
		t.Errorf("expected: %v, got: %v", expected, xs)
	}
}
