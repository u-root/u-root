// Copyright 2016 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// created by Manoel Vilela <manoel_vilela@engineer.com>

package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

const (
	testPath   = "."
	removeTest = true
)

type result struct {
	symlink       bool
	name, linksTo string
}

type test struct {
	conf    config   // conf with flags
	args    []string // to pass for ln
	results []result // expected results
	cmdline string   // cmdline ln equivalent
}

func loadTests() []test {
	return []test{
		{
			// covers usage:
			// ln [OPTIONS]... [-T] TARGET LINK_NAME   (1st form) (posix)
			config{},
			[]string{"a", "b"},
			[]result{
				result{name: "b", linksTo: "a"},
			},
			"$ ln a b",
		},
		{
			config{symlink: true},
			[]string{"a", "b"},
			[]result{
				{symlink: true, name: "b", linksTo: "a"},
			},
			"$ ln -s a b",
		},
		{
			// covers usage:
			// ln [OPTIONS]... TARGET   (2nd form) (gnu)
			config{symlink: true},
			[]string{"/bin/cp"},
			[]result{
				{symlink: true, name: "cp", linksTo: "/bin/cp"},
			},
			"$ ln -s /usr/bin/cp",
		},
		{
			// covers usage:
			// ln [OPTIONS]... TARGET... DIRECTORY  (3rd form) (posix)
			config{symlink: true},
			[]string{"/bin/cp", "/bin/ls", "/bin/ln", "."},
			[]result{
				{symlink: true, name: "cp", linksTo: "/bin/cp"},
				{symlink: true, name: "ls", linksTo: "/bin/ls"},
				{symlink: true, name: "ln", linksTo: "/bin/ln"},
			},
			"$ ln -s /bin/cp /bin/ls /bin/ln .",
		},
		{
			// covers usage:
			// ln [OPTIONS]... -t DIRECTORY TARGET...  (4th form) (gnu)
			config{symlink: true, dirtgt: "folder"},
			[]string{"/bin/cp", "/bin/ls", "/bin/ln"},
			[]result{
				{symlink: true, name: "folder/cp", linksTo: "/bin/cp"},
				{symlink: true, name: "folder/ls", linksTo: "/bin/ls"},
				{symlink: true, name: "folder/ln", linksTo: "/bin/ln"},
			},
			"$ ln -s -v -t folder /bin/cp /bin/ls /bin/ln",
		},

		{
			// -i -f mutually exclusive (f overwrite evers)
			config{force: true, prompt: true},
			[]string{"a", "overwrite"},
			[]result{
				{name: "overwrite", linksTo: "a"},
			},
			"$ ln -i -f a overwrite",
		},
	}
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
// 'target' and 'linkName' must exists
// linkName -> target
func testHardLink(linkName, target string, t *testing.T) {
	linkStat, err := os.Stat(linkName)
	if err != nil {
		t.Fatalf("stat %q failed: %v", linkName, err)
	}
	targetStat, err := os.Stat(target)
	if err != nil {
		t.Fatalf("stat %q failed: %v", target, err)
	}
	if !os.SameFile(linkStat, targetStat) {
		t.Errorf("link %q, %q did not create hard link", linkName, target)
	}
}

// test if symlink creation was sucessful
// 'target' and 'linkName' must exists
// linkName -> target
func testSymlink(linkName, target string, t *testing.T) {
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
	if s != target {
		t.Fatalf("after symlink %q != %q", s, target)
	}
	file, err := os.Open(target)
	if err != nil {
		t.Fatalf("open %q failed: %v", target, err)
	}
	file.Close()
}

// Alpha test using tabdriven
// Generic evaluation
func TestLn(t *testing.T) {
	tabDriven := loadTests()
	testDir := newDir("TestLnGeneric", t)
	if removeTest {
		defer os.RemoveAll(testDir)
	}

	for caseNum, testCase := range tabDriven {
		fis, err := ioutil.ReadDir(testDir)
		if err != nil {
			t.Fatal(err)
		}
		for _, fi := range fis {
			os.Remove(filepath.Join(testDir, fi.Name()))
		}

		if testCase.conf.dirtgt != "" {
			d := filepath.Join(testDir, testCase.conf.dirtgt)
			if err := os.Mkdir(d, 0750); err != nil {
				t.Errorf("Creation of target dir %q fails: %v", d, err)
				continue
			}
		}

		if err := ioutil.WriteFile(filepath.Join(testDir, "a"), []byte{'.'}, 0640); err != nil {
			t.Fatal(err)
		}

		// write a file to overwrite
		if testCase.conf.force {
			if err := ioutil.WriteFile(filepath.Join(testDir, "overwrite"), []byte{'.'}, 0640); err != nil {
				t.Fatal(err)
			}
		}

		// executing ln on isolated tempDir
		if err := os.Chdir(testDir); err != nil {
			t.Errorf("Changing directory for %q fails: %v", testDir, err)
		}
		if err := testCase.conf.ln(testCase.args); err != nil {
			t.Errorf("test [%d]. %v", caseNum, err)
			continue
		}

		t.Logf("Testing cmdline: %q", testCase.cmdline)
		for _, expected := range testCase.results {
			if expected.symlink {
				t.Logf("%q -> %q (symlink)", expected.name, expected.linksTo)
				testSymlink(expected.name, expected.linksTo, t)
			} else {
				t.Logf("%q -> %q (hardlink)", expected.name, expected.linksTo)
				testHardLink(expected.name, expected.linksTo, t)
			}
		}
		// underfolder
		if err := os.Chdir(".."); err != nil {
			t.Errorf("Backing to root fails: %v", err)
		}
	}
}
