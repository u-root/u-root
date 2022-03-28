// Copyright 2021 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package upath

import (
	"testing"
)

func TestSafeFilepathJoin(t *testing.T) {
	for _, tt := range []struct {
		name         string
		path1, path2 string
		wantPath     string
		wantErr      bool
	}{
		{
			name:     "safe relative paths",
			path1:    "a",
			path2:    "b",
			wantPath: "a/b",
		},
		{
			name:     "safe relative paths 2",
			path1:    "./a",
			path2:    "./b",
			wantPath: "a/b",
		},
		{
			name:    "unsafe absolute paths",
			path1:   "/a",
			path2:   "/b",
			wantErr: true,
		},
		{
			name:     "safe absolute path",
			path1:    "/a",
			path2:    "b",
			wantPath: "/a/b",
		},
		{
			name:    "unsafe absolute path",
			path1:   "a",
			path2:   "/b",
			wantErr: true,
		},
		{
			name:    "unsafe dotdot escape",
			path1:   "/a",
			path2:   "../b",
			wantErr: true,
		},
		{
			name:    "unsafe dotdot escape 2",
			path1:   "/a",
			path2:   "c/d/../../../b",
			wantErr: true,
		},
		{
			name:    "unsafe dotdot escape 3",
			path1:   "/a",
			path2:   "c/d/../../../a/b",
			wantErr: true,
		},
		{
			name:     "safe dotdot",
			path1:    "/a",
			path2:    "c/../b",
			wantPath: "/a/b",
		},
		{
			name:     "safe dotdot 2",
			path1:    "/a",
			path2:    "c/d/../b",
			wantPath: "/a/c/b",
		},
		{
			name:     "safe dotdot 3",
			path1:    "../a",
			path2:    "c/d/../b",
			wantPath: "../a/c/b",
		},
		{
			name:     "safe missing path",
			path1:    "",
			path2:    "b",
			wantPath: "b",
		},
		{
			name:     "safe missing path 2",
			path1:    "a",
			path2:    "",
			wantPath: "a",
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			got, err := SafeFilepathJoin(tt.path1, tt.path2)
			if (err == nil) == tt.wantErr {
				t.Errorf("safePathJoin(%q, %q) = err %v; wantErr=%v", tt.path1, tt.path2, err, tt.wantErr)
			}
			if got != tt.wantPath {
				t.Errorf("safePathJoin(%q, %q) = %q; want %q", tt.path1, tt.path2, got, tt.wantPath)
			}
		})
	}
}
