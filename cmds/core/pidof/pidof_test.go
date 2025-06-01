// Copyright 2025 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"errors"
	"testing"
)

func TestPidof(t *testing.T) {
	tests := []struct {
		name     string
		procDir  string
		expected string
		args     []string
		single   bool
	}{
		{
			name:     "multiple processes",
			procDir:  "./testdata",
			args:     []string{"init", "bash"},
			expected: "1 2\n",
		},
		{
			name:     "multiple pids with single flag",
			procDir:  "./testdata",
			single:   true,
			args:     []string{"process"},
			expected: "3\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stdout := bytes.Buffer{}

			err := run(&stdout, tt.procDir, tt.single, tt.args)
			if err != nil {
				t.Fatalf("expected nil got %v", err)
			}

			if stdout.String() != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, stdout.String())
			}
		})
	}
}

func TestPidofMissing(t *testing.T) {
	stdout := bytes.Buffer{}
	err := run(&stdout, "./testdata", false, []string{"goooo"})
	if !errors.Is(err, errNotFound) {
		t.Fatalf("expected %v got %v", errNotFound, err)
	}

	if stdout.String() != "" {
		t.Errorf("expected empty string got %q", stdout.String())
	}
}
