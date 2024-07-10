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

func TestSHASum(t *testing.T) {
	// Creating tmp files with data to hash
	tmpdir := t.TempDir()
	file1, err := os.Create(filepath.Join(tmpdir, "file1"))
	if err != nil {
		t.Errorf("failed to create tmp file1: %v", err)
	}
	defer file1.Close()
	if _, err := file1.WriteString("abcdef\n"); err != nil {
		t.Errorf("failed to write string to file1: %v", err)
	}
	file2, err := os.Create(filepath.Join(tmpdir, "file2"))
	if err != nil {
		t.Errorf("failed to create tmp file2: %v", err)
	}
	defer file2.Close()
	if _, err := file2.WriteString("pqra\n"); err != nil {
		t.Errorf("failed to write string to file2: %v", err)
	}

	for _, tt := range []struct {
		name      string
		args      []string
		algorithm int
		want      string
	}{
		{
			name:      "bufIn as input with sha1 sum",
			args:      []string{},
			algorithm: 1,
			want:      "bdc37c074ec4ee6050d68bc133c6b912f36474df -\n",
		},
		{
			name:      "bufIn as input with sha256 sum",
			args:      []string{},
			algorithm: 256,
			want:      "ae0666f161fed1a5dde998bbd0e140550d2da0db27db1d0e31e370f2bd366a57 -\n",
		},
		{
			name: "wrong path file",
			args: []string{"testfile"},
			want: "open testfile: no such file or directory",
		},
		{
			name: "file1 as input with invalid algorithm",
			args: []string{file1.Name()},
			want: "invalid algorithm, only 1 or 256 are valid",
		},
		{
			name: "stdin as input with invalid algorithm",
			args: []string{},
			want: "invalid algorithm, only 1 or 256 are valid",
		},
		{
			name:      "file1 as input with sha1 sum",
			args:      []string{file1.Name()},
			algorithm: 1,
			want:      fmt.Sprintf("%s %s\n", "bdc37c074ec4ee6050d68bc133c6b912f36474df", file1.Name()),
		},
		{
			name:      "file2 as input with sha1 sum",
			args:      []string{file2.Name()},
			algorithm: 1,
			want:      fmt.Sprintf("%s %s\n", "e8ed2d487f1dc32152c8590f39c20b7703f9e159", file2.Name()),
		},
		{
			name:      "file1 as input with sha256 sum",
			args:      []string{file1.Name()},
			algorithm: 256,
			want:      fmt.Sprintf("%s %s\n", "ae0666f161fed1a5dde998bbd0e140550d2da0db27db1d0e31e370f2bd366a57", file1.Name()),
		},
		{
			name:      "file2 as input with sha256 sum",
			args:      []string{file2.Name()},
			algorithm: 256,
			want:      fmt.Sprintf("%s %s\n", "db296dd0bcb796df9b327f44104029da142c8fff313a25bd1ac7c3b7562caea9", file2.Name()),
		},
		{
			name:      "file1 and file 2 as input with sha256 sum",
			args:      []string{file1.Name(), file2.Name()},
			algorithm: 256,
			want: fmt.Sprintf("%s %s\n%s %s\n", "ae0666f161fed1a5dde998bbd0e140550d2da0db27db1d0e31e370f2bd366a57", file1.Name(),
				"db296dd0bcb796df9b327f44104029da142c8fff313a25bd1ac7c3b7562caea9", file2.Name()),
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			// Setting flags
			algorithm = tt.algorithm
			bufIn := &bytes.Buffer{}
			if _, err := bufIn.WriteString("abcdef\n"); err != nil {
				t.Errorf("failed to write string to bufIn: %v", err)
			}
			bufOut := &bytes.Buffer{}
			if got := shasum(bufOut, bufIn, tt.args...); got != nil {
				if got.Error() != tt.want {
					t.Errorf("shasum() = %q, want: %q", got.Error(), tt.want)
				}
			} else {
				if bufOut.String() != tt.want {
					t.Errorf("shasum() = %q, want: %q", bufOut.String(), tt.want)
				}
			}
		})
	}
}
