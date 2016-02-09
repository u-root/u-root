// Copyright 2016 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"crypto/md5"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"os"
	"path"
	"reflect"
	"testing"
	"time"
)

var (
	testPath    = "."
	remove      = true                            // if true removeAll the testPath on the end
	indentifier = time.Now().Format(time.Kitchen) // simple label for src and dst
)

const (
	maxDirDepth = 5
	maxFiles    = 5
	maxSizeFile = 1024
	rangeRand   = 100
)

// reset the cp flags to default
func resetFlags() {
	nwork = 1
	recursive = false
	ask = false
	force = false
	verbose = false
	link = false
}

// create a random file with random content
func randomFile(fpath, prefix string) (*os.File, error) {
	f, err := ioutil.TempFile(fpath, prefix)
	if err != nil {
		return nil, err
	}
	// generate random content for files
	bytes := []byte{}
	for i := 0; i < rand.Intn(maxSizeFile); i++ {
		bytes = append(bytes, 'a'+(byte((rand.Intn(rangeRand)))%26))
	}
	f.Write(bytes)

	return f, err
}

// create a random files tree
func createTreeFiles(root string, maxDepth, depth int) error {
	// create more one dir if don't achieve the maxDepth
	if depth < maxDepth {
		newDir, err := ioutil.TempDir(root, fmt.Sprintf("cpdir_%d_", depth))
		if err != nil {
			return err
		}

		if err = createTreeFiles(newDir, maxDepth, depth+1); err != nil {
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

// get the path and fname each file of dir pathToRead
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
		fpath := path.Join(pathToRead, fname)
		if !exists {
			slc := []string{fpath}
			mapFiles[fname] = slc
		} else {
			mapFiles[fname] = append(mapFiles[fname], fpath)
		}
	}
	return err
}

// get a file and return your md5sum
func getMd5Sum(file *os.File) ([]byte, error) {
	hash := md5.New()
	_, err := io.Copy(hash, file)
	if err != nil {
		return nil, err
	}

	return hash.Sum(nil), nil

}

// compare two files by checksum
func isEqualFile(f1, f2 *os.File, t *testing.T) bool {
	checksum1, err := getMd5Sum(f1)
	if err != nil {
		t.Errorf("Failed to compute the md5sum of file %v: %v", f1, err)
		return false
	}
	checksum2, err := getMd5Sum(f2)
	if err != nil {
		t.Errorf("Failed to compute the md5sum of file %v: %v", f2, err)
		return false
	}
	t.Logf("source file %v get md5sum %v", f1.Name(), checksum1)
	t.Logf("destination file %v get md5sum %v", f2.Name(), checksum2)
	return reflect.DeepEqual(checksum1, checksum2)
}

// compare the content between of src and dst paths
func isEqualTree(src, dst string, t *testing.T) bool {
	mapFiles := map[string][]string{}
	if err := readDirs(src, mapFiles); err != nil {
		t.Fatalf("cannot read dir %v: %v", src, err)
	}
	if err := readDirs(dst, mapFiles); err != nil {
		t.Fatalf("cannot read dir %v: %v", dst, err)
	}

	equalTree := true
	for _, files := range mapFiles {
		if len(files) < 2 {
			return false
		}
		fpath1, fpath2 := files[0], files[1]
		file1, err := os.Open(fpath1)
		if err != nil {
			t.Fatalf("cannot open file %v: %v", fpath1, err)
		}
		file2, err := os.Open(fpath2)
		if err != nil {
			t.Fatalf("cannot open file %v: %v", fpath2, err)
		}

		stat1, err := file1.Stat()
		if err != nil {
			t.Fatalf("cannot stat file %v: %v\n", file1, err)
		}
		stat2, err := file2.Stat()
		if err != nil {
			t.Fatalf("cannot stat file %v: %v\n", file2, err)

		}
		if stat1.IsDir() && stat2.IsDir() {
			equalDirs := isEqualTree(fpath1, fpath2, t)
			equalTree = equalTree && equalDirs
			if !equalDirs {
				t.Logf("Dirs %v %v: %v", fpath1, fpath2, equalDirs)
			}

		} else {
			equalFiles := isEqualFile(file1, file2, t)
			equalTree = equalTree && equalFiles
			if !equalTree {
				t.Logf("File %v == %v: %v", file1, file2, equalFiles)
			}

		}
		if !equalTree {
			break
		}
		file1.Close()
		file2.Close()

	}

	return equalTree

}

// make a simple test for copy file-to-file
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
	dstFpath := path.Join(tempDir, dstFname)

	if err := copyFile(srcFpath, dstFpath, false); err != nil {
		t.Fatalf("copyFile %v -> %v failed: %v", dstFpath, dstFpath, err)
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
	if !isEqualFile(s, d, t) {
		t.Fatalf("checksum are different; copies failed %v -> %v", srcFpath, dstFpath)
	}
}

// Test the recursive mode copy
// Copy dir hierarchies src-dir to dst-dir
// whose src-dir and dst-dir already exists
// cmd-line equivalent: $ cp -R src-dir/ dst-dir/
func TestCpRecursive(t *testing.T) {
	recursive = true
	defer resetFlags()
	tempDir, err := ioutil.TempDir(testPath, "TestCpRecursive")
	if err != nil {
		t.Fatal("Failed on build tmp dir %v: %v\n", testPath, err)
	}
	if remove {
		defer os.RemoveAll(tempDir)
	}
	srcPrefix := fmt.Sprintf("TestCpSrc_%v_", indentifier)
	dstPrefix := fmt.Sprintf("TestCpDst_%v_copied", indentifier)
	srcTest, err := ioutil.TempDir(tempDir, srcPrefix)
	if err != nil {
		t.Fatal("Failed on build directory %v: %v\n", srcTest, err)
	}
	if err = createTreeFiles(srcTest, maxDirDepth, 0); err != nil {
		t.Fatalf("cannot create files tree on directory %v: %v", srcTest)
	}

	dstTest, err := ioutil.TempDir(tempDir, dstPrefix)
	if err != nil {
		t.Fatalf("Failed on build directory %v: %v\n", dstTest, err)
	}
	if err := copyFile(srcTest, dstTest, false); err != nil {
		t.Fatalf("copyFile %v %v failed: %v", srcTest, dstTest, err)
	}

	if !isEqualTree(srcTest, dstTest, t) {
		t.Fatalf("The copy %v -> %v failed, trees are different", srcTest, dstTest)
	}
}

// whose src-dir exists but dst-dir no
// cmd-line equivalent: $ cp -R some-dir/ new-dir/
func TestCpRecursiveNew(t *testing.T) {
	recursive = true
	defer resetFlags()
	tempDir, err := ioutil.TempDir(testPath, "TestCpRecursiveNew")
	if err != nil {
		t.Fatal("failed on build tmp directory at %v: %v\n", tempDir, err)
	}
	if remove {
		defer os.RemoveAll(tempDir)
	}
	srcPrefix := fmt.Sprintf("TestCpSrc_%v_", indentifier)
	dstPrefix := fmt.Sprintf("TestCpDst_%v_new", indentifier)
	srcTest, err := ioutil.TempDir(tempDir, srcPrefix)
	if err != nil {
		t.Fatal("failed on build tmp directory %v: %v\n", srcPrefix, err)
	}

	if err = createTreeFiles(srcTest, maxDirDepth, 0); err != nil {
		t.Fatalf("cannot create files tree on directory %v: %v", srcTest)
	}

	dstTest := path.Join(tempDir, dstPrefix)
	copyFile(srcTest, dstTest, false)
	isEqual := isEqualTree(srcTest, dstTest, t)
	if !isEqual {
		t.Fatalf("The copy %v -> %v failed, ", srcTest, dstTest)
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
		t.Fatal("Failed on build tmp directory %v: %v\n", testPath, err)
	}
	if remove {
		defer os.RemoveAll(tempDir)
	}

	dstPrefix := fmt.Sprintf("TestCpDst_%v_container", indentifier)
	dstTest, err := ioutil.TempDir(tempDir, dstPrefix)
	if err != nil {
		t.Fatal("Failed on build directory %v: %v\n", dstTest, err)
	}
	// create multiple random directories sources
	srcDirs := []string{}
	for i := 0; i < maxDirDepth; i++ {
		srcPrefix := fmt.Sprintf("TestCpSrc_%v_", indentifier)
		srcTest, err := ioutil.TempDir(tempDir, srcPrefix)
		if err != nil {
			t.Fatal("Failed on build directory %v: %v\n", srcTest, err)
		}
		if err = createTreeFiles(srcTest, maxDirDepth, 0); err != nil {
			t.Fatalf("cannot create files tree on directory %v: %v", srcTest)
		}

		srcDirs = append(srcDirs, srcTest)

	}

	args := srcDirs
	args = append(args, dstTest)
	if err := cp(args); err != nil {
		t.Fatalf("cp %v exit with error: %v", args, err)
	}
	for _, src := range srcDirs {
		_, srcFile := path.Split(src)
		dst := path.Join(dstTest, srcFile)
		isEqual := isEqualTree(src, dst, t)
		if !isEqual {
			t.Fatalf("The copy %v -> %v failed, trees are different", src, dst)
		}
	}
}
