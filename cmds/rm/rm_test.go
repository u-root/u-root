// Copyright 2012 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
//  created by Manoel Vilela (manoel_vilela@engineer.com)

package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"testing"
)

type makeit struct {
	n string      // name
	m os.FileMode // mode
	s string      // for symlinks or content
}

/* Tests contains standard file permissions */
var tests = []makeit{
	{
		n: "hi1.txt",
		m: 0666,
		s: "",
	},
	{
		n: "hi2.txt",
		m: 0777,
		s: "",
	},
	{
		n: "go.txt",
		m: 0555,
		s: "",
	},
}

/* Tests2 contains files, directories, links, and block devices */
var tests2 = []makeit{
	{
		n: "file.txt",
		m: 0000,
		s: "",
	},
	{
		n: "dir",
	},
	{
		n: "dirWfiles",
	},
	{
		n: "block",
	},
	{
		n: "symFile",
	},
}

func setup() (string, error) {
	fmt.Println(":: Creating simulating data...")
	d, err := ioutil.TempDir(os.TempDir(), "hi.dir")
	if err != nil {
		return "", err
	}

	/* Creates some files with content */
	for i := range tests {
		if err := ioutil.WriteFile(path.Join(d, tests[i].n), []byte("Go is cool!"), tests[i].m); err != nil {
			return "", err
		}
	}

	/* Creates empty file and directory s */
	if err := ioutil.WriteFile(path.Join(d, tests2[0].n), nil, tests2[0].m); err != nil {
		return "", err
	}
	if ioutil.TempDir(d, tests2[1].n); err != nil {
		return "", err
	}

	/* Create file, symlink, and block device in a new directory */
	newD, err := ioutil.TempDir(d, tests2[2].n)
	if err != nil {
		return "", err
	}
	if err := ioutil.WriteFile(path.Join(newD, tests2[0].n), []byte(""), tests2[0].m); err != nil {
		return "", err
	}

	if err := os.Symlink(path.Join(d, tests2[0].n), path.Join(newD, "symlink")); err != nil {
		return "", err
	}
	/* Does not work on Travis
	if err := syscall.Mknod(path.Join(newD, tests2[3].n), 0777, 64); err != nil {
		return "", err
	}
	*/
	return d, nil
}

func printFiles(d string) (filenames []string, err error) {
	var nameArray []string
	fileNames, err := ioutil.ReadDir(d)
	if err != nil {
		return nil, err
	}
	for _, file := range fileNames {
		fmt.Println(file.Name())
		nameArray = append(nameArray, path.Join(d, file.Name()))
	}
	return nameArray, nil
}

// TEST 1
// Flags: -v; only delete files from tests

func Test_rm_1(t *testing.T) {
	fmt.Println("TEST 1:")
	d, err := setup()
	if err != nil {
		t.Fatalf("Error on setup of the test: creating files and folders: %s", err)
	}
	fmt.Println("== Deleting files and empty folders (no args) ...")
	files := []string{path.Join(d, "hi1.txt"), path.Join(d, "hi2.txt"), path.Join(d, "go.txt")}
	var flags rmFlags
	flags.verbose = true
	if err := rm(files, flags); err != nil {
		t.Error(err)
	}
	os.RemoveAll(d)
}

// TEST 2
// Flags: -v -r; only delete files from test1
func Test_rm_2(t *testing.T) {
	fmt.Println("TEST 2:")
	d, err := setup()
	if err != nil {
		t.Fatalf("Error on setup of the test: creating files and folders.")
	}
	var flags rmFlags
	flags.verbose = true
	flags.recursive = true
	fmt.Println("== Deleting folders recursively (using -r flag) ...")
	files := []string{d}
	if err := rm(files, flags); err != nil {
		t.Error(err)
	}
	defer os.RemoveAll(d)
}

// TEST 3
// Flags: -r none; delete files from test1 and 2 and print all files in directory
func Test_rm_3(t *testing.T) {
	fmt.Println("TEST 3:")
	d, err := setup()
	if err != nil {
		t.Fatalf("Error on setup of the test: creating files and folders.")
	}
	fmt.Printf("All files in directory")
	filename, err := printFiles(d)
	if err != nil {
		t.Error(err)
	}
	fmt.Printf("end of all files")
	//defer os.RemoveAll(d)
	fmt.Printf("%v", filename)
	var flags rmFlags
	flags.recursive = true
	if err := rm(filename, flags); err != nil {
		t.Error(err)
	}
	printFiles(d)
	defer os.RemoveAll(d)
}

//TEST 4
/*
//does not work with Travis, but is still a good test
func Test_rm_4(t *testing.T) {
	fmt.Println("TEST 4:")
	d, err := setup()
	if err != nil {
		t.Fatalf("Error on setup of the test: creating files and folders.")
	}
	var flags rmFlags
	if err := rm([]string{"dnefile"}, flags); err != nil {
		t.Error(err)
	}
	printFiles(d)

}*/
