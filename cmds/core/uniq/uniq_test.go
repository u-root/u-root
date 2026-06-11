// Copyright 2016 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"errors"
	"io"
	"log"
	"os"
	"strings"
	"testing"
	"testing/iotest"
)

func TestUniq(t *testing.T) {
	for _, tt := range []struct {
		name       string
		args       []string
		unique     bool
		duplicates bool
		count      bool
		ignoreCase bool
		want       string
		wantErr    error
		stdin      io.Reader
	}{
		{
			name:    "file 1 with wrong file",
			args:    []string{"filedoesnotexist"},
			wantErr: os.ErrNotExist,
		},
		{
			name: "file 1 without any flag",
			args: []string{"testdata/file1.txt"},
			want: "test\ngo\ncoool\ncool\nlegaal\ntest\n",
		},
		{
			name:  "file 1 count == true",
			args:  []string{"testdata/file1.txt"},
			count: true,
			want:  "2\ttest\n3\tgo\n2\tcoool\n1\tcool\n1\tlegaal\n1\ttest\n",
		},
		{
			name:   "file 1 uniques == true",
			args:   []string{"testdata/file1.txt"},
			unique: true,
			want:   "cool\nlegaal\ntest\n",
		},
		{
			name:       "file 1 duplicates == true",
			args:       []string{"testdata/file1.txt"},
			duplicates: true,
			want:       "test\ngo\ncoool\n",
		},
		{
			name: "file 2 without any flag",
			args: []string{"testdata/file2.txt"},
			want: "u-root\nuniq\nron\nteam\nbinaries\ntest\nTest\n\n",
		},
		{
			name:  "file 2 count == true",
			args:  []string{"testdata/file2.txt"},
			count: true,
			want:  "1\tu-root\n1\tuniq\n2\tron\n1\tteam\n1\tbinaries\n1\ttest\n1\tTest\n5\t\n",
		},
		{
			name:   "file 2 uniques == true",
			args:   []string{"testdata/file2.txt"},
			unique: true,
			want:   "u-root\nuniq\nteam\nbinaries\ntest\nTest\n",
		},
		{
			name:       "file 2 duplicates == true",
			args:       []string{"testdata/file2.txt"},
			duplicates: true,
			want:       "ron\n\n",
		},
		{
			name:       "file 2 ignore case == true",
			args:       []string{"testdata/file2.txt"},
			ignoreCase: true,
			want:       "u-root\nuniq\nron\nteam\nbinaries\ntest\n\n",
		},
		{
			name:    "error reading from stdin 1",
			args:    nil,
			wantErr: iotest.ErrTimeout,
			stdin:   iotest.ErrReader(iotest.ErrTimeout),
		},
		{
			name:    "error reading from stdin 2",
			args:    nil,
			wantErr: iotest.ErrTimeout,
			stdin:   io.MultiReader(strings.NewReader("go\n"), iotest.ErrReader(iotest.ErrTimeout)),
		},
		{
			name:  "stdin empty input",
			args:  nil,
			want:  "",
			stdin: strings.NewReader(""),
		},
		{
			name:  "stdin one string",
			args:  nil,
			want:  "go\n",
			stdin: strings.NewReader("go\n"),
		},
		{
			name:  "stdin three identical strings",
			args:  nil,
			want:  "go\n",
			stdin: strings.NewReader("go\ngo\ngo\n"),
		},
		{
			name:   "stdin and unique",
			args:   nil,
			unique: true,
			stdin:  strings.NewReader("go\nu-root\ngo\ngo\ngo\n"),
			want:   "go\nu-root\n",
		},
		{
			name:  "same strings but no new line",
			args:  nil,
			stdin: strings.NewReader("go\ngo"),
			want:  "go\n",
		},
	} {
		buf := &bytes.Buffer{}
		log.SetOutput(buf)
		t.Run(tt.name, func(t *testing.T) {
			if got := run(tt.stdin, buf, tt.unique, tt.duplicates, tt.count, tt.ignoreCase, tt.args); got != nil {
				if !errors.Is(got, tt.wantErr) {
					t.Fatalf("runUniq() = %q, want %q", got.Error(), tt.wantErr)
				}
			} else {
				if tt.wantErr != nil {
					t.Fatalf("runUniq() = <nil>, want %q", tt.wantErr)
				}
				if buf.String() != tt.want {
					t.Errorf("runUniq() = %q, want %q", buf.String(), tt.want)
				}
			}
		})
	}
}
