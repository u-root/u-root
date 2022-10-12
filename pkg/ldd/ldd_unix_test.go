// Copyright 2009-2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build freebsd || linux
// +build freebsd linux

package ldd

import (
	"testing"
)

var (
	cases = []struct {
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
	}
)

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
