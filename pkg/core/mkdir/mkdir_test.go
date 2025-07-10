// Copyright 2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mkdir

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"golang.org/x/sys/unix"
)

type testFlags struct {
	mode    string
	mkall   bool
	verbose bool
}

func TestMkdir(t *testing.T) {
	d := t.TempDir()
	for _, tt := range []struct {
		name      string
		flags     testFlags
		args      []string
		wantMode  string
		wantPrint string
		want      error
	}{
		{
			name:     "Create 1 directory",
			flags:    testFlags{mode: "755"},
			args:     []string{filepath.Join(d, "stub0")},
			wantMode: "drwxr-xr-x",
		},
		{
			name:      "Directory already exists",
			flags:     testFlags{mode: "755"},
			args:      []string{filepath.Join(d, "stub0")},
			wantMode:  "drwxr-xr-x",
			wantPrint: fmt.Sprintf("%s: %s file exists", filepath.Join(d, "stub0"), filepath.Join(d, "stub0")),
		},
		{
			name: "Create 1 directory verbose",
			flags: testFlags{
				mode:    "755",
				verbose: true,
			},
			args:     []string{filepath.Join(d, "stub1")},
			wantMode: "drwxr-xr-x",
		},
		{
			name:     "Create 2 directories",
			flags:    testFlags{mode: "755"},
			args:     []string{filepath.Join(d, "stub2"), filepath.Join(d, "stub3")},
			wantMode: "drwxr-xr-x",
		},
		{
			name: "Create a sub directory directly",
			flags: testFlags{
				mode:  "755",
				mkall: true,
			},
			args:     []string{filepath.Join(d, "stub4"), filepath.Join(d, "stub4/subdir")},
			wantMode: "drwxr-xr-x",
		},
		{
			name:  "Perm Mode Bits over 7 Error",
			flags: testFlags{mode: "7778"},
			args:  []string{filepath.Join(d, "stub1")},
			want:  fmt.Errorf(`invalid mode "7778"`),
		},
		{
			name:     "More than 4 Perm Mode Bits Error",
			flags:    testFlags{mode: "11111"},
			args:     []string{filepath.Join(d, "stub1")},
			wantMode: "drwxrwxr-x",
			want:     fmt.Errorf(`invalid mode "11111"`),
		},
		{
			name:     "Custom Perm in Octal Form",
			flags:    testFlags{mode: "0777"},
			args:     []string{filepath.Join(d, "stub6")},
			wantMode: "drwxrwxrwx",
		},
		{
			name:     "Custom Perm not in Octal Form",
			flags:    testFlags{mode: "777"},
			args:     []string{filepath.Join(d, "stub7")},
			wantMode: "drwxrwxrwx",
		},
		{
			name:     "Custom Perm with Sticky Bit",
			flags:    testFlags{mode: "1777"},
			args:     []string{filepath.Join(d, "stub8")},
			wantMode: "dtrwxrwxrwx",
		},
		{
			name:     "Custom Perm with SGID Bit",
			flags:    testFlags{mode: "2777"},
			args:     []string{filepath.Join(d, "stub9")},
			wantMode: "dgrwxrwxrwx",
		},
		{
			name:     "Custom Perm with SUID Bit",
			flags:    testFlags{mode: "4777"},
			args:     []string{filepath.Join(d, "stub10")},
			wantMode: "durwxrwxrwx",
		},
		{
			name:     "Custom Perm with Sticky Bit and SUID Bit",
			flags:    testFlags{mode: "5777"},
			args:     []string{filepath.Join(d, "stub11")},
			wantMode: "dutrwxrwxrwx",
		},
		{
			name:     "Custom Perm for 2 Directories",
			flags:    testFlags{mode: "5777"},
			args:     []string{filepath.Join(d, "stub12"), filepath.Join(d, "stub13")},
			wantMode: "dutrwxrwxrwx",
		},
		{
			name:     "Default creation mode",
			args:     []string{filepath.Join(d, "stub14")},
			wantMode: "drwxr-xr-x",
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			cmd := New().(*command)
			var stdout, stderr bytes.Buffer
			cmd.SetIO(bytes.NewReader(nil), &stdout, &stderr)

			// don't depend on system umask value, if mode is not specified
			if tt.flags.mode == "" {
				m := unix.Umask(unix.S_IWGRP | unix.S_IWOTH)
				defer func() {
					unix.Umask(m)
				}()
			}

			f := flags{
				mode:    tt.flags.mode,
				mkall:   tt.flags.mkall,
				verbose: tt.flags.verbose,
			}

			got := cmd.mkdirFiles(f, tt.args)
			if got != nil {
				if tt.want == nil || got.Error() != tt.want.Error() {
					t.Errorf("mkdirFiles() = '%v', want: '%v'", got, tt.want)
				}
			} else {
				if stderr.String() != "" {
					if !strings.Contains(stderr.String(), "file exist") {
						t.Errorf("Stderr = '%v', want to contain 'file exist'", stderr.String())
					}
				}
				for _, name := range tt.args {
					if stat, err := os.Stat(name); err == nil {
						if stat.Mode().String() != tt.wantMode {
							t.Errorf("Mode = '%v', want: '%v'", stat.Mode().String(), tt.wantMode)
						}
					}
				}
			}
		})
	}
}

