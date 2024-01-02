// Copyright 2022 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"os"
	"testing"

	"github.com/hugelgupf/vmtest/guest"
)

func TestSync(t *testing.T) {
	guest.SkipIfNotInVM(t)

	d := t.TempDir()
	file1, err := os.CreateTemp(d, "file1")
	if err != nil {
		t.Errorf("failed to create tmp file1: %v", err)
	}
	file2, err := os.CreateTemp(d, "file2")
	if err != nil {
		t.Errorf("failed to create tmp file2: %v", err)
	}

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
			if got := sync(tt.data, tt.filesystem, tt.input); got != nil {
				if tt.want.Error() != got.Error() {
					t.Errorf("sync() = '%v', want: '%v'", got, tt.want)
				}
			}
		})
	}
}
