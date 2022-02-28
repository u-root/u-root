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

func TestCksum(t *testing.T) {
	tmpDir := t.TempDir()
	for _, tt := range []struct {
		name    string
		args    []string
		data    []byte
		help    bool
		version bool
		want    string
	}{
		{
			name: "file does not exist",
			args: []string{"cksum", filepath.Join(tmpDir, "files")},
			want: fmt.Sprintf("open %s: no such file or directory", filepath.Join(tmpDir, "files")),
		},
		{
			name: "help flag true",
			help: true,
			want: "Usage:\ncksum <File Name>\n",
		},
		{
			name:    "version flag true",
			version: true,
			want:    "cksum utility, URoot Version.\n",
		},
		{
			name: "small string abcdef from buffer",
			data: []byte("abcdef\n"),
			want: fmt.Sprintf("%s %d \n", "3512391007", 7),
		},
		{
			name: "small string pqra from buffer",
			data: []byte("pqra\n"),
			want: fmt.Sprintf("%s %d \n", "1063566492", 5),
		},
		{
			name: "big string from buffer",
			data: []byte("abcdef\nafdsfsfgdglfdgkdlvcmdfposvpomfvmlcvdfk;lgkd'f;k;lvcvcv\n" +
				"afdsfsfgdglfdgkdlvcmdfposvpomfvmlcvdfk;lgkd'f;k;lvcvcv\nafdsfsfg" +
				"dglfdgkdlvcmdfposvpomfvmlcvdfk;lgkd'f;k;lvcvcv\nafdsfsfgdglfdgkd" +
				"lvcmdfposvpomfvmlcvdfk;lgkd'f;k;lvcvcv\nsdddsfsfsdfsdfsdasaarwre" +
				"mazadsfssfsfsfsafsadfsfdsadfsafsafsfsafdsfsdfsfdsdf"),
			want: fmt.Sprintf("%s %d \n", "689622513", 302),
		},
		{
			name: "small string abcdef from file",
			args: []string{"cksum", filepath.Join(tmpDir, "file")},
			data: []byte("abcdef\n"),
			want: fmt.Sprintf("%s %d %s\n", "3512391007", 7, filepath.Join(tmpDir, "file")),
		},
		{
			name: "small string pqra from file",
			args: []string{"cksum", filepath.Join(tmpDir, "file")},
			data: []byte("pqra\n"),
			want: fmt.Sprintf("%s %d %s\n", "1063566492", 5, filepath.Join(tmpDir, "file")),
		},
		{
			name: "big string from file",
			args: []string{"cksum", filepath.Join(tmpDir, "file")},
			data: []byte("abcdef\nafdsfsfgdglfdgkdlvcmdfposvpomfvmlcvdfk;lgkd'f;k;lvcvcv\n" +
				"afdsfsfgdglfdgkdlvcmdfposvpomfvmlcvdfk;lgkd'f;k;lvcvcv\nafdsfsfg" +
				"dglfdgkdlvcmdfposvpomfvmlcvdfk;lgkd'f;k;lvcvcv\nafdsfsfgdglfdgkd" +
				"lvcmdfposvpomfvmlcvdfk;lgkd'f;k;lvcvcv\nsdddsfsfsdfsdfsdasaarwre" +
				"mazadsfssfsfsfsafsadfsfdsadfsafsafsfsafdsfsdfsfdsdf"),
			want: fmt.Sprintf("%s %d %s\n", "689622513", 302, filepath.Join(tmpDir, "file")),
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			// Set flags
			*help = tt.help
			*version = tt.version

			// Create Buffer
			bufIn := &bytes.Buffer{}
			bufOut := &bytes.Buffer{}
			if _, err := bufIn.Write(tt.data); err != nil {
				t.Errorf("failed to write to buffer: %v", err)
			}

			// Create file
			file, err := os.Create(filepath.Join(tmpDir, "file"))
			if err != nil {
				t.Errorf("failed to create file: %v", err)
			}
			defer file.Close()
			// Write data in file that should be compared
			if _, err = file.Write(tt.data); err != nil {
				t.Errorf("failed to write to file: %v", err)
			}

			if got := cksum(bufOut, bufIn, tt.args...); got != nil {
				if got.Error() != tt.want {
					t.Errorf("cksum = %q, want: %q", got.Error(), tt.want)
				}
			} else {
				if bufOut.String() != tt.want {
					t.Errorf("cksum = %q, want: %q", bufOut.String(), tt.want)
				}
			}
		})
	}
}
