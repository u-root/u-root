// Copyright 2025 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

func createSampleScript(t *testing.T, dir, name string) string {
	t.Helper()
	scriptPath := filepath.Join(dir, name)
	scriptContent := `#!/bin/sh
### BEGIN INIT INFO
# Provides: example
# Default-Start: 2 3
# Default-Stop: 0 6
# Required-Start: network
# Required-Stop: network
### END INIT INFO
case "$1" in
  start) echo "Service started";;
  stop) echo "Service stopped";;
  restart) echo "Service restarted";;
  status) echo "Service status: running"; exit 0;;
  *) echo "Invalid command"; exit 1;;
esac`

	if err := os.WriteFile(scriptPath, []byte(scriptContent), 0o755); err != nil {
		t.Fatalf("failed to create script: %v", err)
	}
	return scriptPath
}

func TestDefaults(t *testing.T) {
	t.Parallel()

	for name, tc := range map[string]struct {
		setupOpts func(t *testing.T) options
		expectErr bool
	}{
		"valid script creates symlinks": {
			setupOpts: func(t *testing.T) options {
				tmpDir := t.TempDir()
				etcDir := filepath.Join(tmpDir, "etc")
				serviceDir := filepath.Join(tmpDir, "services")
				os.MkdirAll(etcDir, 0o755)
				os.MkdirAll(serviceDir, 0o755)
				for i := 0; i <= 6; i++ {
					os.MkdirAll(filepath.Join(etcDir, fmt.Sprintf("rc%d.d", i)), 0o755)
				}
				createSampleScript(t, serviceDir, "testscript")
				return options{etc: etcDir, serviceDir: serviceDir}
			},
			expectErr: false,
		},
		"missing script returns error": {
			setupOpts: func(t *testing.T) options {
				tmpDir := t.TempDir()
				etcDir := filepath.Join(tmpDir, "etc")
				serviceDir := filepath.Join(tmpDir, "services")
				os.MkdirAll(etcDir, 0o755)
				os.MkdirAll(serviceDir, 0o755)
				return options{etc: etcDir, serviceDir: serviceDir}
			},
			expectErr: true,
		},
	} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			opts := tc.setupOpts(t)
			ctx := context.Background()
			err := defaults(ctx, "testscript", opts)

			if tc.expectErr && err == nil {
				t.Fatalf(
					"defaults(ctx, \"testscript\", %v) = nil, wanted error",
					opts,
				)
			} else if !tc.expectErr && err != nil {
				t.Fatalf(
					"defaults(ctx, \"testscript\", %v) = %v, wanted nil",
					opts, err,
				)
			}
		})
	}
}

func TestDefaultsDisable(t *testing.T) {
	t.Parallel()

	for name, tc := range map[string]struct {
		setupOpts func(t *testing.T) options
		expectErr bool
	}{
		"valid script disables symlinks": {
			setupOpts: func(t *testing.T) options {
				tmpDir := t.TempDir()
				etcDir := filepath.Join(tmpDir, "etc")
				serviceDir := filepath.Join(tmpDir, "services")
				os.MkdirAll(etcDir, 0o755)
				os.MkdirAll(serviceDir, 0o755)
				for i := 2; i <= 5; i++ {
					os.MkdirAll(filepath.Join(etcDir, fmt.Sprintf("rc%d.d", i)), 0o755)
				}
				createSampleScript(t, serviceDir, "testscript")
				return options{etc: etcDir, serviceDir: serviceDir}
			},
			expectErr: false,
		},
	} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			opts := tc.setupOpts(t)
			ctx := context.Background()
			err := defaultsDisable(ctx, "testscript", opts)

			if tc.expectErr && err == nil {
				t.Fatalf(
					"defaultsDisable(ctx, \"testscript\", %v) = nil, wanted error",
					opts,
				)
			} else if !tc.expectErr && err != nil {
				t.Fatalf(
					"defaultsDisable(ctx, \"testscript\", %v) = %v, wanted nil",
					opts, err,
				)
			}
		})
	}
}

