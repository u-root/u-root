// Copyright 2023 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// By Ahmed Kamal <email.ahmedkamal@googlemail.com>

package main

import (
	"fmt"
	"io"
	"os"
	"testing"
)

func TestTmpWriter(t *testing.T) {
	tmpDir := t.TempDir()
	f1, err := os.CreateTemp(tmpDir, "f1")
	if err != nil {
		t.Fatal(err)
	}

	_, err = f1.WriteString("hix\nnix\n")
	if err != nil {
		t.Fatal(err)
	}

	err = os.Chdir(tmpDir)
	if err != nil {
		t.Fatal(err)
	}

	testcases := []struct {
		filename string
		content  string
		err      error
	}{
		{
			filename: "/tmp/tw",
			content:  "foo\nbar",
			err:      nil,
		},
	}

	for idx, tc := range testcases {
		test := tc
		t.Run(fmt.Sprintf("case_%d", idx), func(t *testing.T) {
			tw, err := newTmpWriter(test.filename)
			if err != nil {
				t.Errorf("failed to create tempWriter: %v", err)
			}
			fmt.Fprint(tw, test.content)
			tw.Close()
			fh, _ := os.Open(tc.filename)
			fc, _ := io.ReadAll(fh)
			fcontent := string(fc)
			if fcontent != tc.content {
				t.Errorf("got %#v, want %#v", fcontent, tc.content)
			}
		})
	}
}
