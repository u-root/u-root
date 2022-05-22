// Copyright 2014-2022 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"testing"
)

func TestRun(t *testing.T) {
	for _, tt := range []struct {
		name    string
		setup   func(path string, t *testing.T) string
		list    bool
		read    string
		delete  string
		write   string
		content string
		wantErr error
	}{
		{
			name: "list no efivarfs",
			setup: func(path string, t *testing.T) string {
				t.Helper()
				return ""
			},
			list:    true,
			wantErr: fs.ErrNotExist,
		},
		{
			name: "read no efivarfs",
			setup: func(path string, t *testing.T) string {
				t.Helper()
				return ""
			},
			read:    "TestVar-bc54d3fb-ed45-462d-9df8-b9f736228350",
			wantErr: fs.ErrNotExist,
		},
		{
			name: "delete no efivarfs",
			setup: func(path string, t *testing.T) string {
				t.Helper()
				return ""
			},
			delete:  "TestVar-bc54d3fb-ed45-462d-9df8-b9f736228350",
			wantErr: fs.ErrNotExist,
		},
		{
			name: "write malformed var",
			setup: func(path string, t *testing.T) string {
				t.Helper()
				return ""
			},
			write:   "TestVar-bc54d3fb-ed45-462d-9df8-b9f736228350000",
			wantErr: fs.ErrNotExist,
		},
		{
			name: "write no content",
			setup: func(path string, t *testing.T) string {
				t.Helper()
				return ""
			},
			write:   "TestVar-bc54d3fb-ed45-462d-9df8-b9f736228350",
			content: "/bogus",
			wantErr: fs.ErrNotExist,
		},
		{
			name:  "write no guid no efivarfs",
			write: "TestVar",
			setup: func(path string, t *testing.T) string {
				t.Helper()
				f, err := os.Create(filepath.Join(path, "content"))
				if err != nil {
					t.Errorf("Failed to create file: %v", err)
				}
				s := f.Name()
				if err := f.Close(); err != nil {
					t.Errorf("Failed to close file: %v", err)
				}
				return s
			},
			wantErr: os.ErrPermission,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			tt.content = tt.setup(t.TempDir(), t)
			if err := run(tt.list, tt.read, tt.delete, tt.write, tt.content); err != nil {
				t.Logf("Try: %v, err %v", tt, err)
				if !errors.Is(err, tt.wantErr) {
					t.Errorf("Got: %v(%T), Want: %v(%T)", err, err, tt.wantErr, tt.wantErr)
				}
			}
		})
	}
}