func TestDisable(t *testing.T) {
	t.Parallel()

	for name, tc := range map[string]struct {
		setupOpts func(t *testing.T) options
		expectErr bool
	}{
		"rename start links to stop links": {
			setupOpts: func(t *testing.T) options {
				tmpDir := t.TempDir()
				etcDir := filepath.Join(tmpDir, "etc")
				serviceDir := filepath.Join(tmpDir, "services")
				os.MkdirAll(etcDir, 0o755)
				os.MkdirAll(serviceDir, 0o755)
				for i := 0; i <= 6; i++ {
					dir := filepath.Join(etcDir, fmt.Sprintf("rc%d.d", i))
					os.MkdirAll(dir, 0o755)
					os.Symlink(
						filepath.Join(serviceDir, "testscript"),
						filepath.Join(dir, "S30testscript"),
					)
				}
				createSampleScript(t, serviceDir, "testscript")
				return options{etc: etcDir, serviceDir: serviceDir}
			},
			expectErr: false,
		},
	} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			opts := tc.setupOpts(t)
			ctx := context.Background()
			err := disable(ctx, "testscript", opts)

			if tc.expectErr && err == nil {
				t.Fatalf(
					"disable(ctx, \"testscript\", %v) = nil, wanted error",
					opts,
				)
			} else if !tc.expectErr && err != nil {
				t.Fatalf(
					"disable(ctx, \"testscript\", %v) = %v, wanted nil",
					opts, err,
				)
			}
		})
	}
}

func TestEnable(t *testing.T) {
	t.Parallel()

	for name, tc := range map[string]struct {
		setupOpts func(t *testing.T) options
		expectErr bool
	}{
		"rename stop links to start links": {
			setupOpts: func(t *testing.T) options {
				tmpDir := t.TempDir()
				etcDir := filepath.Join(tmpDir, "etc")
				serviceDir := filepath.Join(tmpDir, "services")
				os.MkdirAll(etcDir, 0o755)
				os.MkdirAll(serviceDir, 0o755)
				for i := 0; i <= 6; i++ {
					dir := filepath.Join(etcDir, fmt.Sprintf("rc%d.d", i))
					os.MkdirAll(dir, 0o755)
					os.Symlink(
						filepath.Join(serviceDir, "testscript"),
						filepath.Join(dir, "K70testscript"),
					)
				}
				createSampleScript(t, serviceDir, "testscript")
				return options{etc: etcDir, serviceDir: serviceDir}
			},
			expectErr: false,
		},
	} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			opts := tc.setupOpts(t)
			ctx := context.Background()
			err := enable(ctx, "testscript", opts)

			if tc.expectErr && err == nil {
				t.Fatalf(
					"enable(ctx, \"testscript\", %v) = nil, wanted error",
					opts,
				)
			} else if !tc.expectErr && err != nil {
				t.Fatalf(
					"enable(ctx, \"testscript\", %v) = %v, wanted nil",
					opts, err,
				)
			}
		})
	}
}

func TestRemove(t *testing.T) {
	t.Parallel()

	for name, tc := range map[string]struct {
		setupOpts func(t *testing.T) options
		expectErr bool
	}{
		"remove existing symlinks": {
			setupOpts: func(t *testing.T) options {
				tmpDir := t.TempDir()
				etcDir := filepath.Join(tmpDir, "etc")
				serviceDir := filepath.Join(tmpDir, "services")
				os.MkdirAll(etcDir, 0o755)
				os.MkdirAll(serviceDir, 0o755)
				for i := 0; i <= 6; i++ {
					dir := filepath.Join(etcDir, fmt.Sprintf("rc%d.d", i))
					os.MkdirAll(dir, 0o755)
					os.Symlink(
						filepath.Join(serviceDir, "testscript"),
						filepath.Join(dir, "S30testscript"),
					)
					os.Symlink(
						filepath.Join(serviceDir, "testscript"),
						filepath.Join(dir, "K70testscript"),
					)
				}
				createSampleScript(t, serviceDir, "testscript")
				return options{etc: etcDir, serviceDir: serviceDir, force: true}
			},
			expectErr: false,
		},
	} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			opts := tc.setupOpts(t)
			ctx := context.Background()
			err := remove(ctx, "testscript", opts)

			if tc.expectErr && err == nil {
				t.Fatalf(
					"remove(ctx, \"testscript\", %v) = nil, wanted error",
					opts,
				)
			} else if !tc.expectErr && err != nil {
				t.Fatalf(
					"remove(ctx, \"testscript\", %v) = %v, wanted nil",
					opts, err,
				)
			}
		})
	}
}
