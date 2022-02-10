// Copyright 2016-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"testing"
)

func TestWc(t *testing.T) {
	for _, tt := range []struct {
		name    string
		args    []string
		lines   bool
		words   bool
		runes   bool
		broken  bool
		chars   bool
		want    string
		wantErr string
	}{
		{
			name: "no flags no file",
			want: "1 1 6\n",
		},
		{
			name:    "file does not exist",
			args:    []string{"filedoesnotexist"},
			wantErr: "open filedoesnotexist: no such file or directory",
		},
		{
			name: "file1 no flag",
			args: []string{"testfiles/file1"},
			want: "3 6 36 testfiles/file1\n",
		},
		{
			name:  "count words in file1",
			args:  []string{"testfiles/file1"},
			words: true,
			want:  "6 testfiles/file1\n",
		},
		{
			name:  "count lines in file1",
			args:  []string{"testfiles/file1"},
			lines: true,
			want:  "3 testfiles/file1\n",
		},
		{
			name:  "count lines in file1",
			args:  []string{"testfiles/file1"},
			chars: true,
			want:  "36 testfiles/file1\n",
		},
		{
			name:  "count runes in file1",
			args:  []string{"testfiles/file1"},
			runes: true,
			want:  "36 testfiles/file1\n",
		},
		{
			name:   "count broken in file1",
			args:   []string{"testfiles/file1"},
			broken: true,
			want:   "0 testfiles/file1\n",
		},
		{
			name: "count in 2 files",
			args: []string{"testfiles/file1", "testfiles/file2"},
			want: "3 6 36 testfiles/file1\n0 0 0 testfiles/file2\n3 6 36 total\n",
		},
	} {
		*lines = tt.lines
		*words = tt.words
		*runes = tt.runes
		*broken = tt.broken
		*chars = tt.chars
		t.Run(tt.name, func(t *testing.T) {
			bufOut := &bytes.Buffer{}
			bufIn := &bytes.Buffer{}
			if tt.name == "no flags no file" {
				bufIn.Write([]byte{0x41, 0x42, 0x43, 0x44, 0x45, 0x0a})
			}
			if got := runwc(bufOut, bufIn, tt.args...); got != nil {
				if got.Error() != tt.wantErr {
					t.Errorf("runwc() = %q, want: %q", got.Error(), tt.wantErr)
				}
			} else {
				if bufOut.String() != tt.want {
					t.Errorf("runwc() = %q, want: %q", bufOut.String(), tt.want)
				}
			}
		})
	}
}
