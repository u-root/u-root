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
			name:    "unique no duplicates",
			params:  params{unique: true},
			input:   "a\nb\nc\n",
			want:    "a\nb\nc\n",
			wantErr: nil,
		},
		{
			name:    "unique with duplicates",
			params:  params{unique: true},
			input:   "a\nb\nc\na\n",
			want:    "a\nb\nc\n",
			wantErr: nil,
		},
		{
			name:    "unique and ordered no duplicates",
			params:  params{unique: true, ordered: true},
			input:   "a\nb\nc\n",
			wantErr: nil,
		},
		{
			name:    "unique and ordered with duplicates",
			params:  params{unique: true, ordered: true},
			input:   "a\nb\nc\na\n",
			wantErr: errNotOrdered,
		},
		{
			name:  "ignore case off",
			input: "apple\nOrange\n",
			want:  "Orange\napple\n",
		},
		{
			name:   "ignore case on 1",
			params: params{ignoreCase: true},
			input:  "apple\nOrange\n",
			want:   "apple\nOrange\n",
		},
		{
			name:   "ignore case on 2",
			params: params{ignoreCase: true},
			input:  "apple\nOrange\napple\n",
			want:   "apple\napple\nOrange\n",
		},
		{
			name:   "ordered if ignore case is true",
			params: params{ignoreCase: true, ordered: true},
			input:  "a\nB\nc\nD\ne\n",
		},
		{
			name:    "unique with ignore case",
			params:  params{unique: true, ignoreCase: true},
			input:   "a\nA\n",
			want:    "a\n",
			wantErr: nil,
		},
		{
			name:    "ordered but not unique",
			params:  params{ordered: true, unique: true},
			input:   "a\na\n",
			wantErr: errNotOrdered,
		},
		{
			name:    "unique and ignore case not ordered",
			params:  params{unique: true, ignoreCase: true, ordered: true},
			input:   "A\na\n",
			wantErr: errNotOrdered,
		},
		{
			name:   "ignore blanks",
			params: params{ignoreBlanks: true},
			input:  "  b\nA\n",
			want:   "A\n  b\n",
		},
		{
			name:   "ignore blanks and ignore case",
			params: params{ignoreCase: true, ignoreBlanks: true},
			input:  "  b\nA\n  C\n",
			want:   "A\n  b\n  C\n",
		},
		{
			name:   "ignore blanks, case and unique",
			params: params{ignoreCase: true, ignoreBlanks: true, unique: true},
			input:  " b\nA\n C\nA\nb\n",
			want:   "A\n b\n C\n",
		},
		{
			name:    "ignore blanks (no effect), case and unique and ordered ",
			params:  params{ignoreCase: true, ignoreBlanks: true, ordered: true},
			input:   " b\nA\n C\nA\nb\n",
			wantErr: errNotOrdered,
		},
		{
			name:   "ignore blanks breaking ties",
			params: params{ignoreBlanks: true},
			input:  " {\n  {\n {\n  {\n",
			want:   "  {\n  {\n {\n {\n",
		},
		{
			name:   "ignore blanks breaking ties with and ignore case",
			params: params{ignoreBlanks: true, ignoreCase: true},
			input:  "a\n {\n  {\n {\n  {\nA\n",
			want:   "A\na\n  {\n  {\n {\n {\n",
		},
		{
			name:   "numeric sort",
			params: params{numeric: true},
			input:  "39\n100\n21\n",
			want:   "21\n39\n100\n",
		},
		{
			name:   "numeric sort with floats",
			params: params{numeric: true},
			input:  "39.1\n100.2\n-21.3\n",
			want:   "-21.3\n39.1\n100.2\n",
		},
		{
			name:   "numeric sort with blanks",
			params: params{numeric: true},
			input:  " 39\n  100\n21\n",
			want:   "21\n 39\n  100\n",
		},
		{
			name:   "numeric sort with leading zeros",
			params: params{numeric: true},
			input:  "039\n00100.1\n00021\n",
			want:   "00021\n039\n00100.1\n",
		},
		{
			name:   "numeric sort with non numeric strings",
			params: params{numeric: true},
			input:  "hello\n1.0\n-1.0\n",
			want:   "-1.0\nhello\n1.0\n",
		},
		{
			name:   "numeric sort ordered",
			params: params{numeric: true, ordered: true},
			input:  "-1\n2.1\n3\n",
		},
		{
			name:    "numeric sort unordered",
			params:  params{numeric: true, ordered: true},
			input:   "01\n2.1\n0.2\n",
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
