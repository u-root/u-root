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
		opts  []Set
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
			opts: []Set{WithModeMatch(os.ModeDir, os.ModeDir)},
			names: []string{
				"",
				"/root",
				"/root/xyz",
			},
		},
		{
			name: "just a file",
			opts: []Set{WithModeMatch(0, os.ModeType)},
			names: []string{
				"/root/xyz/0777",
				"/root/xyz/file",
			},
		},
		{
			name:  "file by mode",
			opts:  []Set{WithModeMatch(0o444, os.ModePerm)},
			names: []string{"/root/xyz/0777"},
		},
		{
			name:  "file by name",
			opts:  []Set{WithFilenameMatch("*file")},
			names: []string{"/root/xyz/file"},
		},
		{
			name:  "file by name with debug log",
			opts:  []Set{WithFilenameMatch("*file"), WithDebugLog(func(string, ...any) {})},
			names: []string{"/root/xyz/file"},
		},
		{
			name:  "file by name without error",
			opts:  []Set{WithFilenameMatch("*file"), WithoutError()},
			names: []string{"/root/xyz/file"},
		},
		{
			name:  "file by name with regex",
			opts:  []Set{WithRegexPathMatch("file")},
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
			opts := append([]Set{WithRoot(d)}, tc.opts...)
			files := Find(ctx, opts...)

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

func TestString(t *testing.T) {
	dir := t.TempDir()
	f, err := os.CreateTemp(dir, "")
	if err != nil {
		t.Fatalf("can't create file: %v", err)
	}
	fi, err := f.Stat()
	if err != nil {
		t.Fatalf("can't stat file: %v", err)
	}
	ff := File{f.Name(), fi, nil}
	s := ff.String()
	if !strings.Contains(s, f.Name()) {
		t.Errorf("expected to see %q, got %q", f.Name(), s)
	}
}
