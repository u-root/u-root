// Copyright 2016 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// created by Manoel Vilela <manoel_vilela@engineer.com>

package main

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"path/filepath"
	"reflect"
	"testing"
	"time"
)

var (
	testPath = "."
	// if true removeAll the testPath on the end
	remove = true
	// simple label for src and dst
	indentifier = time.Now().Format(time.Kitchen)
)

const (
	maxSizeFile = buffSize * 2
	rangeRand   = 100
	maxDirDepth = 5
	maxFiles    = 5
)

// resetFlags is used to reset the cp flags to default
func resetFlags() {
	nwork = 1
	recursive = false
	ask = false
	force = false
	verbose = false
	symlink = false
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

// readDirs get the path and fname each file of dir pathToRead
// set on that structure fname:[fpath1, fpath2]
// use that to compare if the files with same name has same content
func readDirs(pathToRead string, mapFiles map[string][]string) error {
	pathFiles, err := ioutil.ReadDir(pathToRead)
	if err != nil {
		return err
	}
	for _, file := range pathFiles {
		fname := file.Name()
		_, exists := mapFiles[fname]
		fpath := filepath.Join(pathToRead, fname)
		if !exists {
			slc := []string{fpath}
			mapFiles[fname] = slc
		} else {
			mapFiles[fname] = append(mapFiles[fname], fpath)
		}
	}
	return err
}

// isEqualFile compare two files by checksum
func isEqualFile(f1, f2 *os.File) (bool, error) {
	bytes1, err := ioutil.ReadAll(f1)
	if err != nil {
		return false, fmt.Errorf("Failed to read the file %v: %v", f1, err)
	}
	bytes2, err := ioutil.ReadAll(f2)
	if err != nil {
		return false, fmt.Errorf("Failed to read the file %v: %v", f2, err)
	}

	if !reflect.DeepEqual(bytes1, bytes2) {
		return false, nil
	}

	return true, nil
}

// isEqualTree compare the content between of src and dst paths
func isEqualTree(src, dst string) (bool, error) {
	mapFiles := map[string][]string{}
	if err := readDirs(src, mapFiles); err != nil {
		return false, fmt.Errorf("cannot read dir %v: %v", src, err)
	}
	if err := readDirs(dst, mapFiles); err != nil {
		return false, fmt.Errorf("cannot read dir %v: %v", dst, err)
	}

	equalTree := true
	for _, files := range mapFiles {
		if len(files) < 2 {
			return false, fmt.Errorf("insufficient files in readDirs(): expected at least 2, got %v", len(files))
		}
		fpath1, fpath2 := files[0], files[1]
		file1, err := os.Open(fpath1)
		if err != nil {
			return false, fmt.Errorf("cannot open file %v: %v", fpath1, err)
		}
		file2, err := os.Open(fpath2)
		if err != nil {
			return false, fmt.Errorf("cannot open file %v: %v", fpath2, err)
		}

		stat1, err := file1.Stat()
		if err != nil {
			return false, fmt.Errorf("cannot stat file %v: %v", file1, err)
		}
		stat2, err := file2.Stat()
		if err != nil {
			return false, fmt.Errorf("cannot stat file %v: %v", file2, err)

		}
		if stat1.IsDir() && stat2.IsDir() {
			equalDirs, err := isEqualTree(fpath1, fpath2)
			if err != nil {
				return false, err
			}
			equalTree = equalTree && equalDirs

		} else {
			equalFiles, err := isEqualFile(file1, file2)
			if err != nil {
				return false, err
			}
			equalTree = equalTree && equalFiles

		}
		if !equalTree {
			break
		}
		file1.Close()
		file2.Close()

	}

	return equalTree, nil

}

// TestCpsSimple make a simple test for copy file-to-file
// cmd-line equivalent: cp file file-copy
func TestCpSimple(t *testing.T) {
	tempDir, err := ioutil.TempDir(testPath, "TestCpSimple")
	if remove {
		defer os.RemoveAll(tempDir)
	}
	srcPrefix := fmt.Sprintf("cpfile_%v_src", indentifier)
	f, err := randomFile(tempDir, srcPrefix)
	if err != nil {
		t.Fatalf("cannot create a random file: %v", err)
	}
	defer f.Close()
	srcFpath := f.Name()

	dstFname := fmt.Sprintf("cpfile_%v_dst_copied", indentifier)
	dstFpath := filepath.Join(tempDir, dstFname)

	if err := copyFile(srcFpath, dstFpath, false); err != nil {
		t.Fatalf("copyFile %v -> %v failed: %v", srcFpath, dstFpath, err)
	}
	s, err := os.Open(srcFpath)
	if err != nil {
		t.Fatalf("cannot open the file %v", srcFpath)
	}
	defer s.Close()
	d, err := os.Open(dstFpath)
	if err != nil {
		t.Fatalf("cannot open the file %v", dstFpath)
	}
	defer d.Close()
	if equal, err := isEqualFile(s, d); !equal || err != nil {
		t.Fatalf("checksum are different; copies failed %q -> %q: %v", srcFpath, dstFpath, err)
	}
}

// TestCpRecursive tests the recursive mode copy
// Copy dir hierarchies src-dir to dst-dir
// whose src-dir and dst-dir already exists
// cmd-line equivalent: $ cp -R src-dir/ dst-dir/
func TestCpRecursive(t *testing.T) {
	recursive = true
	defer resetFlags()
	tempDir, err := ioutil.TempDir(testPath, "TestCpRecursive")
	if err != nil {
		t.Fatalf("Failed on build tmp dir %q: %v\n", testPath, err)
	}
	if remove {
		defer os.RemoveAll(tempDir)
	}
	srcPrefix := fmt.Sprintf("TestCpSrc_%v_", indentifier)
	dstPrefix := fmt.Sprintf("TestCpDst_%v_copied", indentifier)
	srcTest, err := ioutil.TempDir(tempDir, srcPrefix)
	if err != nil {
		t.Fatalf("Failed on build directory %q: %v\n", srcTest, err)
	}
	if err = createFilesTree(srcTest, maxDirDepth, 0); err != nil {
		t.Fatalf("cannot create files tree on directory %q: %v", srcTest, err)
	}

	dstTest, err := ioutil.TempDir(tempDir, dstPrefix)
	if err != nil {
		t.Fatalf("Failed on build directory %q: %v\n", dstTest, err)
	}
	if err := copyFile(srcTest, dstTest, false); err != nil {
		t.Fatalf("copyFile %q -> %q failed: %v", srcTest, dstTest, err)
	}

	if equal, err := isEqualTree(srcTest, dstTest); !equal || err != nil {
		t.Fatalf("The copy %q -> %q failed, trees are different: %v", srcTest, dstTest, err)
	}
}

// whose src-dir exists but dst-dir no
// cmd-line equivalent: $ cp -R some-dir/ new-dir/
func TestCpRecursiveNew(t *testing.T) {
	recursive = true
	defer resetFlags()
	tempDir, err := ioutil.TempDir(testPath, "TestCpRecursiveNew")
	if err != nil {
		t.Fatalf("failed on build tmp directory at %v: %v\n", tempDir, err)
	}
	if remove {
		defer os.RemoveAll(tempDir)
	}
	srcPrefix := fmt.Sprintf("TestCpSrc_%v_", indentifier)
	dstPrefix := fmt.Sprintf("TestCpDst_%v_new", indentifier)
	srcTest, err := ioutil.TempDir(tempDir, srcPrefix)
	if err != nil {
		t.Fatalf("failed on build tmp directory %q: %v\n", srcPrefix, err)
	}

	if err = createFilesTree(srcTest, maxDirDepth, 0); err != nil {
		t.Fatalf("cannot create files tree on directory %q: %v", srcTest, err)
	}

	dstTest := filepath.Join(tempDir, dstPrefix)
	copyFile(srcTest, dstTest, false)
	isEqual, err := isEqualTree(srcTest, dstTest)
	if err != nil {
		t.Fatalf("The test isEqualTree failed")
	}
	if !isEqual && err != nil {
		t.Fatalf("The copy %q -> %q failed, ", srcTest, dstTest)
	}
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
	recursive = true
	defer resetFlags()
	tempDir, err := ioutil.TempDir(testPath, "TestCpRecursiveMultiple")
	if err != nil {
		t.Fatalf("Failed on build tmp directory %v: %v\n", testPath, err)
	}
	if remove {
		defer os.RemoveAll(tempDir)
	}

	dstPrefix := fmt.Sprintf("TestCpDst_%v_container", indentifier)
	dstTest, err := ioutil.TempDir(tempDir, dstPrefix)
	if err != nil {
		t.Fatalf("Failed on build directory %v: %v\n", dstTest, err)
	}
	// create multiple random directories sources
	srcDirs := []string{}
	for i := 0; i < maxDirDepth; i++ {
		srcPrefix := fmt.Sprintf("TestCpSrc_%v_", indentifier)
		srcTest, err := ioutil.TempDir(tempDir, srcPrefix)
		if err != nil {
			t.Fatalf("Failed on build directory %v: %v\n", srcTest, err)
		}
		if err = createFilesTree(srcTest, maxDirDepth, 0); err != nil {
			t.Fatalf("cannot create files tree on directory %v: %v", srcTest, err)
		}

		srcDirs = append(srcDirs, srcTest)

	}
	t.Logf("From: %q", srcDirs)
	t.Logf("To: %q", dstTest)
	args := srcDirs
	args = append(args, dstTest)
	if err := cp(args); err != nil {
		t.Fatalf("cp %q exit with error: %v", args, err)
	}
	for _, src := range srcDirs {
		_, srcFile := filepath.Split(src)
		dst := filepath.Join(dstTest, srcFile)
		if equal, err := isEqualTree(src, dst); !equal || err != nil {
			t.Fatalf("The copy %q -> %q failed, trees are different", src, dst)
		}
	}
}

// using -P don't follow symlinks, create other symlink
// cmd-line equivalent: $ cp -P symlink symlink-copy
func TestCpSymlink(t *testing.T) {
	defer resetFlags()
	symlink = true
	tempDir, err := ioutil.TempDir(testPath, "TestCpSymlink")
	if remove {
		defer os.RemoveAll(tempDir)
	}
	srcPrefix := fmt.Sprintf("cpfile_%v_origin", indentifier)
	f, err := randomFile(tempDir, srcPrefix)
	if err != nil {
		t.Fatalf("cannot create a random file: %v", err)
	}
	defer f.Close()

	srcFpath := f.Name()
	_, srcFname := filepath.Split(srcFpath)

	linkName := srcFname + "_link"
	t.Logf("Enter directory: %q", tempDir)
	// Enter to temp directory to don't have problems
	// with copy links and relative paths
	os.Chdir(tempDir)
	defer os.Chdir("..")
	defer t.Logf("Exiting directory: %q", tempDir)

	t.Logf("Create link %q -> %q", srcFname, linkName)
	if err := os.Symlink(srcFname, linkName); err != nil {
		t.Fatalf("cannot create a link %v -> %v: %v", srcFname, linkName, err)
	}

	dstFname := fmt.Sprintf("cpfile_%v_dst_link", indentifier)
	t.Logf("Copy from: %q", linkName)
	t.Logf("To: %q", dstFname)
	if err := copyFile(linkName, dstFname, false); err != nil {
		t.Fatalf("copyFile %q -> %q failed: %v", linkName, dstFname, err)
	}

	s, err := os.Open(linkName)
	if err != nil {
		t.Fatalf("cannot open the file %v", linkName)
	}
	defer s.Close()

	d, err := os.Open(dstFname)
	if err != nil {
		t.Fatalf("cannot open the file %v", dstFname)
	}
	defer d.Close()
	dStat, err := d.Stat()
	if err != nil {
		t.Fatalf("cannot stat file %v: %v\n", d, err)
	}
	if L := os.ModeSymlink; dStat.Mode()&L == L {
		t.Fatalf("destination file is not a link %v", d.Name())
	}

	if equal, err := isEqualFile(s, d); !equal || err != nil {
		t.Fatalf("checksum are different; copies failed %q -> %q: %v", linkName, dstFname, err)
	}
}
