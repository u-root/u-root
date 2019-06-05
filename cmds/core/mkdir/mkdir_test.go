// Copyright 2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"
	"syscall"
	"testing"

	"github.com/u-root/u-root/pkg/testutil"
)

type MkdirErrorTestcase struct {
	name string
	args []string
	exp  string
}

type MkdirPermTestCase struct {
	name     string
	args     []string
	perm     os.FileMode
	dirNames []string
}

var (
	stubDirNames   = []string{"stub", "stub2"}
	umaskDefault   = 022
	errorTestCases = []MkdirErrorTestcase{
		{
			name: "No Arg Error",
			args: nil,
			exp:  "Usage",
		},
		{
			name: "Perm Mode Bits over 7 Error",
			args: []string{"-m=7778", stubDirNames[0]},
			exp:  `invalid mode '7778'`,
		},
		{
			name: "More than 4 Perm Mode Bits Error",
			args: []string{"-m=11111", stubDirNames[0]},
			exp:  `invalid mode '11111'`,
		},
	}

	regularTestCases = []struct {
		name string
		args []string
		exp  []string
	}{
		{
			name: "Create 1 Directory",
			args: []string{stubDirNames[0]},
			exp:  []string{stubDirNames[0]},
		},
		{
			name: "Create 2 Directories",
			args: stubDirNames,
			exp:  stubDirNames,
		},
	}

	permTestCases = []MkdirPermTestCase{
		{
			name:     "Default Perm",
			args:     []string{stubDirNames[0]},
			perm:     os.FileMode(0755 | os.ModeDir),
			dirNames: []string{stubDirNames[0]},
		},
		{
			name:     "Custom Perm in Octal Form",
			args:     []string{"-m=0777", stubDirNames[0]},
			perm:     os.FileMode(0777 | os.ModeDir),
			dirNames: []string{stubDirNames[0]},
		},
		{
			name:     "Custom Perm not in Octal Form",
			args:     []string{"-m=777", stubDirNames[0]},
			perm:     os.FileMode(0777 | os.ModeDir),
			dirNames: []string{stubDirNames[0]},
		},
		{
			name:     "Custom Perm with Sticky Bit",
			args:     []string{"-m=1777", stubDirNames[0]},
			perm:     os.FileMode(0777 | os.ModeDir | os.ModeSticky),
			dirNames: []string{stubDirNames[0]},
		},
		{
			name:     "Custom Perm with SGID Bit",
			args:     []string{"-m=2777", stubDirNames[0]},
			perm:     os.FileMode(0777 | os.ModeDir | os.ModeSetgid),
			dirNames: []string{stubDirNames[0]},
		},
		{
			name:     "Custom Perm with SUID Bit",
			args:     []string{"-m=4777", stubDirNames[0]},
			perm:     os.FileMode(0777 | os.ModeDir | os.ModeSetuid),
			dirNames: []string{stubDirNames[0]},
		},
		{
			name:     "Custom Perm with Sticky Bit and SUID Bit",
			args:     []string{"-m=5777", stubDirNames[0]},
			perm:     os.FileMode(0777 | os.ModeDir | os.ModeSticky | os.ModeSetuid),
			dirNames: []string{stubDirNames[0]},
		},
		{
			name:     "Custom Perm for 2 Directories",
			args:     []string{"-m=5777", stubDirNames[0], stubDirNames[1]},
			perm:     os.FileMode(0777 | os.ModeDir | os.ModeSticky | os.ModeSetuid),
			dirNames: stubDirNames,
		},
	}
)

func run(c *exec.Cmd) (string, string, error) {
	var o, e bytes.Buffer
	c.Stdout, c.Stderr = &o, &e
	err := c.Run()
	return o.String(), e.String(), err
}

func printError(t *testing.T, testname string, execStmt string, out interface{}, exp interface{}) {
	t.Logf("TEST %v", testname)
	t.Errorf("%s\ngot:%v\nwant:%v", execStmt, out, exp)
}

func findFile(dir string, filename string) (os.FileInfo, error) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	for _, file := range files {
		if file.Name() == filename {
			return file, nil
		}
	}
	return nil, nil
}

func removeCreatedFiles(tmpDir string) {
	for _, dirName := range stubDirNames {
		os.Remove(filepath.Join(tmpDir, dirName))
	}
}

func TestMkdirErrors(t *testing.T) {
	// Set Up
	tmpDir, err := ioutil.TempDir("", "ls")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)
	syscall.Umask(umaskDefault)

	// Error Tests
	for _, test := range errorTestCases {
		removeCreatedFiles(tmpDir)
		c := testutil.Command(t, test.args...)
		execStmt := fmt.Sprintf("exec(mkdir %s)", strings.Trim(fmt.Sprint(test.args), "[]"))
		c.Dir = tmpDir
		_, e, err := run(c)
		if err == nil || !strings.Contains(e, test.exp) {
			printError(t, test.name, execStmt, e, test.exp)
			continue
		}
		f, err := findFile(tmpDir, stubDirNames[0])
		if err != nil {
			printError(t, test.name, execStmt, err, "No error while finding the file")
			continue
		}
		if f != nil {
			printError(t, test.name, execStmt, "A directory was created", "No directory should be created")
		}
	}
}

func TestMkdirRegular(t *testing.T) {
	// Set Up
	tmpDir, err := ioutil.TempDir("", "ls")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)
	syscall.Umask(umaskDefault)

	// Regular Tests
	for _, test := range regularTestCases {
		removeCreatedFiles(tmpDir)
		c := testutil.Command(t, test.args...)
		execStmt := fmt.Sprintf("exec(mkdir %s)", strings.Trim(fmt.Sprint(test.args), "[]"))
		c.Dir = tmpDir
		_, e, err := run(c)
		if err != nil {
			printError(t, test.name, execStmt, e, "No error while mkdir")
			continue
		}
		for _, dirName := range test.exp {
			f, err := findFile(tmpDir, dirName)
			if err != nil {
				printError(t, test.name, execStmt, err, "No error while finding the file")
				break
			}
			if f == nil {
				printError(t, test.name, execStmt, fmt.Sprintf("%s not found", dirName), fmt.Sprintf("%s should have been created", dirName))
				break
			}
		}
	}
}

func TestMkdirPermission(t *testing.T) {
	// Set Up
	tmpDir, err := ioutil.TempDir("", "ls")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)
	syscall.Umask(umaskDefault)

	// Permission Tests
	for _, test := range permTestCases {
		removeCreatedFiles(tmpDir)
		c := testutil.Command(t, test.args...)
		execStmt := fmt.Sprintf("exec(mkdir %s)", strings.Trim(fmt.Sprint(test.args), "[]"))
		c.Dir = tmpDir
		_, e, err := run(c)
		if err != nil {
			printError(t, test.name, execStmt, e, "No error while mkdir")
			continue
		}
		for _, dirName := range test.dirNames {
			f, err := findFile(tmpDir, dirName)
			if err != nil {
				printError(t, test.name, execStmt, err, "No error while finding the file")
				break
			}
			if f == nil {
				printError(t, test.name, execStmt, fmt.Sprintf("%s not found", dirName), fmt.Sprintf("%s should have been created", dirName))
				break
			}
			if f != nil && !reflect.DeepEqual(f.Mode(), test.perm) {
				printError(t, test.name, execStmt, f.Mode(), test.perm)
				break
			}
		}
	}
}

func TestMain(m *testing.M) {
	testutil.Run(m, main)
}
