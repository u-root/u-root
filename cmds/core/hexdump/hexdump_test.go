// Copyright 2017 the u-root Authors. All rights reserved
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

func TestHexdump(t *testing.T) {
	// Creating file and write content into it for testing purposes
	d := t.TempDir()
	f, err := os.Create(filepath.Join(d, "testfile"))
	if err != nil {
		t.Errorf("failed to create tmp file: %v", err)
	}
	if _, err = f.WriteString("abcdefghijklmnopqrstuvwxyz"); err != nil {
		t.Errorf("failed to write string into tmp file: %v", err)
	}

	for _, tt := range []struct {
		name      string
		filenames []string
		readInput string
		want      string
		wantErr   error
	}{
		{
			name:      "hexdump from Stdin",
			filenames: []string{},
			readInput: "abcdefghijklmnopqrstuvwxyz",
			want: `00000000  61 62 63 64 65 66 67 68  69 6a 6b 6c 6d 6e 6f 70  |abcdefghijklmnop|
00000010  71 72 73 74 75 76 77 78  79 7a                    |qrstuvwxyz|
`,
		},
		{
			name:      "hexdump from file",
			filenames: []string{filepath.Join(d, "testfile")},
			want: `00000000  61 62 63 64 65 66 67 68  69 6a 6b 6c 6d 6e 6f 70  |abcdefghijklmnop|
00000010  71 72 73 74 75 76 77 78  79 7a                    |qrstuvwxyz|
`,
		},
		{
			name:      "hexdump from 2 files",
			filenames: []string{filepath.Join(d, "testfile"), filepath.Join(d, "testfile")},
			want: `00000000  61 62 63 64 65 66 67 68  69 6a 6b 6c 6d 6e 6f 70  |abcdefghijklmnop|
00000010  71 72 73 74 75 76 77 78  79 7a 61 62 63 64 65 66  |qrstuvwxyzabcdef|
00000020  67 68 69 6a 6b 6c 6d 6e  6f 70 71 72 73 74 75 76  |ghijklmnopqrstuv|
00000030  77 78 79 7a                                       |wxyz|
`,
		},
		{
			name:      "hexdump from file that does not exist",
			filenames: []string{"error"},
			wantErr:   fmt.Errorf("open %s: no such file or directory", "error"),
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			var readBuf = &bytes.Buffer{}
			var writeBuf = &bytes.Buffer{}
			readBuf.WriteString(tt.readInput)
			if got := hexdump(tt.filenames, readBuf, writeBuf); got != nil {
				if got.Error() != tt.wantErr.Error() {
					t.Errorf("hexdump() = '%v', want: '%v'", got, tt.wantErr)
				}
			} else {
				if writeBuf.String() != tt.want {
					t.Errorf("Console output: '%s', want: '%s'", writeBuf.String(), tt.want)
				}
			}
		})
	}
}
