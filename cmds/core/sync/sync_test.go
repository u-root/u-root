// Copyright 2022 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"os"
	"testing"

	"github.com/u-root/u-root/pkg/testutil"
)

func TestSync(t *testing.T) {
	testutil.SkipIfNotRoot(t)
	d, err := os.MkdirTemp(os.TempDir(), "sync")
	if err != nil {
		t.Errorf("Failed to create tmp folder: %v", err)
	}
	file1, err := os.CreateTemp(d, "file1")
	if err != nil {
		t.Errorf("failed to create tmp file1: %v", err)
	}
	file2, err := os.CreateTemp(d, "file2")
	if err != nil {
		t.Errorf("failed to create tmp file2: %v", err)
	}
	defer os.RemoveAll(d)
	for _, tt := range []struct {
		name       string
		input      []string
		want       error
		data       bool
		filesystem bool
	}{
		{
			name:  "data flag",
			input: []string{file1.Name(), file2.Name()},
			want:  nil,
			data:  true,
		},
		{
			name:  "data flag with wrong path",
			input: []string{"file1"},
			want:  fmt.Errorf("open file1: no such file or directory"),
			data:  true,
		},
		{
			name:       "filesystem flag",
			input:      []string{file1.Name(), file2.Name()},
			want:       nil,
			filesystem: true,
		},
		{
			name:  "default",
			input: []string{file1.Name(), file2.Name()},
			want:  nil,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			*data = tt.data
			*filesystem = tt.filesystem
			if got := sync(tt.input); got != nil {
				if tt.want.Error() != got.Error() {
					t.Errorf("sync() = '%v', want: '%v'", got, tt.want)
				}
			}
		})
	}
}
