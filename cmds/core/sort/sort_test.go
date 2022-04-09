// Copyright 2016 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

func TestSort(t *testing.T) {
	tmpDir := t.TempDir()
	file1, err := os.Create(filepath.Join(tmpDir, "file1"))
	if err != nil {
		t.Errorf("Failed to create tmp file1: %v", err)
	}
	if _, err := file1.WriteString("α\nβ\nγ"); err != nil {
		t.Errorf("failed to write into file1: %v", err)
	}
	file2, err := os.Create(filepath.Join(tmpDir, "file2"))
	if err != nil {
		t.Errorf("Failed to create tmp file2: %v", err)
	}
	if _, err := file2.WriteString("a\nd\nc\n"); err != nil {
		t.Errorf("failed to write into file1: %v", err)
	}
	for _, tt := range []struct {
		name       string
		args       []string
		reverse    bool
		outputFile string
		want       string
		wantErr    string
	}{
		{
			name: "empty input",
			args: []string{},
			want: "",
		},
		{
			name: "input from 2 files",
			args: []string{filepath.Join(tmpDir, "file1"), filepath.Join(tmpDir, "file2")},
			want: "a\nc\nd\nα\nβ\nγ\n",
		},
		{
			name:    "reversed = true",
			args:    []string{filepath.Join(tmpDir, "file1"), filepath.Join(tmpDir, "file2")},
			reverse: true,
			want:    "γ\nβ\nα\nd\nc\na\n",
		},
		{
			name:       "outputfile set",
			args:       []string{filepath.Join(tmpDir, "file1"), filepath.Join(tmpDir, "file2")},
			outputFile: filepath.Join(tmpDir, "outputfile"),
			want:       "a\nc\nd\nα\nβ\nγ\n",
		},
	} {
		*reverse = tt.reverse
		*outputFile = tt.outputFile
		t.Run(tt.name, func(t *testing.T) {
			buf := &bytes.Buffer{}
			f, err := os.Create(filepath.Join(tmpDir, "file"))
			if err != nil {
				t.Errorf("failed to create tmp file: %v", err)
			}
			if got := readInput(buf, f, tt.args...); got != nil {
				if got.Error() != tt.wantErr {
					t.Errorf("readInput() = %q, want: %q", got.Error(), tt.wantErr)
				}
			} else if tt.name == "outputfile set" {
				sort, err := os.ReadFile(filepath.Join(tmpDir, "outputfile"))
				if err != nil {
					t.Errorf("Failed to read file: %v", err)
				}
				if string(sort) != tt.want {
					t.Errorf("readInput() = %q, want: %q", string(sort), tt.want)
				}
			} else {
				if buf.String() != tt.want {
					t.Errorf("readInput() = %q, want: %q", buf.String(), tt.want)
				}
			}
		})
	}
}
