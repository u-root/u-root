// Copyright 2016 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// created by Manoel Vilela <manoel_vilela@engineer.com>

package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/u-root/u-root/pkg/cp"
	"github.com/u-root/u-root/pkg/cp/cmp"
)

const (
	maxSizeFile = 1000
	maxDirDepth = 5
	maxFiles    = 5
)

// resetFlags is used to reset the cp flags to default
func resetFlags() {
	flags.recursive = false
	flags.ask = false
	flags.force = false
	flags.verbose = false
	flags.noFollowSymlinks = false
}

// randomFile create a random file with random content
func randomFile(fpath, prefix string) (*os.File, error) {
	f, err := ioutil.TempFile(fpath, prefix)
	if err != nil {
		return nil, err
	}
	// generate random content for files
	bytes := []byte{}
	for i := 0; i < rand.Intn(maxSizeFile); i++ {
		bytes = append(bytes, byte(i))
	}
	f.Write(bytes)

	return f, nil
}

// createFilesTree create a random files tree
func createFilesTree(root string, maxDepth, depth int) error {
	// create more one dir if don't achieve the maxDepth
	if depth < maxDepth {
		newDir, err := ioutil.TempDir(root, fmt.Sprintf("cpdir_%d_", depth))
		if err != nil {
			return err
		}

		if err = createFilesTree(newDir, maxDepth, depth+1); err != nil {
			return err
		}
	}
	// generate random files
	for i := 0; i < maxFiles; i++ {
		f, err := randomFile(root, fmt.Sprintf("cpfile_%d_", i))
		if err != nil {
			return err
		}
		f.Close()
	}

	return nil
}

// TestCpsSimple make a simple test for copy file-to-file
// cmd-line equivalent: cp file file-copy
func TestCpSimple(t *testing.T) {
	tempDir, err := ioutil.TempDir("", "TestCpSimple")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	f, err := randomFile(tempDir, "src-")
	if err != nil {
		t.Fatalf("cannot create a random file: %v", err)
	}
	defer f.Close()

	srcf := f.Name()
	dstf := filepath.Join(tempDir, "destination")

	if err := cpArgs([]string{srcf, dstf}); err != nil {
		t.Fatalf("copy(%q -> %q) = %v, want nil", srcf, dstf, err)
	}
	if err := cmp.IsEqualTree(cp.Default, srcf, dstf); err != nil {
		t.Fatalf("copy(%q -> %q): file trees not equal: %v", srcf, dstf, err)
	}
}

// TestCpSrcDirectory tests copying source to destination without recursive
// cmd-line equivalent: cp ~/dir ~/dir2
func TestCpSrcDirectory(t *testing.T) {
	flags.recursive = false
	defer resetFlags()

	tempDir, err := ioutil.TempDir("", "TestCpSrcDirectory")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	tempDirTwo, err := ioutil.TempDir("", "TestCpSrcDirectoryTwo")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDirTwo)

	// capture log output to verify
	var logBytes bytes.Buffer
	log.SetOutput(&logBytes)

	if err := cpArgs([]string{tempDir, tempDirTwo}); err != nil {
		t.Fatalf("copy(%q -> %q) = %v, want nil", tempDir, tempDirTwo, err)
	}

	outString := fmt.Sprintf("cp: -r not specified, omitting directory %s", tempDir)
	capturedString := logBytes.String()
	if !strings.Contains(capturedString, outString) {
		t.Fatalf("copy(%q -> %q) = %v, want %v", tempDir, tempDirTwo, capturedString, outString)
	}
}

// TestCpRecursive tests the recursive mode copy
// Copy dir hierarchies src-dir to dst-dir
// whose src-dir and dst-dir already exists
// cmd-line equivalent: $ cp -R src-dir/ dst-dir/
func TestCpRecursive(t *testing.T) {
	flags.recursive = true
	defer resetFlags()

	tempDir, err := ioutil.TempDir("", "TestCpRecursive")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	srcDir := filepath.Join(tempDir, "src")
	if err := os.Mkdir(srcDir, 0755); err != nil {
		t.Fatal(err)
	}
	dstDir := filepath.Join(tempDir, "dst-exists")
	if err := os.Mkdir(dstDir, 0755); err != nil {
		t.Fatal(err)
	}

	if err = createFilesTree(srcDir, maxDirDepth, 0); err != nil {
		t.Fatalf("cannot create files tree on directory %q: %v", srcDir, err)
	}

	t.Run("existing-dst-dir", func(t *testing.T) {
		if err := cpArgs([]string{srcDir, dstDir}); err != nil {
			t.Fatalf("cp(%q -> %q) = %v, want nil", srcDir, dstDir, err)
		}
		// Because dstDir already existed, a new dir was created inside it.
		realDestination := filepath.Join(dstDir, filepath.Base(srcDir))
		if err := cmp.IsEqualTree(cp.Default, srcDir, realDestination); err != nil {
			t.Fatalf("copy(%q -> %q): file trees not equal: %v", srcDir, realDestination, err)
		}
	})

	t.Run("non-existing-dst-dir", func(t *testing.T) {
		notExistDstDir := filepath.Join(tempDir, "dst-does-not-exist")
		if err := cpArgs([]string{srcDir, notExistDstDir}); err != nil {
			t.Fatalf("cp(%q -> %q) = %v, want nil", srcDir, notExistDstDir, err)
		}

		if err := cmp.IsEqualTree(cp.Default, srcDir, notExistDstDir); err != nil {
			t.Fatalf("copy(%q -> %q): file trees not equal: %v", srcDir, notExistDstDir, err)
		}
	})
}

