// Copyright 2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cp

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"syscall"
	"testing"
)

func copyAndTest(t *testing.T, o Options, src, dst string) {
	if err := o.Copy(src, dst); err != nil {
		t.Fatalf("Copy(%q -> %q) = %v, want %v", src, dst, err, nil)
	}
	// if err := cmp.IsEqualTree(o, src, dst); err != nil {
	// 	t.Fatalf("Expected %q and %q to be same, got %v", src, dst, err)
	// }
}

func TestSimpleCopy(t *testing.T) {
	tmpDir := t.TempDir()

	// Copy a directory.
	origd := filepath.Join(tmpDir, "directory")
	if err := os.Mkdir(origd, 0o744); err != nil {
		t.Fatal(err)
	}

	copyAndTest(t, Default, origd, filepath.Join(tmpDir, "directory-copied"))
	copyAndTest(t, NoFollowSymlinks, origd, filepath.Join(tmpDir, "directory-copied-2"))

	// Copy a file.
	origf := filepath.Join(tmpDir, "normal-file")
	if err := os.WriteFile(origf, []byte("F is for fire that burns down the whole town"), 0o766); err != nil {
		t.Fatal(err)
	}

	copyAndTest(t, Default, origf, filepath.Join(tmpDir, "normal-file-copied"))
	copyAndTest(t, NoFollowSymlinks, origf, filepath.Join(tmpDir, "normal-file-copied-2"))

	// Copy a symlink.
	origs := filepath.Join(tmpDir, "foobar")
	// foobar -> normal-file
	if err := os.Symlink(origf, origs); err != nil {
		t.Fatal(err)
	}

	copyAndTest(t, Default, origf, filepath.Join(tmpDir, "foobar-copied"))
	copyAndTest(t, NoFollowSymlinks, origf, filepath.Join(tmpDir, "foobar-copied-just-symlink"))
}

func TestCopyTree(t *testing.T) {
	testfiles := make([]*os.File, 3)
	tmpDir, err := ioutil.TempDir("", "u-root-pkg-cp-")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// Make src directory.
	srcd := filepath.Join(tmpDir, "src")
	if err := os.Mkdir(srcd, 0744); err != nil {
		t.Fatal(err)
	}

	if err := os.Chdir(srcd); err != nil {
		t.Fatal(err)
	}

	// Make some rnd files
	for i := 0; i < 3; i++ {
		testfiles[i], err = os.Create("testfile" + fmt.Sprintf("%d", i))
		if err != nil {
			t.Fatal(err)
		}
	}

	// Make dest directory.
	dest := filepath.Join(tmpDir, "dest")
	if err := os.Mkdir(dest, 0744); err != nil {
		t.Fatal(err)
	}
	// Copy the tree
	if err := CopyTree(srcd, dest); err != nil {
		t.Fatal(err)
	}
}

func TestCopyFile(t *testing.T) {
	// Defining files and vars for the test
	var equalTreeOpts Options

	firstTmpDir, err := os.MkdirTemp("", "u-root-pkg-cp-")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(firstTmpDir)

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

	equalTreeOpts.NoFollowSymlinks = true
	err = os.Symlink(tmpFile1.Name(), filepath.Join(firstTmpDir, "symlink1"))
	if err != nil {
		t.Errorf("err while creating a symlink")
	}

	// Retrieving os.ModeSymlink from symlink
	srcInfo, err := equalTreeOpts.stat(firstTmpDir + "/symlink1")
	if err != nil {
		t.Errorf("err is: %v", err)
	}

	// Read the symlink
	err = copyFile(firstTmpDir+"/symlink1", firstTmpDir+"/symlink2", srcInfo)
	if err != nil {
		t.Errorf("err is: %v", err)
	}

	// Error in reading symlink
	err = copyFile(tmpFile1.Name(), firstTmpDir+"/symlink2", srcInfo)
	if fmt.Sprintf("%v", err) != "readlink "+tmpFile1.Name()+": invalid argument" {
		t.Errorf("Test %s: got: (%s), want: (%s)", "error in os.Readlink", err, "readlink "+tmpFile1.Name()+": invalid argument")
	}
}

func TestCopyRegularFile(t *testing.T) {

	// Faking os.Open function
	oopen := open
	defer func() { open = oopen }()
	open = func(name string) (*os.File, error) {
		if name == "srcto" {
			return nil, fmt.Errorf("error in open src")
		}
		f, err := os.OpenFile(name, syscall.O_RDONLY, 0)
		if err != nil {
			return nil, err
		}
		defer f.Close()
		return f, nil
	}

	err := copyRegularFile("srcto", "dst", nil)
	if fmt.Sprintf("%v", err) != "error in open src" {
		t.Errorf("Test %s: got: (%s), want: (%s)\n", "error open src", err, "error in open src")
	}
}
