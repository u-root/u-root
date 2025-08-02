// Copyright 2012 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rm

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
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
			want: usage,
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
			cmd := New()
			var stdout, stderr bytes.Buffer
			var stdin bytes.Buffer
			stdin.WriteString(tt.iString)

			cmd.SetIO(&stdin, &stdout, &stderr)
			cmd.SetWorkingDir(d)

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

func TestRmWorkingDir(t *testing.T) {
	d := setup(t)

	// Test that working directory is respected
	cmd := New()
	var stdout, stderr bytes.Buffer
	var stdin bytes.Buffer

	cmd.SetIO(&stdin, &stdout, &stderr)
	cmd.SetWorkingDir(d)

	// Remove file using relative path
	err := cmd.Run("go.txt")
	if err != nil {
		t.Errorf("Run() = %v, want nil", err)
	}

	// Verify file was removed
	if _, err := os.Stat(filepath.Join(d, "go.txt")); !os.IsNotExist(err) {
		t.Errorf("File should have been removed")
	}
}

func TestRmInteractive(t *testing.T) {
	d := setup(t)

	// Test interactive mode with "no" response
	cmd := New()
	var stdout, stderr bytes.Buffer
	var stdin bytes.Buffer
	stdin.WriteString("n\n")

	cmd.SetIO(&stdin, &stdout, &stderr)

	err := cmd.Run("-i", filepath.Join(d, "go.txt"))
	if err != nil {
		t.Errorf("Run() = %v, want nil", err)
	}

	// Verify file was NOT removed
	if _, err := os.Stat(filepath.Join(d, "go.txt")); os.IsNotExist(err) {
		t.Errorf("File should not have been removed")
	}

	// Test interactive mode with "yes" response
	cmd2 := New()
	var stdout2, stderr2 bytes.Buffer
	var stdin2 bytes.Buffer
	stdin2.WriteString("y\n")

	cmd2.SetIO(&stdin2, &stdout2, &stderr2)

	err = cmd2.Run("-i", filepath.Join(d, "go.txt"))
	if err != nil {
		t.Errorf("Run() = %v, want nil", err)
	}

	// Verify file was removed
	if _, err := os.Stat(filepath.Join(d, "go.txt")); !os.IsNotExist(err) {
		t.Errorf("File should have been removed")
	}
}
