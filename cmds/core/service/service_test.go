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

	for name, tc := range map[string]struct {
		setupDir    func(t *testing.T) string
		expectErr   bool
		outputCheck func(t *testing.T, output string)
	}{
		"all services report status successfully": {
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
		"directory not found": {
			setupDir: func(t *testing.T) string {
				return "/nonexistent/directory"
			},
			expectErr:   true,
			outputCheck: nil, // no output to check since it fails immediately
		},
	} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			serviceDir := tc.setupDir(t)

			ctx := context.Background()
			err := statusAll(ctx, serviceDir)

			if tc.expectErr && err == nil {
				t.Fatalf(
					"statusAll(ctx, %q) = nil, wanted error",
					serviceDir,
				)
			} else if !tc.expectErr && err != nil {
				t.Fatalf(
					"statusAll(ctx, %q) = %v, wanted nil",
					serviceDir, err,
				)
			}
		})
	}
}

func TestFullRestart(t *testing.T) {
	t.Parallel()

	for name, tc := range map[string]struct {
		setupDir    func(t *testing.T) string
		expectErr   bool
		outputCheck func(t *testing.T, output string)
	}{
		"all services restart successfully": {
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
		"directory not found": {
			setupDir: func(t *testing.T) string {
				return "/nonexistent/directory"
			},
			expectErr:   true,
			outputCheck: nil, // no output to check since it fails immediately
		},
	} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			serviceDir := tc.setupDir(t)

			ctx := context.Background()
			err := fullRestart(ctx, serviceDir)

			if tc.expectErr && err == nil {
				t.Fatalf(
					"fullRestart(ctx, %q) = nil, wanted error",
					serviceDir,
				)
			} else if !tc.expectErr && err != nil {
				t.Fatalf(
					"fullRestart(ctx, %q) = %v, wanted nil",
					serviceDir, err,
				)
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
	err := os.WriteFile(path, []byte(scriptContent), 0o755)
	if err != nil {
		t.Fatalf("failed to create sample script: %v", err)
	}
}
