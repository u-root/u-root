// Copyright 2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestHexdump(t *testing.T) {
	const testString = "abcdefghijklmnopqrstuvwxyz"
	// Creating file and write content into it for testing purposes
	d := t.TempDir()
	n := filepath.Join(d, "testfile")
	if err := os.WriteFile(n, []byte(testString), 0o644); err != nil {
		t.Fatal(err)
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
			readInput: testString,
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
			readBuf := &bytes.Buffer{}
			writeBuf := &bytes.Buffer{}
			readBuf.WriteString(tt.readInput)
			if got := hexdump(tt.filenames, readBuf, writeBuf); got != nil {
				// Different Go compilers deliver
				// different errors: either syscall.ENOENT
				// or os.ErrNotExist.
				// Once that is fixed, we can use errors.Is
				if tt.wantErr == nil {
					t.Fatalf("hexdump() = '%v', want: nil", got)
				}
				if !strings.HasPrefix(got.Error(), tt.wantErr.Error()[:10]) {
					t.Fatalf("hexdump() = '%v', want: '%v'", got, tt.wantErr)
				}
			}
			if writeBuf.String() != tt.want {
				t.Errorf("Console output: %q, want: %q", writeBuf.String(), tt.want)
			}
		})
	}
}
