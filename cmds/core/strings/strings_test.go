// Copyright 2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

func TestStrings(t *testing.T) {
	tmpdir := t.TempDir()
	for _, tt := range []struct {
		name    string
		files   []string
		n       int
		input   string
		want    string
		wantErr string
	}{
		{
			name:  "empty",
			files: []string{},
			input: "",
			n:     4,
			want:  "",
		},
		{
			name:    "n < 1",
			files:   []string{},
			input:   "",
			n:       0,
			wantErr: fmt.Sprintf("strings: invalid minimum string length %v", 0),
		},
		{
			name:  "string printable",
			files: []string{},
			input: "abcdefg",
			n:     4,
			want:  "abcdefg\n",
		},
		{
			name:  "test prevent buffer from growing indefinitely",
			files: []string{},
			input: "qwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwerty",
			n:     4,
			want:  "qwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwerty\n",
		},
		{
			name:  "non ascii char",
			files: []string{},
			input: "abcdefg£",
			n:     4,
			want:  "abcdefg\n",
		},
		{
			name:  "with file",
			files: []string{filepath.Join(tmpdir, "file")},
			input: "abcdefg",
			n:     4,
			want:  "abcdefg\n",
		},
		{
			name:    "with file",
			files:   []string{"filedoesnotexist"},
			input:   "abcdefg",
			n:       4,
			wantErr: "open filedoesnotexist: no such file or directory",
		},
		{
			name:  "strings are too short",
			files: []string{filepath.Join(tmpdir, "file")},
			input: "\n\na\nab\n\nabc\nabc\xff\n01\n",
			n:     4,
			want:  "",
		},
		{
			name:  "string fits perfectly",
			files: []string{filepath.Join(tmpdir, "file")},
			input: "abcd\n",
			n:     4,
			want:  "abcd\n",
		},
		{
			name:  "terminating newline",
			files: []string{filepath.Join(tmpdir, "file")},
			input: "abcdefghijklmnopqrstuvwxyz\n",
			n:     4,
			want:  "abcdefghijklmnopqrstuvwxyz\n",
		},
		{
			name:  "mix of printable and non-printable sequences",
			files: []string{filepath.Join(tmpdir, "file")},
			input: "\n\na123456\nab\n\nabc\nabcde\xff\n01\n",
			n:     4,
			want:  "a123456\nabcde\n",
		},
		{
			name:  "spaces are printable",
			files: []string{filepath.Join(tmpdir, "file")},
			input: " abcdefghijklm nopqrstuvwxyz ",
			n:     4,
			want:  " abcdefghijklm nopqrstuvwxyz \n",
		},
		{
			name:  "shorter value of n",
			files: []string{filepath.Join(tmpdir, "file")},
			input: "\n\na\nab\n\nabc\nabc\xff\n01\n",
			n:     1,
			want:  "a\nab\nabc\nabc\n01\n",
		},
		{
			name:  "larger value of n",
			files: []string{filepath.Join(tmpdir, "file")},
			input: "\n\na123456\nab\n\nabc\nabcde\xff\n01\n",
			n:     6,
			want:  "a123456\n",
		},
	} {
		// Write input into file than seek to the beginning and truncate the file afterwards
		// Create file for input data
		file, err := os.Create(filepath.Join(tmpdir, "file"))
		if err != nil {
			t.Errorf("failed to create tmp file: %v", err)
		}
		_, err = file.WriteString(tt.input)
		if err != nil {
			t.Errorf("failed to write to file: %v", err)
		}

		// Run tests
		t.Run(tt.name, func(t *testing.T) {
			bufIn := &bytes.Buffer{}
			bufIn.WriteString(tt.input)
			bufOut := &bytes.Buffer{}

			if got := command(bufIn, bufOut, params{n: tt.n}, tt.files).run(); got != nil {
				if got.Error() != tt.wantErr {
					t.Errorf("strings() = %q, want: %q", got.Error(), tt.wantErr)
				}
			} else {
				if bufOut.String() != tt.want {
					t.Errorf("strings() = %q, want: %q", bufOut.String(), tt.want)
				}
			}
		})
	}
}
