// Copyright 2013-2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cpio

import (
	"testing"
)

func TestNormalize(t *testing.T) {
	for _, tt := range []struct {
		path string
		want string
	}{
		{
			path: "/foo/bar",
			want: "foo/bar",
		},
		{
			path: "foo////bar",
			want: "foo/bar",
		},
		{
			path: "/foo/bar/../baz",
			want: "foo/baz",
		},
		{
			path: "foo/bar/../baz",
			want: "foo/baz",
		},
		{
			path: "./foo/bar",
			want: "foo/bar",
		},
		{
			path: "foo/../../bar",
			want: "../bar",
		},
		{
			path: "",
			want: ".",
		},
		{
			path: ".",
			want: ".",
		},
	} {
		if got := Normalize(tt.path); got != tt.want {
			t.Errorf("Normalize(%q) = %q, want %q", tt.path, got, tt.want)
		}
	}
}
