// Copyright 2015-2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/u-root/u-root/pkg/core/mv"
)

func setup(t *testing.T) string {
	t.Helper()
	d := t.TempDir()
	for _, tt := range []struct {
		name    string
		content []byte
		mode    os.FileMode
	}{
		{
			name:    "hi1.txt",
			mode:    0o666,
			content: []byte("hi"),
		},
		{
			name:    "hi2.txt",
			mode:    0o777,
			content: []byte("hi"),
		},
		{
			name:    "old.txt",
			mode:    0o777,
			content: []byte("old"),
		},
		{
			name:    "new.txt",
			mode:    0o777,
			content: []byte("new"),
		},
	} {
		if err := os.WriteFile(filepath.Join(d, tt.name), tt.content, tt.mode); err != nil {
			t.Fatalf("setup failed: %v", err)
		}
	}
	return d
}

func TestMove(t *testing.T) {
	d := setup(t)

	for _, tt := range []struct {
		name     string
		args     []string
		wantErr  bool
		errCheck func(string) bool
	}{
		{
			name:    "Multiple files to non-directory",
			args:    []string{filepath.Join(d, "hi1.txt"), filepath.Join(d, "hi2.txt"), filepath.Join(d, "old.txt")},
			wantErr: true,
			errCheck: func(err string) bool {
				return strings.Contains(err, "not a directory")
			},
		},
		{
			name:    "Source file does not exist",
			args:    []string{filepath.Join(d, "nonexistent.txt"), filepath.Join(d, "hi2.txt")},
			wantErr: true,
			errCheck: func(err string) bool {
				return strings.Contains(err, "no such file or directory")
			},
		},
		{
			name:    "Simple move",
			args:    []string{filepath.Join(d, "hi1.txt"), filepath.Join(d, "moved.txt")},
			wantErr: false,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			cmd := mv.New()
			var stderr bytes.Buffer
			cmd.SetIO(nil, io.Discard, &stderr)

			err := cmd.Run(tt.args...)

			if tt.wantErr {
				if err == nil {
					t.Errorf("Expected error, got none")
				}
				if tt.errCheck != nil {
					errOutput := stderr.String()
					if err != nil {
						errOutput += err.Error()
					}
					if !tt.errCheck(errOutput) {
						t.Errorf("Error check failed for output: %q", errOutput)
					}
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
			}
		})
	}
}

func TestMvFlags(t *testing.T) {
	for _, tt := range []struct {
		name    string
		args    func(string) []string // Function that takes temp dir and returns args
		wantErr bool
		setup   func(string)      // Additional setup for the test, takes temp dir
		check   func(string) bool // Check the result, takes temp dir
	}{
		{
			name: "Update flag - newer source",
			args: func(d string) []string {
				return []string{"-u", filepath.Join(d, "new.txt"), filepath.Join(d, "old.txt")}
			},
			wantErr: false,
			check: func(d string) bool {
				// Check that the file was moved (new.txt should not exist, old.txt should have "new" content)
				if _, err := os.Stat(filepath.Join(d, "new.txt")); !os.IsNotExist(err) {
					return false
				}
				content, err := os.ReadFile(filepath.Join(d, "old.txt"))
				return err == nil && string(content) == "new"
			},
		},
		{
			name: "No clobber flag",
			args: func(d string) []string {
				return []string{"-n", filepath.Join(d, "hi2.txt"), filepath.Join(d, "old.txt")}
			},
			wantErr: false,
			check: func(d string) bool {
				// Check that old.txt still has "old" content (wasn't overwritten)
				content, err := os.ReadFile(filepath.Join(d, "old.txt"))
				return err == nil && string(content) == "old"
			},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			d := setup(t) // Each subtest gets its own fresh setup

			if tt.setup != nil {
				tt.setup(d)
			}

			cmd := mv.New()
			var stderr bytes.Buffer
			cmd.SetIO(nil, io.Discard, &stderr)

			err := cmd.Run(tt.args(d)...)

			if tt.wantErr {
				if err == nil {
					t.Errorf("Expected error, got none")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if tt.check != nil && !tt.check(d) {
					t.Errorf("Post-move check failed")
				}
			}
		})
	}
}

func TestMvToDirectory(t *testing.T) {
	d := setup(t)

	// Create a directory
	subdir := filepath.Join(d, "subdir")
	if err := os.Mkdir(subdir, 0o755); err != nil {
		t.Fatalf("Failed to create test directory: %v", err)
	}

	cmd := mv.New()
	var stderr bytes.Buffer
	cmd.SetIO(nil, io.Discard, &stderr)

	// Move multiple files to directory
	args := []string{filepath.Join(d, "hi1.txt"), filepath.Join(d, "hi2.txt"), subdir}
	err := cmd.Run(args...)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	// Check that files were moved to the directory
	if _, err := os.Stat(filepath.Join(subdir, "hi1.txt")); err != nil {
		t.Errorf("hi1.txt was not moved to directory: %v", err)
	}
	if _, err := os.Stat(filepath.Join(subdir, "hi2.txt")); err != nil {
		t.Errorf("hi2.txt was not moved to directory: %v", err)
	}

	// Check that original files no longer exist
	if _, err := os.Stat(filepath.Join(d, "hi1.txt")); !os.IsNotExist(err) {
		t.Errorf("Original hi1.txt still exists")
	}
	if _, err := os.Stat(filepath.Join(d, "hi2.txt")); !os.IsNotExist(err) {
		t.Errorf("Original hi2.txt still exists")
	}
}
