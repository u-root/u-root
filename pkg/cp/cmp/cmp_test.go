// Copyright 2021 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cmp

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/u-root/u-root/pkg/cp"
)

func TestCMP(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "u-root-pkg-cmp-")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	tmpFile1, err := ioutil.TempFile(tmpDir, "file1")
	if err != nil {
		t.Fatal(err)
	}
	defer tmpFile1.Close()

	tmpFile2, err := ioutil.TempFile(tmpDir, "file2")
	if err != nil {
		t.Fatal(err)
	}
	defer tmpFile2.Close()

	tmpFile3, err := ioutil.TempFile(tmpDir, "file3")
	if err != nil {
		t.Fatal(err)
	}
	defer tmpFile3.Close()

	if err := os.WriteFile(tmpFile1.Name(), []byte("F is for fire that burns down the whole town"), 0766); err != nil {
		t.Fatal(err)
	}

	if err := os.WriteFile(tmpFile2.Name(), []byte("F is for fire that burns down the whole town"), 0766); err != nil {
		t.Fatal(err)
	}

	if err := os.WriteFile(tmpFile3.Name(), []byte("nwot elohw eht nwod snrub taht erif rof si F"), 0766); err != nil {
		t.Fatal(err)
	}

	//function isEqualFile
	var tests = []struct {
		n     string
		file1 string
		file2 string
		err   string
	}{
		{n: "file1 does not exist", file1: "file1", file2: tmpFile2.Name(), err: "open file1: no such file or directory"},
		{n: "file2 does not exist", file1: tmpFile1.Name(), file2: "file2", err: "open file2: no such file or directory"},
		{n: "files are not equal", file1: tmpFile1.Name(), file2: tmpFile3.Name(), err: fmt.Sprintf("%q and %q do not have equal content", tmpFile1.Name(), tmpFile3.Name())},
	}

	for _, tt := range tests {
		err := isEqualFile(tt.file1, tt.file2)
		if err.Error() != tt.err {
			t.Errorf("Test %s: got: (%s), want: (%s)", tt.n, err.Error(), tt.err)
		}
	}
	err = isEqualFile(tmpFile1.Name(), tmpFile2.Name())
	if err != nil {
		t.Errorf("got: (%s), want: (%s)", err.Error(), "")
	}

	//function readDirNames
	names, err := readDirNames(tmpDir)
	if len(names) != 3 || names[0] != filepath.Base(tmpFile1.Name()) || names[1] != filepath.Base(tmpFile2.Name()) || names[2] != filepath.Base(tmpFile3.Name()) || err != nil {
		t.Errorf("file amount: %d, files: %v, files created %s, %s, %s", len(names), names, filepath.Base(tmpFile1.Name()), filepath.Base(tmpFile2.Name()), filepath.Base(tmpFile3.Name()))
	}
	_, err = readDirNames("dir")
	if err.Error() != "open dir: no such file or directory" {
		t.Errorf("got: (%s), want: (%s)", err.Error(), "")
	}

	// function stat
	options := cp.Default
	_, err = stat(options, tmpFile1.Name())
	if err != nil {
		t.Fatal(err)
	}
	options.NoFollowSymlinks = true
	_, err = stat(options, tmpFile1.Name())
	if err != nil {
		t.Fatal(err)
	}
}
