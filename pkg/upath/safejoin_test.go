// Copyright 2020 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package upath

import (
	"testing"
)

func TestSafeFilepathJoin(t *testing.T) {
	for _, tt := range []struct {
		name               string
		path1, path2, want string
	}{
		{
			path1: "a",
			path2: "b",
			want:  "a/b",
		},
		{
			path1: "/a",
			path2: "b",
			want:  "/a/b",
		},
		{
			path1: "/a",
			path2: "/b",
			want:  "/a/b",
		},
		{
			path1: "/a",
			path2: "../b",
			want:  "/a/b",
		},
		{
			path1: "/a",
			path2: "c/../b",
			want:  "/a/b",
		},
		{
			path1: "/a",
			path2: "c/d/../b",
			want:  "/a/c/b",
		},
		{
			path1: "/a",
			path2: "c/d/../../../b",
			want:  "/a/b",
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			got, err := SafeFilepathJoin(tt.path1, tt.path2)
			if err != nil {
				t.Fatal(err)
			}
			if got != tt.want {
				t.Errorf("safePathJoin(%q, %q) = %q; want %q", tt.path1, tt.path2, got, tt.want)
			}
		})
	}
}

func TestSafeResolveSymlink(t *testing.T) {
	for _, tt := range []struct {
		name               string
		path1, path2, want string
	}{
		{
			path1: "a",
			path2: "b",
			want:  "a/b",
		},
		{
			path1: "/a",
			path2: "b",
			want:  "/a/b",
		},
		{
			path1: "/a",
			path2: "/b",
			want:  "/a/b",
		},
		{
			path1: "/a",
			path2: "../b",
			want:  "/a/b",
		},
		{
			path1: "/a",
			path2: "c/../b",
			want:  "/a/b",
		},
		{
			path1: "/a",
			path2: "c/d/../b",
			want:  "/a/c/b",
		},
		{
			path1: "/a",
			path2: "c/d/../../../b",
			want:  "/a/b",
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			got, err := SafeFilepathJoin(tt.path1, tt.path2)
			if err != nil {
				t.Fatal(err)
			}
			if got != tt.want {
				t.Errorf("safePathJoin(%q, %q) = %q; want %q", tt.path1, tt.path2, got, tt.want)
			}
		})
	}
}
