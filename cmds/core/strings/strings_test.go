// Copyright 2018-2023 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"testing"
)

func TestStrings(t *testing.T) {
	tmpdir := t.TempDir()
	for _, tt := range []struct {
		name    string
		files   []string
		p       params
		input   string
		want    string
		wantErr error
	}{
		{
			name:  "empty",
			files: []string{},
			input: "",
			p:     params{n: 4},
			want:  "",
		},
		{
			name:    "n < 1",
			files:   []string{},
			input:   "",
			p:       params{n: 0},
			wantErr: errInvalidMinLength,
		},
		{
			name:  "string printable",
			files: []string{},
			input: "abcdefg",
			p:     params{n: 4},
			want:  "abcdefg\n",
		},
		{
			name:  "test prevent buffer from growing indefinitely",
			files: []string{},
			input: "qwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwerty",
			p:     params{n: 4},
			want:  "qwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwertyqwerty\n",
		},
		{
			name:  "non ascii char",
			files: []string{},
			input: "abcdefg£",
			p:     params{n: 4},
			want:  "abcdefg\n",
		},
		{
			name:  "with file",
			files: []string{filepath.Join(tmpdir, "file")},
			input: "abcdefg",
			p:     params{n: 4},
			want:  "abcdefg\n",
		},
		{
			name:    "with file",
			files:   []string{"filedoesnotexist"},
			input:   "abcdefg",
			p:       params{n: 4},
			wantErr: fs.ErrNotExist,
		},
		{
			name:  "strings are too short",
			files: []string{filepath.Join(tmpdir, "file")},
			input: "\n\na\nab\n\nabc\nabc\xff\n01\n",
			p:     params{n: 4},
			want:  "",
		},
		{
			name:  "string fits perfectly",
			files: []string{filepath.Join(tmpdir, "file")},
			input: "abcd\n",
			p:     params{n: 4},
			want:  "abcd\n",
		},
		{
			name:  "terminating newline",
			files: []string{filepath.Join(tmpdir, "file")},
			input: "abcdefghijklmnopqrstuvwxyz\n",
			p:     params{n: 4},
			want:  "abcdefghijklmnopqrstuvwxyz\n",
		},
		{
			name:  "mix of printable and non-printable sequences",
			files: []string{filepath.Join(tmpdir, "file")},
			input: "\n\na123456\nab\n\nabc\nabcde\xff\n01\n",
			p:     params{n: 4},
			want:  "a123456\nabcde\n",
		},
		{
			name:  "spaces are printable",
			files: []string{filepath.Join(tmpdir, "file")},
			input: " abcdefghijklm nopqrstuvwxyz ",
			p:     params{n: 4},
			want:  " abcdefghijklm nopqrstuvwxyz \n",
		},
		{
			name:  "shorter value of n",
			files: []string{filepath.Join(tmpdir, "file")},
			input: "\n\na\nab\n\nabc\nabc\xff\n01\n",
			p:     params{n: 1},
			want:  "a\nab\nabc\nabc\n01\n",
		},
		{
			name:  "larger value of n",
			files: []string{filepath.Join(tmpdir, "file")},
			input: "\n\na123456\nab\n\nabc\nabcde\xff\n01\n",
			p:     params{n: 6},
			want:  "a123456\n",
		},
		{
			name:    "wrong format",
			p:       params{t: "wrong", n: 1},
			wantErr: errInvalidFormatArgument,
		},
		{
			name:  "offset all readable no newline",
			input: "hello",
			p:     params{t: "d", n: 5},
			want:  "0 hello\n",
		},
		{
			name:  "offset all readable",
			input: "hello\n",
			p:     params{t: "d", n: 5},
			want:  "0 hello\n",
		},
		{
			name:  "offset with hex",
			input: "\xffhelloworld\xffhello",
			p:     params{t: "x", n: 5},
			want:  "1 helloworld\nc hello\n",
		},
		{
			name:  "offset with octal",
			input: "\xffhelloworld\xffhello",
			p:     params{t: "o", n: 5},
			want:  "1 helloworld\n14 hello\n",
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

			if err := command(bufIn, bufOut, tt.p, tt.files).run(); err != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Errorf("strings() = %q, want: %q", err.Error(), tt.wantErr)
				}
			} else {
				if bufOut.String() != tt.want {
					t.Errorf("strings() = %q, want: %q", bufOut.String(), tt.want)
				}
			}
		})
	}
}
