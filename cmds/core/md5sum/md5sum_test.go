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

func TestMd5Sum(t *testing.T) {
	// Creating tmp files with data to hash
	tmpdir := t.TempDir()
	file1, err := os.Create(filepath.Join(tmpdir, "file1"))
	if err != nil {
		t.Errorf("failed to create tmp file1: %v", err)
	}
	if _, err := file1.WriteString("abcdef\n"); err != nil {
		t.Errorf("failed to write string to file1: %v", err)
	}
	file2, err := os.Create(filepath.Join(tmpdir, "file2"))
	if err != nil {
		t.Errorf("failed to create tmp file2: %v", err)
	}
	if _, err := file2.WriteString("pqra\n"); err != nil {
		t.Errorf("failed to write string to file2: %v", err)
	}

	for _, tt := range []struct {
		name string
		args []string
		want string
	}{
		{
			name: "bufIn as input",
			args: []string{},
			want: "8d777f385d3dfec8815d20f7496026dc\n",
		},
		{
			name: "wrong path file",
			args: []string{"testfile"},
			want: "open testfile: no such file or directory",
		},
		{
			name: "file1 as input",
			args: []string{file1.Name()},
			want: fmt.Sprintf("%s %s\n", "5ab557c937e38f15291c04b7e99544ad", file1.Name()),
		},
		{
			name: "file2 as input",
			args: []string{file2.Name()},
			want: fmt.Sprintf("%s %s\n", "721d6b135656aa83baca6ebdbd2f6c86", file2.Name()),
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			// Setting flags
			bufIn := &bytes.Buffer{}
			if _, err := bufIn.WriteString("data"); err != nil {
				t.Errorf("failed to write string to bufIn: %v", err)
			}
			bufOut := &bytes.Buffer{}
			if got := md5Sum(bufOut, bufIn, tt.args...); got != nil {
				if got.Error() != tt.want {
					t.Errorf("md5Sum() = %q, want: %q", got.Error(), tt.want)
				}
			} else {
				if bufOut.String() != tt.want {
					t.Errorf("md5Sum() = %q, want: %q", bufOut.String(), tt.want)
				}
			}
		})
	}
}
