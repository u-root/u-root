// Copyright 2023 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/u-root/u-root/pkg/core/touch"
)

func TestParseParamsDate(t *testing.T) {
	cmd := touch.New()
	var stdout, stderr bytes.Buffer
	cmd.SetIO(bytes.NewReader(nil), &stdout, &stderr)

	// Test valid date
	err := cmd.Run("-d", "2021-01-01T00:00:00Z", "/tmp/test_touch_date")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Clean up
	os.Remove("/tmp/test_touch_date")

	// Test invalid date
	cmd2 := touch.New()
	var stdout2, stderr2 bytes.Buffer
	var stdin2 bytes.Buffer
	cmd2.SetIO(&stdin2, &stdout2, &stderr2)

	err = cmd2.Run("-d", "invalid", "/tmp/test_touch_invalid")
	if err == nil {
		t.Error("expected error for invalid date, got nil")
	}
}

var tests = []struct {
	err  error
	name string
	args []string
}{
	{
		name: "create is true, no new files created",
		args: []string{"-c", "a1", "a2"},
	},
	{
		name: "create is false, files should be created",
		args: []string{"a1", "a2"},
	},
	{
		name: "no such file or directory",
		args: []string{"no/such/file/or/direcotry"},
		err:  os.ErrNotExist,
	},
}

func TestTouchEmptyDir(t *testing.T) {
	for _, test := range tests {
		temp := t.TempDir()
		var args []string
		args = append(args, test.args[0]) // "touch"
		for i := 1; i < len(test.args); i++ {
			arg := test.args[i]
			if !strings.HasPrefix(arg, "-") {
				args = append(args, filepath.Join(temp, arg))
			} else {
				args = append(args, arg)
			}
		}

		cmd := touch.New()
		var stdout, stderr bytes.Buffer
		cmd.SetIO(bytes.NewReader(nil), &stdout, &stderr)

		err := cmd.Run(args...)
		if test.err != nil {
			if !errors.Is(err, test.err) {
				t.Fatalf("Run() expected %v, got %v", test.err, err)
			}
			continue
		}

		if err != nil {
			t.Fatalf("Run() expected no error, got %v", err)
		}

		// Check if files were created (only for non-error cases)
		for i := 1; i < len(test.args); i++ {
			arg := test.args[i]
			if !strings.HasPrefix(arg, "-") {
				fullPath := filepath.Join(temp, arg)
				_, err := os.Stat(fullPath)
				if strings.Contains(test.name, "create is true") {
					// With -c flag, files should not be created if they don't exist
					if !os.IsNotExist(err) {
						t.Errorf("expected %s to not exist", fullPath)
					}
				} else {
					// Without -c flag, files should be created
					if err != nil {
						t.Errorf("expected %s to exist, got %v", fullPath, err)
					}
				}
			}
		}
	}
}