// Other test to verify the CopyRecursive
// whose dir$n and dst-dir already exists
// cmd-line equivalent: $ cp -R dir1/ dir2/ dir3/ dst-dir/
//
// dst-dir will content dir{1, 3}
// $ dst-dir/
// ..	dir1/
// ..	dir2/
// ..   dir3/
func TestCpRecursiveMultiple(t *testing.T) {
	flags.recursive = true
	defer resetFlags()
	tempDir, err := ioutil.TempDir("", "TestCpRecursiveMultiple")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	dstTest := filepath.Join(tempDir, "destination")
	if err := os.Mkdir(dstTest, 0755); err != nil {
		t.Fatalf("Failed on build directory %v: %v", dstTest, err)
	}

	// create multiple random directories sources
	srcDirs := []string{}
	for i := 0; i < maxDirDepth; i++ {
		srcTest, err := ioutil.TempDir(tempDir, "src-")
		if err != nil {
			t.Fatalf("Failed on build directory %v: %v\n", srcTest, err)
		}
		if err = createFilesTree(srcTest, maxDirDepth, 0); err != nil {
			t.Fatalf("cannot create files tree on directory %v: %v", srcTest, err)
		}

		srcDirs = append(srcDirs, srcTest)

	}

	args := srcDirs
	args = append(args, dstTest)
	if err := cpArgs(args); err != nil {
		t.Fatalf("cp %q exit with error: %v", args, err)
	}
	for _, src := range srcDirs {
		_, srcFile := filepath.Split(src)

		dst := filepath.Join(dstTest, srcFile)
		if err := cmp.IsEqualTree(cp.Default, src, dst); err != nil {
			t.Fatalf("The copy %q -> %q failed: %v", src, dst, err)
		}
	}
}

// using -P don't follow symlinks, create other symlink
// cmd-line equivalent: $ cp -P symlink symlink-copy
func TestCpSymlink(t *testing.T) {
	tempDir, err := ioutil.TempDir("", "TestCpSymlink")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	f, err := randomFile(tempDir, "src-")
	if err != nil {
		t.Fatalf("cannot create a random file: %v", err)
	}
	defer f.Close()

	srcFpath := f.Name()
	srcFname := filepath.Base(srcFpath)

	newName := filepath.Join(tempDir, srcFname+"_link")
	if err := os.Symlink(srcFname, newName); err != nil {
		t.Fatalf("cannot create a link %q with target %q: %v", newName, srcFname, err)
	}

	t.Run("no-follow-symlink", func(t *testing.T) {
		defer resetFlags()
		flags.noFollowSymlinks = true

		dst := filepath.Join(tempDir, "dst-no-follow")
		if err := cpArgs([]string{newName, dst}); err != nil {
			t.Fatalf("cp(%q -> %q) = %v, want nil", newName, dst, err)
		}
		if err := cmp.IsEqualTree(cp.NoFollowSymlinks, newName, dst); err != nil {
			t.Fatalf("The copy %q -> %q failed: %v", newName, dst, err)
		}
	})

	t.Run("follow-symlink", func(t *testing.T) {
		defer resetFlags()
		flags.noFollowSymlinks = false

		dst := filepath.Join(tempDir, "dst-follow")
		if err := cpArgs([]string{newName, dst}); err != nil {
			t.Fatalf("cp(%q -> %q) = %v, want nil", newName, dst, err)
		}
		if err := cmp.IsEqualTree(cp.Default, newName, dst); err != nil {
			t.Fatalf("The copy %q -> %q failed: %v", newName, dst, err)
		}
	})
}
