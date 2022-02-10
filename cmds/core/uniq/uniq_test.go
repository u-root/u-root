// Copyright 2016 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"log"
	"testing"
)

func TestUniq(t *testing.T) {
	for _, tt := range []struct {
		name       string
		args       []string
		uniques    bool
		duplicates bool
		count      bool
		want       string
		wantErr    string
	}{
		{
			name:    "file 1 with wrong file",
			args:    []string{"filedoesnotexist"},
			wantErr: "open filedoesnotexist: no such file or directory",
		},
		{
			name: "file 1 without any flag",
			args: []string{"testfiles/file1"},
			want: "test\ngo\ncoool\ncool\nlegaal\ntest\n",
		},
		{
			name:  "file 1 count == true",
			args:  []string{"testfiles/file1"},
			count: true,
			want:  "2\ttest\n3\tgo\n2\tcoool\n1\tcool\n1\tlegaal\n1\ttest\n",
		},
		{
			name:    "file 1 uniques == true",
			args:    []string{"testfiles/file1"},
			uniques: true,
			want:    "cool\nlegaal\ntest\n",
		},
		{
			name:       "file 1 duplicates == true",
			args:       []string{"testfiles/file1"},
			duplicates: true,
			want:       "test\ngo\ncoool\n",
		},
		{
			name: "file 2 without any flag",
			args: []string{"testfiles/file2"},
			want: "u-root\nuniq\nron\nteam\nbinaries\ntest\n\n",
		},
		{
			name:  "file 2 count == true",
			args:  []string{"testfiles/file2"},
			count: true,
			want:  "1\tu-root\n1\tuniq\n2\tron\n1\tteam\n1\tbinaries\n1\ttest\n5\t\n",
		},
		{
			name:    "file 2 uniques == true",
			args:    []string{"testfiles/file2"},
			uniques: true,
			want:    "u-root\nuniq\nteam\nbinaries\ntest\n",
		},
		{
			name:       "file 2 duplicates == true",
			args:       []string{"testfiles/file2"},
			duplicates: true,
			want:       "ron\n\n",
		},
	} {
		*uniques = tt.uniques
		*duplicates = tt.duplicates
		*count = tt.count
		buf := &bytes.Buffer{}
		log.SetOutput(buf)
		t.Run(tt.name, func(t *testing.T) {
			if got := runUniq(buf, tt.args...); got != nil {
				if got.Error() != tt.wantErr {
					t.Errorf("runUniq() = %q, want %q", got.Error(), tt.wantErr)
				}
			} else {
				if buf.String() != tt.want {
					t.Errorf("runUniq() = %q, want %q", buf.String(), tt.want)
				}
			}
		})
	}
}
