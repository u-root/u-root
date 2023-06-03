// Copyright 2016-2023 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"testing"
)

func TestSortStdin(t *testing.T) {
	for _, tt := range []struct {
		name    string
		params  params
		input   string
		want    string
		wantErr error
	}{
		{
			name:    "unique no dublicates",
			params:  params{unique: true},
			input:   "a\nb\nc\n",
			want:    "a\nb\nc\n",
			wantErr: nil,
		},
		{
			name:    "unique with dublicates",
			params:  params{unique: true},
			input:   "a\nb\nc\na\n",
			want:    "a\nb\nc\n",
			wantErr: nil,
		},
		{
			name:    "unique and ordered no dublicates",
			params:  params{unique: true, ordered: true},
			input:   "a\nb\nc\n",
			wantErr: nil,
		},
		{
			name:    "unique and ordered with dublicates",
			params:  params{unique: true, ordered: true},
			input:   "a\nb\nc\na\n",
			wantErr: errNotOrdered,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			stdin := io.NopCloser(strings.NewReader(tt.input))
			stdout := &bytes.Buffer{}

			err := command(stdin, stdout, nil, tt.params, nil).run()
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("sort err: = %q, want: %q", err, tt.wantErr)
			}

			if stdout.String() != tt.want {
				t.Errorf("sort = %q, want: %q", stdout.String(), tt.want)
			}
		})
	}
}

func TestSortFiles(t *testing.T) {
	tmpDir := t.TempDir()
	file1, err := os.Create(filepath.Join(tmpDir, "file1"))
	if err != nil {
		t.Fatalf("Failed to create tmp file1: %v", err)
	}
	if _, err := file1.WriteString("α\nβ\nγ"); err != nil {
		t.Fatalf("failed to write into file1: %v", err)
	}
	file2, err := os.Create(filepath.Join(tmpDir, "file2"))
	if err != nil {
		t.Fatalf("Failed to create tmp file2: %v", err)
	}
	if _, err := file2.WriteString("a\nd\nc\n"); err != nil {
		t.Fatalf("failed to write into file1: %v", err)
	}

	for _, tt := range []struct {
		name    string
		args    []string
		params  params
		want    string
		wantErr error
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
			name:   "reversed = true",
			args:   []string{filepath.Join(tmpDir, "file1"), filepath.Join(tmpDir, "file2")},
			params: params{reverse: true},
			want:   "γ\nβ\nα\nd\nc\na\n",
		},
		{
			name:   "outputfile set",
			args:   []string{filepath.Join(tmpDir, "file1"), filepath.Join(tmpDir, "file2")},
			params: params{outputFile: filepath.Join(tmpDir, "outputfile")},
			want:   "a\nc\nd\nα\nβ\nγ\n",
		},
		{
			name:    "no such file or directory",
			args:    []string{"nosuchfile"},
			wantErr: syscall.Errno(2),
		},
		{
			name:   "ordered",
			args:   []string{filepath.Join(tmpDir, "file1")},
			params: params{ordered: true},
		},
		{
			name:    "not ordered",
			args:    []string{filepath.Join(tmpDir, "file2")},
			params:  params{ordered: true},
			wantErr: errNotOrdered,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			f, err := os.Create(filepath.Join(tmpDir, "file"))
			if err != nil {
				t.Fatalf("failed to create tmp file: %v", err)
			}
			stdout := &bytes.Buffer{}
			c := command(f, stdout, nil, tt.params, tt.args)

			err = c.run()
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("sort err: = %q, want: %q", err, tt.wantErr)
			}

			if tt.params.outputFile != "" {
				res, err := os.ReadFile(tt.params.outputFile)
				if err != nil {
					t.Fatalf("failed to read file: %v", err)
				}
				if string(res) != tt.want {
					t.Errorf("sort = %q, want: %q", string(res), tt.want)
				}
			} else {
				if stdout.String() != tt.want {
					t.Errorf("sort = %q, want: %q", stdout.String(), tt.want)
				}
			}
		})
	}
}
