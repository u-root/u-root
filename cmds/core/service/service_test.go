// Copyright 2025 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"context"
	"os"
	"path/filepath"
	"testing"
)

func TestStatusAll(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		setupDir    func(t *testing.T) string
		expectErr   bool
		outputCheck func(t *testing.T, output string)
	}{
		{
			name: "all services report status successfully",
			setupDir: func(t *testing.T) string {
				dir := t.TempDir()
				createSampleScript(t, dir, "service1.sh")
				createSampleScript(t, dir, "service2.sh")
				return dir
			},
			expectErr: false,
			outputCheck: func(t *testing.T, output string) {
				if output != "" {
					t.Errorf("expected no error output, got: %s", output)
				}
			},
		},
		{
			name: "directory not found",
			setupDir: func(t *testing.T) string {
				return "/nonexistent/directory"
			},
			expectErr:   true,
			outputCheck: nil, // no output to check since it fails immediately
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			serviceDir := tt.setupDir(t)

			ctx := context.Background()
			err := statusAll(ctx, serviceDir)

			if tt.expectErr && err == nil {
				t.Errorf("expected error but got none")
			} else if !tt.expectErr && err != nil {
				t.Errorf("did not expect error but got: %v", err)
			}
		})
	}
}

func TestFullRestart(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		setupDir    func(t *testing.T) string
		expectErr   bool
		outputCheck func(t *testing.T, output string)
	}{
		{
			name: "all services restart successfully",
			setupDir: func(t *testing.T) string {
				dir := t.TempDir()
				createSampleScript(t, dir, "service1.sh")
				createSampleScript(t, dir, "service2.sh")
				return dir
			},
			expectErr: false,
			outputCheck: func(t *testing.T, output string) {
				if output != "" {
					t.Errorf("expected no error output, got: %s", output)
				}
			},
		},
		{
			name: "directory not found",
			setupDir: func(t *testing.T) string {
				return "/nonexistent/directory"
			},
			expectErr:   true,
			outputCheck: nil, // no output to check since it fails immediately
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			serviceDir := tt.setupDir(t)

			ctx := context.Background()
			err := fullRestart(ctx, serviceDir)

			if tt.expectErr && err == nil {
				t.Errorf("expected error but got none")
			} else if !tt.expectErr && err != nil {
				t.Errorf("did not expect error but got: %v", err)
			}
		})
	}
}

func createSampleScript(t *testing.T, dir, name string) {
	t.Helper()
	scriptContent := `#!/bin/sh
case "$1" in
  start) echo "Service started";;
  stop) echo "Service stopped";;
  restart) echo "Service restarted";;
  status) echo "Service status: running"; exit 0;;
  *) echo "Invalid command"; exit 1;;
esac`
	path := filepath.Join(dir, name)
	err := os.WriteFile(path, []byte(scriptContent), 0755)
	if err != nil {
		t.Fatalf("failed to create sample script: %v", err)
	}
}
