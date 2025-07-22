// Copyright 2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mktemp

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestMkTemp(t *testing.T) {
	tmpDir := os.TempDir()
	tests := []struct {
		name    string
		args    []string
		wantOut string
		wantErr bool
	}{
		{
			name:    "basic mktemp",
			args:    []string{},
			wantOut: tmpDir,
			wantErr: false,
		},
		{
			name:    "directory mode",
			args:    []string{"-d"},
			wantOut: tmpDir,
			wantErr: false,
		},
		{
			name:    "with template",
			args:    []string{"foofoo.XXXX"},
			wantOut: filepath.Join(tmpDir, "foofoo"),
			wantErr: false,
		},
		{
			name:    "with suffix",
			args:    []string{"--suffix", "baz", "foo.XXXX"},
			wantOut: filepath.Join(tmpDir, "foo.baz"),
			wantErr: false,
		},
		{
			name:    "dry run",
			args:    []string{"-u"},
			wantOut: "",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var stdout bytes.Buffer
			var stderr bytes.Buffer

			cmd := New()
			cmd.SetIO(nil, &stdout, &stderr)

			err := cmd.Run(tt.args...)
			if (err != nil) != tt.wantErr {
				t.Errorf("Run() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			output := stdout.String()
			if !strings.HasPrefix(output, tt.wantOut) {
				t.Errorf("stdout got:\n%s\nwant starting with:\n%s", output, tt.wantOut)
			}
		})
	}
}

func TestMkTempFlags(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		checkDir bool
	}{
		{
			name:     "directory flag",
			args:     []string{"-d"},
			checkDir: true,
		},
		{
			name:     "file creation",
			args:     []string{},
			checkDir: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var stdout bytes.Buffer
			var stderr bytes.Buffer

			cmd := New()
			cmd.SetIO(nil, &stdout, &stderr)

			err := cmd.Run(tt.args...)
			if err != nil {
				t.Fatalf("Run() error = %v", err)
			}

			output := strings.TrimSpace(stdout.String())
			if output == "" {
				return // dry-run case
			}

			info, err := os.Stat(output)
			if err != nil {
				t.Fatalf("created path %s does not exist: %v", output, err)
			}

			if tt.checkDir && !info.IsDir() {
				t.Errorf("expected directory, got file")
			}
			if !tt.checkDir && info.IsDir() {
				t.Errorf("expected file, got directory")
			}

			// Clean up
			if tt.checkDir {
				os.RemoveAll(output)
			} else {
				os.Remove(output)
			}
		})
	}
}