func TestMkdirCommand(t *testing.T) {
	tempDir := t.TempDir()
	testDir := filepath.Join(tempDir, "testdir")

	cmd := New()
	var stdout, stderr bytes.Buffer
	cmd.SetIO(bytes.NewReader(nil), &stdout, &stderr)

	// Test creating a new directory
	err := cmd.Run(testDir)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Verify directory was created
	if stat, err := os.Stat(testDir); err != nil {
		t.Errorf("Expected directory to be created, got %v", err)
	} else if !stat.IsDir() {
		t.Errorf("Expected %s to be a directory", testDir)
	}
}

func TestMkdirCommandWithMode(t *testing.T) {
	tempDir := t.TempDir()
	testDir := filepath.Join(tempDir, "testdir_mode")

	cmd := New()
	var stdout, stderr bytes.Buffer
	cmd.SetIO(bytes.NewReader(nil), &stdout, &stderr)

	// Test with specific mode
	err := cmd.Run("-m", "755", testDir)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Verify directory was created with correct mode
	stat, err := os.Stat(testDir)
	if err != nil {
		t.Fatalf("Expected directory to exist, got %v", err)
	}

	expectedMode := "drwxr-xr-x"
	if stat.Mode().String() != expectedMode {
		t.Errorf("Expected mode %s, got %s", expectedMode, stat.Mode().String())
	}
}

func TestMkdirCommandVerbose(t *testing.T) {
	tempDir := t.TempDir()
	testDir := filepath.Join(tempDir, "testdir_verbose")

	cmd := New()
	var stdout, stderr bytes.Buffer
	cmd.SetIO(bytes.NewReader(nil), &stdout, &stderr)

	// Test with verbose flag
	err := cmd.Run("-v", testDir)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Verify verbose output
	if !strings.Contains(stdout.String(), testDir) {
		t.Errorf("Expected verbose output to contain %s, got %s", testDir, stdout.String())
	}
}

func TestMkdirCommandParents(t *testing.T) {
	tempDir := t.TempDir()
	testDir := filepath.Join(tempDir, "parent", "child", "grandchild")

	cmd := New()
	var stdout, stderr bytes.Buffer
	cmd.SetIO(bytes.NewReader(nil), &stdout, &stderr)

	// Test with -p flag (create parents)
	err := cmd.Run("-p", testDir)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Verify directory was created
	if stat, err := os.Stat(testDir); err != nil {
		t.Errorf("Expected directory to be created, got %v", err)
	} else if !stat.IsDir() {
		t.Errorf("Expected %s to be a directory", testDir)
	}
}

func TestMkdirCommandNoArgs(t *testing.T) {
	cmd := New()
	var stdout, stderr bytes.Buffer
	cmd.SetIO(bytes.NewReader(nil), &stdout, &stderr)

	// Test with no arguments
	err := cmd.Run()
	if err == nil {
		t.Error("Expected error for no arguments")
	}
}

func TestMkdirWorkingDir(t *testing.T) {
	tempDir := t.TempDir()
	testDir := "relative_test_dir"

	cmd := New()
	var stdout, stderr bytes.Buffer
	cmd.SetIO(bytes.NewReader(nil), &stdout, &stderr)
	cmd.SetWorkingDir(tempDir)

	// Test with relative path
	err := cmd.Run(testDir)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Verify directory was created in the working directory
	fullPath := filepath.Join(tempDir, testDir)
	if stat, err := os.Stat(fullPath); err != nil {
		t.Errorf("Expected directory to be created in working directory, got %v", err)
	} else if !stat.IsDir() {
		t.Errorf("Expected %s to be a directory", fullPath)
	}
}

func TestMkdirInvalidMode(t *testing.T) {
	tempDir := t.TempDir()
	testDir := filepath.Join(tempDir, "invalid_mode")

	cmd := New()
	var stdout, stderr bytes.Buffer
	cmd.SetIO(bytes.NewReader(nil), &stdout, &stderr)

	// Test with invalid mode
	err := cmd.Run("-m", "invalid", testDir)
	if err == nil {
		t.Error("Expected error for invalid mode")
	}
}
