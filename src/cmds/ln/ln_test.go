// Copyright 2016 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"testing"
)

const (
	// if true remove the test creation on the end
	remove   = true
	testPath = "."
)

func resetFlags() {
	flags.symlink = false
	flags.verbose = false
	flags.force = false
	flags.nondir = false
	flags.prompt = false
	flags.logical = false
	flags.physical = false
	flags.relative = false
	flags.dirtgt = ""
}

// create a temp file
func newFile(path, testName string, t *testing.T) (f *os.File) {
	f, err := ioutil.TempFile(path, "Go_"+testName)
	if err != nil {
		t.Fatalf("TempFile %s: %s", testName, err)
	}

	return
}

// create a temp dir
func newDir(testName string, t *testing.T) (name string) {
	name, err := ioutil.TempDir(testPath, "Go_"+testName)
	if err != nil {
		t.Fatalf("TempDir %s: %s", testName, err)
	}
	return
}

// test if hardlink crealinkNamen was sucessful
// 'from' and 'linkName' must exists
func testHardLink(from, linkName string, t *testing.T) {
	linkStat, err := os.Stat(linkName)
	if err != nil {
		t.Fatalf("stat %q failed: %v", linkName, err)
	}
	fromstat, err := os.Stat(from)
	if err != nil {
		t.Fatalf("stat %q failed: %v", from, err)
	}
	if !os.SameFile(linkStat, fromstat) {
		t.Errorf("link %q, %q did not create hard link", linkName, from)
	}
}

// test if symlink creation was sucessful
// 'from' and 'linkName' must exists
func testSymlink(target, linkName string, t *testing.T) {
	linkStat, err := os.Stat(linkName)
	if err != nil {
		t.Fatalf("stat %q failed: %v", linkName, err)
	}
	targetStat, err := os.Stat(target)
	if err != nil {
		t.Fatalf("stat %q failed: %v", target, err)
	}
	if !os.SameFile(linkStat, targetStat) {
		t.Errorf("symlink %q, %q did not create symlink", linkName, target)
	}
	targetStat, err = os.Stat(target)
	if err != nil {
		t.Fatalf("lstat %q failed: %v", target, err)
	}

	if targetStat.Mode()&os.ModeSymlink == os.ModeSymlink {
		t.Fatalf("symlink %q, %q did not create symlink", linkName, target)
	}

	targetStat, err = os.Stat(target)
	if err != nil {
		t.Fatalf("stat %q failed: %v", target, err)
	}
	if targetStat.Mode()&os.ModeSymlink != 0 {
		t.Fatalf("stat %q did not follow symlink", target)
	}
	s, err := os.Readlink(linkName)
	if err != nil {
		t.Fatalf("readlink %q failed: %v", target, err)
	}
	if s != filepath.Base(target) {
		t.Fatalf("after symlink %q != %q", s, target)
	}
	file, err := os.Open(target)
	if err != nil {
		t.Fatalf("open %q failed: %v", target, err)
	}
	file.Close()
}

// Ln default behavior (without flags)
// cmd-line equivalent: $ ln target link_name
func TestLnHardLink2Args(t *testing.T) {
	d := newDir("TestLnHardLink2Arg", t)
	if remove {
		defer os.RemoveAll(d)
	}

	f := newFile(d, "target", t)

	target := f.Name()
	linkName := target + "_hardlink"

	args := []string{target, linkName}
	if err := ln(args); err != nil {
		t.Errorf("Ln execution fails: %s", err)
	}

	testHardLink(target, linkName, t)

}

// cmd-line equivalent: $ ln -s target link_name
func TestLnSymlink2Arg(t *testing.T) {
	flags.symlink = true
	defer resetFlags()

	d := newDir("TestLnSymlink2Arg", t)
	if remove {
		defer os.RemoveAll(d)
	}

	f := newFile(d, "target", t)

	target := f.Name()
	linkName := path.Join(filepath.Dir(target), "symlink")
	targetLocal := filepath.Base(target)

	// i get the base because the test will be in dir 'd'
	// so, we must have: linkName -> fileNameTarget
	// [lerax@starfox Go_TestLnSymlink2Arg415244648 (ln-dev)]$ ls -l
	// total 0
	// -rw------- 1 lerax users  0 Feb 15 18:34 Go_target154139047
	// lrwxrwxrwx 1 lerax users 18 Feb 15 18:34 Go_target154139047_symlink -> Go_target154139047
	// like calling $ln -s Go_TestLnSymlink2Arg415244648/Go_target
	args := []string{targetLocal, linkName}
	if err := ln(args); err != nil {
		t.Errorf("Ln execution fails: %s", err)
	}

	//rigorous test about symlink validation
	testSymlink(target, linkName, t)
}
