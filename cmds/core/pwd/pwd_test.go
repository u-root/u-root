// Copyright 2012-2023 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestPWD(t *testing.T) {
	tests := []struct {
		name     string
		physical bool
	}{
		{
			name:     "follow-symlinks",
			physical: true,
		},
		{
			name:     "no-follow-symlinks",
			physical: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := t.TempDir()
			err := os.Chdir(dir)
			if err != nil {
				t.Fatalf("os.Chdir(%q): got %v, want nil", dir, err)
			}

			// required on macOS, where tempDir returns path with symlink
			// and PWD is defined on bash and zsh
			err = os.Setenv("PWD", dir)
			if err != nil {
				t.Fatalf("os.Setenv(%q): got %v, want nil", dir, err)
			}

			expected := dir
			if tt.physical {
				expected, err = filepath.EvalSymlinks(dir)
				if err != nil {
					t.Fatalf("filepath.EvalSymlinks(%q): got %v, want nil", dir, err)
				}
			}

			res, err := pwd(tt.physical)
			if err != nil {
				t.Fatalf("pwd(%t): error got %v, want nil", tt.physical, err)
			}
			if res != expected {
				t.Errorf("pwd(%t): got %q, want %q", tt.physical, res, expected)
			}
		})
	}
}
