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
	// Creating all tmp dirs and files for testing purpose
	dirPath, err := os.MkdirTemp("", "cmp_test")
	if err != nil {
		t.Fatalf("Failed to create tmp dir: %v", err)
	}
	defer os.RemoveAll(dirPath)

	for i := 1; i < 7; i++ {
		if err := os.Mkdir(filepath.Join(dirPath, fmt.Sprint(i)), 0700); err != nil {
			t.Fatalf("Failed to create %s: %v", filepath.Join(dirPath, fmt.Sprint(i)), err)
		}
		if i == 5 {
			if err := os.Mkdir(filepath.Join(dirPath, fmt.Sprint(i), "7"), 0700); err != nil {
				t.Fatalf("Failed to create %s: %v", filepath.Join(dirPath, fmt.Sprint(i), "7"), err)
			}
		}
		for j := 1; j < 5; j++ {
			if j < 3 {
				if err := ioutil.WriteFile(filepath.Join(dirPath, fmt.Sprint(i), fmt.Sprint(j)), []byte("F is for fire that burns down the whole town"), 0766); err != nil {
					t.Fatalf("Failed to write to %s: %v", filepath.Join(dirPath, fmt.Sprint(i), fmt.Sprint(j)), err)
				}
			} else {
				if err := ioutil.WriteFile(filepath.Join(dirPath, fmt.Sprint(i), fmt.Sprint(j)), []byte("nwot elohw eht nwod snrub taht erif rof si F"), 0766); err != nil {
					t.Fatalf("Failed to write to %s: %v", filepath.Join(dirPath, fmt.Sprint(i), fmt.Sprint(j)), err)
				}
			}
		}
	}

	// Tests start here
	//
	//
	// Struct for testing isEqualFile and readDirNames
	var testTable1 = []struct {
		name  string
		file1 string
		file2 string
		err   string
	}{
		{
			name:  "1 does not exist",
			file1: "1",
			file2: filepath.Join(dirPath, "1", "2"),
			err:   "open 1: no such file or directory",
		},
		{
			name:  "2 does not exist",
			file1: filepath.Join(dirPath, "1", "1"),
			file2: "2",
			err:   "open 2: no such file or directory",
		},
		{
			name:  "files are not equal",
			file1: filepath.Join(dirPath, "1", "1"),
			file2: filepath.Join(dirPath, "1", "3"),
			err:   fmt.Sprintf("%q and %q do not have equal content", filepath.Join(dirPath, "1", "1"), filepath.Join(dirPath, "1", "3")),
		},
	}

	// Testing isEqualFile
	t.Run("Test isEqualFile", func(t *testing.T) {
		for _, tt := range testTable1 {
			err := isEqualFile(tt.file1, tt.file2)
			if err.Error() != tt.err {
				t.Errorf("Test %s: got: (%s), want: (%s)", tt.name, err.Error(), tt.err)
			}
		}
		err := isEqualFile(filepath.Join(dirPath, "1", "1"), filepath.Join(dirPath, "1", "2"))
		if err != nil {
			t.Errorf("got: (%s), want: (%s)", err.Error(), "")
		}
	})

	// Testing readDirNames
	t.Run("Test readDirNames", func(t *testing.T) {
		names, err := readDirNames(filepath.Join(dirPath, "1"))
		if len(names) != 4 || names[0] != filepath.Base(filepath.Join(dirPath, "1", "1")) || names[1] != filepath.Base(filepath.Join(dirPath, "1", "2")) ||
			names[2] != filepath.Base(filepath.Join(dirPath, "1", "3")) || err != nil {
			t.Errorf("file amount: %d, files: %v, files created %s, %s, %s, %s",
				len(names), names, filepath.Base(filepath.Join(dirPath, "1", "1")), filepath.Base(filepath.Join(dirPath, "1", "2")),
				filepath.Base(filepath.Join(dirPath, "1", "3")), filepath.Base(filepath.Join(dirPath, "1", "4")))
		}
		_, err = readDirNames("dir")
		if err.Error() != "open dir: no such file or directory" {
			t.Errorf("got: (%s), want: (%s)", err.Error(), "")
		}
	})

	// Default option var
	equalTreeOpts := cp.Default

	// Testing stats and the IsEqualTree 2 dirs equal
	t.Run("Test stats and the IsEqualTree 2 dirs equal", func(t *testing.T) {
		// Struct for testing isEqualTree
		var testTable2 = []struct {
			name string
			src  string
			dst  string
			err  string
		}{
			{
				name: "stat src err",
				src:  filepath.Join(dirPath, "1", "1"),
				dst:  "",
				err:  "stat : no such file or directory",
			},
			{
				name: "stat dst err",
				src:  "",
				dst:  filepath.Join(dirPath, "1", "2"),
				err:  "stat : no such file or directory",
			},
			{
				name: "2 dirs are equal",
				src:  filepath.Join(dirPath, "4"),
				dst:  filepath.Join(dirPath, "6"),
				err:  "<nil>",
			},
		}

		for _, tt := range testTable2 {
			_, _, _, err := stats(equalTreeOpts, tt.src, tt.dst)

			if fmt.Sprintf("%v", err) != tt.err {
				t.Errorf("Test %s: got: (%s), want: (%s)", tt.name, err, tt.err)
			}
		}
	})

	// Testing IsEqualTree for case dir
	t.Run("Test IsEqualTree for case dir", func(t *testing.T) {
		// Test case that 2 dirs are equal
		err := IsEqualTree(equalTreeOpts, filepath.Join(dirPath, "4"), filepath.Join(dirPath, "6"))
		if fmt.Sprintf("%v", err) != "<nil>" {
			t.Errorf("Test %s: got: (%s), want: (%s)\n", "2 dirs are equal", err, "<nil>")
		}

		// retrieve sm and dm for err checking
		sm, dm, _, err := stats(equalTreeOpts, filepath.Join(dirPath, "1"), filepath.Join(dirPath, "1", "2"))
		if err != nil {
			t.Errorf("err is: %v", err)
		}

		// retrive srcEntries and dstEntries
		srcEntries, err := readDirNames(filepath.Join(dirPath, "3"))
		if err != nil {
			t.Errorf("err is: %v", err)
		}
		dstEntries, err := readDirNames(filepath.Join(dirPath, "4"))
		if err != nil {
			t.Errorf("err is: %v", err)
		}

		// Struct for testing isEqualTree
		var testTable3 = []struct {
			name string
			src  string
			dst  string
			err  string
		}{
			{
				name: "mismatched mode, 1 dir 1 file",
				src:  filepath.Join(dirPath, "1"),
				dst:  filepath.Join(dirPath, "1", "2"),
				err:  fmt.Sprintf("mismatched mode: %q has mode %s while %q has mode %s", filepath.Join(dirPath, "1"), sm, filepath.Join(dirPath, "1", "2"), dm),
			},
			{
				name: "err in first readDirName",
				src:  filepath.Join(dirPath, "1"),
				dst:  filepath.Join(dirPath, "3"),
				err:  "error in readDirNames",
			},
			{
				name: "err in second readDirName",
				src:  filepath.Join(dirPath, "3"),
				dst:  filepath.Join(dirPath, "1"),
				err:  "error in readDirNames",
			},
			{
				name: "directory content is different",
				src:  filepath.Join(dirPath, "3"),
				dst:  filepath.Join(dirPath, "4"),
				err:  fmt.Sprintf("directory contents did not match:\n%q had %v\n%q had %v", filepath.Join(dirPath, "3"), srcEntries, filepath.Join(dirPath, "4"), dstEntries),
			},
			{
				name: "tree content is different",
				src:  filepath.Join(dirPath, "4"),
				dst:  filepath.Join(dirPath, "5"),
				err:  "could not get the stat for src or dst",
			},
		}

		for _, tt := range testTable3 {
			err := IsEqualTree(equalTreeOpts, tt.src, tt.dst)

			if fmt.Sprintf("%v", err) != tt.err {
				t.Errorf("Test %s: got: (%s), want: (%s)", tt.name, err, tt.err)
			}
		}
	})

	// Symlink
	// Creating Symlinks and adapt the opts symlink value
	equalTreeOpts.NoFollowSymlinks = true
	err = os.Symlink(filepath.Join(dirPath, "1", "1"), filepath.Join(dirPath, "1", "symlink1"))
	if err != nil {
		t.Errorf("err while creating a symlink")
	}
	err = os.Symlink(filepath.Join(dirPath, "1", "2"), filepath.Join(dirPath, "1", "symlink2"))
	if err != nil {
		t.Errorf("err while creating a symlink")
	}
	err = os.Symlink(filepath.Join(dirPath, "1", "3"), filepath.Join(dirPath, "1", "symlink3"))
	if err != nil {
		t.Errorf("err while creating a symlink")
	}
	err = os.Symlink(filepath.Join(dirPath, "1", "1"), filepath.Join(dirPath, "1", "symlink4"))
	if err != nil {
		t.Errorf("got: (%s), want: (%s)", err.Error(), "")
	}

	srcTarget, err := os.Readlink(filepath.Join(dirPath, "1", "symlink3"))
	if err != nil {
		t.Errorf("err is: %v", err)
	}
	dstTarget, err := os.Readlink(filepath.Join(dirPath, "1", "symlink4"))
	if err != nil {
		t.Errorf("err is: %v", err)
	}

	var testTable4 = []struct {
		name string
		src  string
		dst  string
		err  string
	}{
		{
			name: "symlinks are not equal",
			src:  filepath.Join(dirPath, "1", "symlink3"),
			dst:  filepath.Join(dirPath, "1", "symlink4"),
			err: fmt.Sprintf("target mismatch: symlink %q had target %q, while %q had target %q", filepath.Join(dirPath, "1", "symlink3"),
				srcTarget, filepath.Join(dirPath, "1", "symlink4"), dstTarget),
		},
		{
			name: "symlinks are equal",
			src:  filepath.Join(dirPath, "1", "symlink3"),
			dst:  filepath.Join(dirPath, "1", "symlink3"),
			err:  "<nil>",
		},
	}

	// Testing IsEqualTree for case symlink
	t.Run("Test IsEqualTree for case symlink", func(t *testing.T) {
		for _, tt := range testTable4 {
			err := IsEqualTree(equalTreeOpts, tt.src, tt.dst)

			if fmt.Sprintf("%v", err) != tt.err {
				t.Errorf("Test %s: got: (%s), want: (%s)", tt.name, err, tt.err)
			}
		}
	})

	var testTable5 = []struct {
		name string
		src  string
		dst  string
		err  string
	}{
		{
			name: "first read link err",
			src:  filepath.Join(dirPath, "1", "symlink1"),
			dst:  filepath.Join(dirPath, "1", "symlink2"),
			err:  "error in readlink",
		},
		{
			name: "second read link err",
			src:  filepath.Join(dirPath, "1", "symlink3"),
			dst:  filepath.Join(dirPath, "1", "symlink2"),
			err:  "error in readlink",
		},
	}

	// Testing IsEqualTree for case symlink errors
	t.Run("Test IsEqualTree for case symlink errors", func(t *testing.T) {
		for _, tt := range testTable5 {
			err := IsEqualTree(equalTreeOpts, tt.src, tt.dst)

			if fmt.Sprintf("%v", err) != tt.err {
				t.Errorf("Test %s: got: (%s), want: (%s)", tt.name, err, tt.err)
			}
		}
	})

	// Testing  IsEqualTree case regular file
	t.Run("Test IsEqualTree case regular file", func(t *testing.T) {
		err = IsEqualTree(equalTreeOpts, filepath.Join(dirPath, "5", "4"), filepath.Join(dirPath, "5", "4"))
		if err != nil {
			t.Errorf("err is: %v", err)
		}
	})

}
