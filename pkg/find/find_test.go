// Copyright 2015-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package find

import (
	"context"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"syscall"
	"testing"
)

func TestSimple(t *testing.T) {
	type tests struct {
		name  string
		opts  Set
		names []string
	}

	testCases := []tests{
		{
			name: "basic find",
			opts: nil,
			names: []string{
				"",
				"/root",
				"/root/xyz",
				"/root/xyz/0777",
				"/root/xyz/file",
			},
		},
		{
			name: "just a dir",
			opts: WithModeMatch(os.ModeDir, os.ModeDir),
			names: []string{
				"",
				"/root",
				"/root/xyz",
			},
		},
		{
			name: "just a file",
			opts: WithModeMatch(0, os.ModeType),
			names: []string{
				"/root/xyz/0777",
				"/root/xyz/file",
			},
		},
		{
			name:  "file by mode",
			opts:  WithModeMatch(0o444, os.ModePerm),
			names: []string{"/root/xyz/0777"},
		},
		{
			name:  "file by name",
			opts:  WithFilenameMatch("*file"),
			names: []string{"/root/xyz/file"},
		},
	}
	d := t.TempDir()

	// Make sure files are actually created with the permissions we ask for.
	syscall.Umask(0)
	if err := os.MkdirAll(filepath.Join(d, "root/xyz"), 0o775); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(d, "root/xyz/file"), nil, 0o664); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(d, "root/xyz/0777"), nil, 0o444); err != nil {
		t.Fatal(err)
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()
			files := Find(ctx, WithRoot(d), tc.opts)

			var names []string
			for o := range files {
				if o.Err != nil {
					t.Errorf("%v: got %v, want nil", o.Name, o.Err)
				}
				names = append(names, strings.TrimPrefix(o.Name, d))
			}

			if len(names) != len(tc.names) {
				t.Errorf("Find output: got %d bytes, want %d bytes", len(names), len(tc.names))
			}
			if !reflect.DeepEqual(names, tc.names) {
				t.Errorf("Find output: got %v, want %v", names, tc.names)
			}
		})
	}
}
