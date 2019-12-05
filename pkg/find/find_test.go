// Copyright 2015-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package find

import (
	"context"
	"io/ioutil"
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

	var testCases = []tests{
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
			opts:  WithModeMatch(0444, os.ModePerm),
			names: []string{"/root/xyz/0777"},
		},
		{
			name:  "file by name",
			opts:  WithFilenameMatch("*file"),
			names: []string{"/root/xyz/file"},
		},
	}
	d, err := ioutil.TempDir(os.TempDir(), "u-root.cmds.find")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(d)

	// Make sure files are actually created with the permissions we ask for.
	syscall.Umask(0)
	if err := os.MkdirAll(filepath.Join(d, "root/xyz"), 0775); err != nil {
		t.Fatal(err)
	}
	if err := ioutil.WriteFile(filepath.Join(d, "root/xyz/file"), nil, 0664); err != nil {
		t.Fatal(err)
	}
	if err := ioutil.WriteFile(filepath.Join(d, "root/xyz/0777"), nil, 0444); err != nil {
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
