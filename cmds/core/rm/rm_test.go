// Copyright 2012 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/u-root/u-root/pkg/core/rm"
)

func setup(t *testing.T) string {
	d := t.TempDir()
	fbody := []byte("Go is cool!")
	for _, f := range []struct {
		name  string
		mode  os.FileMode
		isdir bool
	}{
		{
			name:  "hi",
			mode:  0o755,
			isdir: true,
		},
		{
			name: "hi/one.txt",
			mode: 0o666,
		},
		{
			name: "hi/two.txt",
			mode: 0o777,
		},
		{
			name: "go.txt",
			mode: 0o555,
		},
	} {
		var (
			err      error
			filepath = filepath.Join(d, f.name)
		)
		if f.isdir {
			err = os.Mkdir(filepath, f.mode)
		} else {
			err = os.WriteFile(filepath, fbody, f.mode)
		}
		if err != nil {
			t.Fatal(err)
		}
	}
	return d
}

func TestRm(t *testing.T) {
	for _, tt := range []struct {
		name        string
		args        []string
		interactive bool
		iString     string
		verbose     bool
		recursive   bool
		force       bool
		want        string
	}{
		{
			name: "no args",
			args: nil,
			want: "rm [-Rrvif] file...",
		},
		{
			name: "rm one file",
			args: []string{"go.txt"},
			want: "",
		},
		{
			name:    "rm one file verbose",
			args:    []string{"-v", "go.txt"},
			verbose: true,
			want:    "",
		},
		{
			name: "fail to rm one file",
			args: []string{"go"},
			want: "no such file or directory",
		},
		{
			name:  "fail to rm one file forced to trigger continue",
			args:  []string{"-f", "go"},
			force: true,
			want:  "",
		},
		{
			name:        "rm one file interactive",
			args:        []string{"-i", "go.txt"},
			interactive: true,
			iString:     "y\n",
			want:        "",
		},
		{
			name:        "rm one file interactive continue triggered",
			args:        []string{"-i", "go.txt"},
			interactive: true,
			iString:     "\n",
			want:        "",
		},
		{
			name:      "rm dir recursively",
			args:      []string{"-r", "hi"},
			recursive: true,
		},
		{
			name: "rm dir not recursively",
			args: []string{"hi"},
			want: "directory not empty",
		},
	} {
		d := setup(t)

		t.Run(tt.name, func(t *testing.T) {
			cmd := rm.New()
			var stdout, stderr bytes.Buffer
			cmd.SetIO(strings.NewReader(tt.iString), &stdout, &stderr)

			// Update args to use absolute paths for files
			args := make([]string, len(tt.args))
			copy(args, tt.args)
			for i := range args {
				if !strings.HasPrefix(args[i], "-") {
					args[i] = filepath.Join(d, args[i])
				}
			}

			err := cmd.Run(args...)

			if tt.want != "" {
				if err == nil || !strings.Contains(err.Error(), tt.want) {
					t.Errorf("Run() = %v, want error containing: %q", err, tt.want)
				}
				return
			}

			if err != nil {
				t.Errorf("Run() = %v, want nil", err)
			}

			// Check verbose output
			if tt.verbose && stdout.Len() == 0 {
				t.Errorf("Expected verbose output, got none")
			}
		})
	}
}
