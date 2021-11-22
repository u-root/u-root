// Copyright 2021 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cmp

import (
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/u-root/u-root/pkg/cp"
)

func TestCMP(t *testing.T) {
	// Creating all tmp dirs and files for testing purpose
	firstTmpDir, err := os.MkdirTemp("", "u-root-pkg-cmp-")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(firstTmpDir)

	secondTmpDir, err := os.MkdirTemp("", "u-root-pkg-cmp-")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(secondTmpDir)

	thirdTmpDir, err := os.MkdirTemp("", "u-root-pkg-cmp-")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(thirdTmpDir)

	fourthTmpDir, err := os.MkdirTemp("", "u-root-pkg-cmp-")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(fourthTmpDir)

	fifthTmpDir, err := os.MkdirTemp("", "u-root-pkg-cmp-")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(fifthTmpDir)

	sixthTmpDir, err := os.MkdirTemp("", "u-root-pkg-cmp-")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(sixthTmpDir)

	seventhTmpDir, err := os.MkdirTemp(fifthTmpDir, "u-root-pkg-cmp-")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(seventhTmpDir)

	tmpFile1, err := ioutil.TempFile(firstTmpDir, "file1")
	if err != nil {
		t.Fatal(err)
	}
	defer tmpFile1.Close()

	tmpFile2, err := ioutil.TempFile(firstTmpDir, "file2")
	if err != nil {
		t.Fatal(err)
	}
	defer tmpFile2.Close()

	tmpFile3, err := ioutil.TempFile(firstTmpDir, "file3")
	if err != nil {
		t.Fatal(err)
	}
	defer tmpFile3.Close()

	tmpFile4, err := ioutil.TempFile(fifthTmpDir, "file4")
	if err != nil {
		t.Fatal(err)
	}
	defer tmpFile4.Close()

	if err := os.WriteFile(tmpFile1.Name(), []byte("F is for fire that burns down the whole town"), 0766); err != nil {
		t.Fatal(err)
	}

	if err := os.WriteFile(tmpFile2.Name(), []byte("F is for fire that burns down the whole town"), 0766); err != nil {
		t.Fatal(err)
	}

	if err := os.WriteFile(tmpFile3.Name(), []byte("nwot elohw eht nwod snrub taht erif rof si F"), 0766); err != nil {
		t.Fatal(err)
	}

	if err := os.WriteFile(tmpFile4.Name(), []byte("nwot elohw eht nwod snrub taht erif rof si F"), fs.ModeExclusive); err != nil {
		t.Fatal(err)
	}

	// Tests start here
	//
	//
	// Struct for testing isEqualFile and readDirNames
	var test1 = []struct {
		n     string
		file1 string
		file2 string
		err   string
	}{
		{n: "file1 does not exist", file1: "file1", file2: tmpFile2.Name(), err: "open file1: no such file or directory"},
		{n: "file2 does not exist", file1: tmpFile1.Name(), file2: "file2", err: "open file2: no such file or directory"},
		{n: "files are not equal", file1: tmpFile1.Name(), file2: tmpFile3.Name(), err: fmt.Sprintf("%q and %q do not have equal content", tmpFile1.Name(), tmpFile3.Name())},
	}

	// Testing isEqualFile
	t.Run("Test isEqualFile", func(t *testing.T) {
		for _, tt := range test1 {
			err := isEqualFile(tt.file1, tt.file2)
			if err.Error() != tt.err {
				t.Errorf("Test %s: got: (%s), want: (%s)", tt.n, err.Error(), tt.err)
			}
		}
		err = isEqualFile(tmpFile1.Name(), tmpFile2.Name())
		if err != nil {
			t.Errorf("got: (%s), want: (%s)", err.Error(), "")
		}
	})

	// Testing readDirNames
	t.Run("Test isEqualFile", func(t *testing.T) {
		names, err := readDirNames(firstTmpDir)
		if len(names) != 3 || names[0] != filepath.Base(tmpFile1.Name()) || names[1] != filepath.Base(tmpFile2.Name()) || names[2] != filepath.Base(tmpFile3.Name()) || err != nil {
			t.Errorf("file amount: %d, files: %v, files created %s, %s, %s", len(names), names, filepath.Base(tmpFile1.Name()), filepath.Base(tmpFile2.Name()), filepath.Base(tmpFile3.Name()))
		}
		_, err = readDirNames("dir")
		if err.Error() != "open dir: no such file or directory" {
			t.Errorf("got: (%s), want: (%s)", err.Error(), "")
		}
	})

	// Testing stat
	t.Run("Test stat", func(t *testing.T) {
		statOpts := cp.Default
		_, err = stat(statOpts, tmpFile1.Name())
		if err != nil {
			t.Fatal(err)
		}
		statOpts.NoFollowSymlinks = true
		_, err = stat(statOpts, tmpFile1.Name())
		if err != nil {
			t.Fatal(err)
		}
	})

	// Default option var
	equalTreeOpts := cp.Default

	// Testing isEqualTree
	t.Run("Test isEqualTree stat fail + dirs are equal", func(t *testing.T) {
		// Struct for testing isEqualTree
		var test2 = []struct {
			n   string
			src string
			dst string
			err string
		}{
			{n: "stat src err", src: tmpFile1.Name(), dst: "", err: "stat : no such file or directory"},
			{n: "stat dst err", src: "", dst: tmpFile2.Name(), err: "stat : no such file or directory"},
			{n: "two dirs are equal", src: fourthTmpDir, dst: sixthTmpDir, err: "<nil>"},
		}

		for _, tt := range test2 {
			_, _, _, err := stats(equalTreeOpts, tt.src, tt.dst)

			if fmt.Sprintf("%v", err) != tt.err {
				t.Errorf("Test %s: got: (%s), want: (%s)", tt.n, err, tt.err)
			}
		}
	})

	t.Run("Test stat", func(t *testing.T) {
		// Test case that two dirs are equal
		err := IsEqualTree(equalTreeOpts, fourthTmpDir, sixthTmpDir)
		if fmt.Sprintf("%v", err) != "<nil>" {
			t.Errorf("Test %s: got: (%s), want: (%s)\n", "two dirs are equal", err, "<nil>")
		}

		// Faking readDirNames function
		oReadDirName := readDirName
		defer func() { readDirName = oReadDirName }()
		readDirName = func(path string) ([]string, error) {
			if path == firstTmpDir {
				return nil, fmt.Errorf("error in readDirNames")
			}
			if path == secondTmpDir {
				return nil, fmt.Errorf("error in readDirNames")
			}
			var basename = []string{"test1", "test2"}
			if path == thirdTmpDir {
				basename[0] = "test3"
			}
			return basename, nil
		}

		// retrieve sm and dm for err checking
		sm, dm, _, err := stats(equalTreeOpts, firstTmpDir, tmpFile2.Name())
		if err != nil {
			t.Errorf("err is: %v", err)
		}

		// retrive srcEntries and dstEntries
		srcEntries, err := readDirName(thirdTmpDir)
		if err != nil {
			t.Errorf("err is: %v", err)
		}
		dstEntries, err := readDirName(fourthTmpDir)
		if err != nil {
			t.Errorf("err is: %v", err)
		}

		// Struct for testing isEqualTree
		var test3 = []struct {
			n   string
			src string
			dst string
			err string
		}{
			{n: "mismatched mode, one dir one file", src: firstTmpDir, dst: tmpFile2.Name(), err: fmt.Sprintf("mismatched mode: %q has mode %s while %q has mode %s", firstTmpDir, sm, tmpFile2.Name(), dm)},
			{n: "err in first readDirName", src: firstTmpDir, dst: thirdTmpDir, err: "error in readDirNames"},
			{n: "err in second readDirName", src: thirdTmpDir, dst: secondTmpDir, err: "error in readDirNames"},
			{n: "directory content is different", src: thirdTmpDir, dst: fourthTmpDir, err: fmt.Sprintf("directory contents did not match:\n%q had %v\n%q had %v", thirdTmpDir, srcEntries, fourthTmpDir, dstEntries)},
			{n: "tree content is different", src: fourthTmpDir, dst: fifthTmpDir, err: "could not get the stat for src or dst"},
		}

		for _, tt := range test3 {
			err := IsEqualTree(equalTreeOpts, tt.src, tt.dst)

			if fmt.Sprintf("%v", err) != tt.err {
				t.Errorf("Test %s: got: (%s), want: (%s)", tt.n, err, tt.err)
			}
		}
	})

	// Symlink

	// Creating Symlinks
	equalTreeOpts.NoFollowSymlinks = true
	err = os.Symlink(tmpFile1.Name(), filepath.Join(firstTmpDir, "symlink1"))
	if err != nil {
		t.Errorf("err while creating a symlink")
	}
	err = os.Symlink(tmpFile2.Name(), filepath.Join(firstTmpDir, "symlink2"))
	if err != nil {
		t.Errorf("err while creating a symlink")
	}
	err = os.Symlink(tmpFile3.Name(), filepath.Join(firstTmpDir, "symlink3"))
	if err != nil {
		t.Errorf("err while creating a symlink")
	}
	err = os.Symlink(tmpFile1.Name(), filepath.Join(firstTmpDir, "symlink4"))
	if err != nil {
		t.Errorf("err while creating a symlink")
	}

	srcTarget, err := readLink(firstTmpDir + "/symlink3")
	if err != nil {
		t.Errorf("err is: %v", err)
	}
	dstTarget, err := readLink(firstTmpDir + "/symlink4")
	if err != nil {
		t.Errorf("err is: %v", err)
	}

	var test4 = []struct {
		n   string
		src string
		dst string
		err string
	}{
		{n: "symlinks are not equal", src: firstTmpDir + "/symlink3", dst: firstTmpDir + "/symlink4", err: fmt.Sprintf("target mismatch: symlink %q had target %q, while %q had target %q", firstTmpDir+"/symlink3", srcTarget, firstTmpDir+"/symlink4", dstTarget)},
		{n: "symlinks are equal", src: firstTmpDir + "/symlink3", dst: firstTmpDir + "/symlink3", err: "<nil>"},
	}

	for _, tt := range test4 {
		err := IsEqualTree(equalTreeOpts, tt.src, tt.dst)

		if fmt.Sprintf("%v", err) != tt.err {
			t.Errorf("Test %s: got: (%s), want: (%s)", tt.n, err, tt.err)
		}
	}

	// Fake the readLink func
	oReadLink := readLink
	defer func() { readLink = oReadLink }()
	readLink = func(name string) (string, error) {
		if name == firstTmpDir+"/symlink1" {
			return "", fmt.Errorf("error in readlink")
		} else if name == firstTmpDir+"/symlink2" {
			return "", fmt.Errorf("error in readlink")
		}
		return "test", nil
	}

	var test5 = []struct {
		n   string
		src string
		dst string
		err string
	}{
		{n: "first read link err", src: firstTmpDir + "/symlink1", dst: firstTmpDir + "/symlink2", err: "error in readlink"},
		{n: "second read link err", src: firstTmpDir + "/symlink3", dst: firstTmpDir + "/symlink2", err: "error in readlink"},
	}

	for _, tt := range test5 {
		err := IsEqualTree(equalTreeOpts, tt.src, tt.dst)

		if fmt.Sprintf("%v", err) != tt.err {
			t.Errorf("Test %s: got: (%s), want: (%s)", tt.n, err, tt.err)
		}
	}

	err = IsEqualTree(equalTreeOpts, tmpFile4.Name(), tmpFile4.Name())
	if err != nil {
		t.Errorf("err is: %w", err)
	}
}
