// Copyright 2016-2023 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"fmt"
	"os"
	"testing"
)

func TestFiles(t *testing.T) {
	tmpDir := t.TempDir()
	f1, err := os.CreateTemp(tmpDir, "")
	if err != nil {
		t.Fatal(err)
	}
	f2, err := os.CreateTemp(tmpDir, "")
	if err != nil {
		t.Fatal(err)
	}
	_, err = f1.Write([]byte("simple test count words\nlines\nlines\n"))
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name       string
		want       string
		wantStderr string
		args       []string
		p          params
	}{
		{
			name: "count in 2 files",
			args: []string{f1.Name(), f2.Name()},
			want: "3 6 36 " + f1.Name() + "\n0 0 0 " + f2.Name() + "\n3 6 36 total\n",
		},
		{
			name:       "file does not exist",
			args:       []string{"filedoesnotexist"},
			wantStderr: "wc: filedoesnotexist: open filedoesnotexist: no such file or directory\n",
		},
		{
			name:       "both files do not exist but total printed",
			args:       []string{"filedoesnotexist1", "filedoesnotexist2"},
			wantStderr: "wc: filedoesnotexist1: open filedoesnotexist1: no such file or directory\n" + "wc: filedoesnotexist2: open filedoesnotexist2: no such file or directory\n",
			want:       "0 0 0 total\n",
		},
		{
			name: "file1 no flag",
			args: []string{f1.Name()},
			want: fmt.Sprintf("3 6 36 %s\n", f1.Name()),
		},
		{
			name: "count words in file1",
			args: []string{f1.Name()},
			p:    params{words: true},
			want: fmt.Sprintf("6 %s\n", f1.Name()),
		},
		{
			name: "count lines in file1",
			args: []string{f1.Name()},
			p:    params{lines: true},
			want: fmt.Sprintf("3 %s\n", f1.Name()),
		},
		{
			name: "count lines in file1",
			args: []string{f1.Name()},
			p:    params{chars: true},
			want: fmt.Sprintf("36 %s\n", f1.Name()),
		},
		{
			name: "count runes in file1",
			args: []string{f1.Name()},
			p:    params{runes: true},
			want: fmt.Sprintf("36 %s\n", f1.Name()),
		},
		{
			name: "count broken in file1",
			args: []string{f1.Name()},
			p:    params{broken: true},
			want: fmt.Sprintf("0 %s\n", f1.Name()),
		},
		{
			name: "count in 2 files",
			args: []string{f1.Name(), f2.Name()},
			want: fmt.Sprintf("3 6 36 %s\n0 0 0 %s\n3 6 36 total\n", f1.Name(), f2.Name()),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			stdout := &bytes.Buffer{}
			stderr := &bytes.Buffer{}
			if err := command(nil, stdout, stderr, test.p, test.args).run(); err != nil {
				t.Fatal(err)
			}
			if stdout.String() != test.want {
				t.Errorf("wc stdout = %q, want: %q", stdout.String(), test.want)
			}
			if stderr.String() != test.wantStderr {
				t.Errorf("wc stderr = %q, want: %q", stderr.String(), test.wantStderr)
			}
		})
	}
}

func TestStdin(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
		args  []string
		p     params
	}{
		{
			name:  "ascii",
			input: "hello world\n",
			want:  "1 2 12\n",
		},
		{
			name:  "utf8",
			input: "München\n",
			want:  "1 1 9\n",
		},
		{
			name:  "empty",
			input: "",
			want:  "0 0 0\n",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			stdin := bytes.NewBufferString(test.input)
			stdout := &bytes.Buffer{}

			if err := command(stdin, stdout, nil, test.p, test.args).run(); err != nil {
				t.Fatal(err)
			}

			if got := stdout.String(); got != test.want {
				t.Errorf("wc %q = %q, want: %q", test.input, got, test.want)
			}
		})
	}
}
