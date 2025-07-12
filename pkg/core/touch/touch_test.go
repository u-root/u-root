// Copyright 2023 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package touch

import (
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestParseParamsDate(t *testing.T) {
	cmd := New().(*command)
	date := "2021-01-01T00:00:00Z"
	expected, err := time.Parse(time.RFC3339, date)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	p, err := cmd.parseParams(date, false, false, false)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if !expected.Equal(p.time) {
		t.Errorf("expected %v, got %v", expected, p.time)
	}

	date = "invalid"
	_, err = cmd.parseParams(date, false, false, false)
	if err == nil {
		t.Errorf("expected error, got nil")
	}
}

func TestParseParams(t *testing.T) {
	cmd := New().(*command)
	tests := []struct {
		expected     params
		access       bool
		modification bool
		create       bool
	}{
		{
			access:       false,
			modification: false,
			create:       false,
			expected: params{
				access:       true,
				modification: true,
				create:       false,
			},
		},
		{
			access:       true,
			modification: false,
			create:       false,
			expected: params{
				access:       true,
				modification: false,
				create:       false,
			},
		},
		{
			access:       false,
			modification: true,
			create:       true,
			expected: params{
				access:       false,
				modification: true,
				create:       true,
			},
		},
	}

	for _, test := range tests {
		p, err := cmd.parseParams("", test.access, test.modification, test.create)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if p.access != test.expected.access {
			t.Errorf("expected %v, got %v", test.expected.access, p.access)
		}
		if p.modification != test.expected.modification {
			t.Errorf("expected %v, got %v", test.expected.modification, p.modification)
		}
		if p.create != test.expected.create {
			t.Errorf("expected %v, got %v", test.expected.create, p.create)
		}
	}
}

var tests = []struct {
	err  error
	p    params
	name string
	args []string
}{
	{
		name: "create is true, no new files created",
		args: []string{"a1", "a2"},
		p: params{
			access:       true,
			modification: true,
			create:       true,
			time:         time.Now(),
		},
	},
	{
		name: "create is false, files should be created",
		args: []string{"a1", "a2"},
		p: params{
			access:       true,
			modification: true,
			create:       false,
			time:         time.Now(),
		},
	},
	{
		name: "no such file or directory",
		args: []string{"no/such/file/or/direcotry"},
		p: params{
			create: false,
			time:   time.Now(),
		},
		err: os.ErrNotExist,
	},
}

func TestTouchEmptyDir(t *testing.T) {
	for _, test := range tests {
		temp := t.TempDir()
		var args []string
		for _, arg := range test.args {
			args = append(args, filepath.Join(temp, arg))
		}

		cmd := New().(*command)
		err := cmd.touchFiles(test.p, args)
		if !errors.Is(err, test.err) {
			t.Fatalf("touchFiles() expected %v, got %v", test.err, err)
		}
		if test.err != nil {
			continue
		}

		for _, arg := range args {
			_, err := os.Stat(arg)
			if test.p.create {
				if !os.IsNotExist(err) {
					t.Errorf("expected %s to not exist", arg)
				}
			} else {
				if err != nil {
					t.Errorf("expected %s to exist, got %v", arg, err)
				}

				stat, err := os.Stat(arg)
				if err != nil {
					t.Fatalf("expected no error, got %v", err)
				}

				if test.p.modification {
					if stat.ModTime().Unix() != test.p.time.Unix() {
						t.Errorf("expected %s to have mod time %v, got %v", arg, test.p.time, stat.ModTime())
					}
				}
			}
		}
	}
}

func TestTouchCommand(t *testing.T) {
	tempDir := t.TempDir()
	testFile := filepath.Join(tempDir, "testfile.txt")

	cmd := New()
	var stdout, stderr bytes.Buffer
	cmd.SetIO(bytes.NewReader(nil), &stdout, &stderr)

	// Test creating a new file
	err := cmd.Run(testFile)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Verify file was created
	if _, err := os.Stat(testFile); os.IsNotExist(err) {
		t.Errorf("Expected file to be created")
	}
}

func TestTouchCommandWithDate(t *testing.T) {
	tempDir := t.TempDir()
	testFile := filepath.Join(tempDir, "testfile.txt")

	cmd := New()
	var stdout, stderr bytes.Buffer
	cmd.SetIO(bytes.NewReader(nil), &stdout, &stderr)

	// Test with specific date
	err := cmd.Run("-d", "2021-01-01T00:00:00Z", testFile)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Verify file was created with correct time
	stat, err := os.Stat(testFile)
	if err != nil {
		t.Fatalf("Expected file to exist, got %v", err)
	}

	expectedTime, _ := time.Parse(time.RFC3339, "2021-01-01T00:00:00Z")
	if stat.ModTime().Unix() != expectedTime.Unix() {
		t.Errorf("Expected mod time %v, got %v", expectedTime, stat.ModTime())
	}
}

func TestTouchCommandNoCreate(t *testing.T) {
	tempDir := t.TempDir()
	testFile := filepath.Join(tempDir, "nonexistent.txt")

	cmd := New()
	var stdout, stderr bytes.Buffer
	cmd.SetIO(bytes.NewReader(nil), &stdout, &stderr)

	// Test with -c flag (don't create)
	err := cmd.Run("-c", testFile)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Verify file was NOT created
	if _, err := os.Stat(testFile); !os.IsNotExist(err) {
		t.Errorf("Expected file to not be created")
	}
}

func TestTouchCommandNoArgs(t *testing.T) {
	cmd := New()
	var stdout, stderr bytes.Buffer
	cmd.SetIO(bytes.NewReader(nil), &stdout, &stderr)

	// Test with no arguments
	err := cmd.Run()
	if err == nil {
		t.Error("Expected error for no arguments")
	}
}

func TestTouchWorkingDir(t *testing.T) {
	tempDir := t.TempDir()
	testFile := "relative_test.txt"

	cmd := New()
	var stdout, stderr bytes.Buffer
	cmd.SetIO(bytes.NewReader(nil), &stdout, &stderr)
	cmd.SetWorkingDir(tempDir)

	// Test with relative path
	err := cmd.Run(testFile)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Verify file was created in the working directory
	fullPath := filepath.Join(tempDir, testFile)
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		t.Errorf("Expected file to be created in working directory")
	}
}
