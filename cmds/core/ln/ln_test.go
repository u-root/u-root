// Copyright 2016 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// created by Manoel Vilela <manoel_vilela@engineer.com>

package main

import (
	"os"
	"path/filepath"
	"testing"
)

type create struct {
	name string
	dir  bool
}

type result struct {
	symlink       bool
	name, linksTo string
}

type test struct {
	conf    config   // conf with flags
	args    []string // to pass for ln
	results []result // expected results
	files   []create // previous files for testing
	cmdline string   // cmdline ln equivalent
}

// loadTests loads the main table driven tests
// for ln command tests
func loadTests() []test {
	return []test{
		{
			// covers usage:
			// ln [OPTIONS]... [-T] TARGET LINK_NAME   (1st form) (posix)
			config{},
			[]string{"a", "b"},
			[]result{{name: "b", linksTo: "a"}},
			[]create{{name: "a"}},
			"ln a b",
		},
		{
			config{symlink: true},
			[]string{"a", "b"},
			[]result{{symlink: true, name: "b", linksTo: "a"}},
			[]create{{name: "a"}},
			"ln -s a b",
		},
		{
			// covers usage:
			// ln [OPTIONS]... TARGET   (2nd form) (gnu)
			config{symlink: true},
			[]string{"bin/cp"},
			[]result{
				{symlink: true, name: "cp", linksTo: "bin/cp"},
			},
			[]create{
				{name: "bin", dir: true},
				{name: "bin/cp"},
			},
			"ln -s bin/cp",
		},
		{
			// covers usage:
			// ln [OPTIONS]... TARGET... DIRECTORY  (3rd form) (posix)
			config{symlink: true},
			[]string{"bin/cp", "bin/ls", "bin/ln", "."},
			[]result{
				{symlink: true, name: "cp", linksTo: "bin/cp"},
				{symlink: true, name: "ls", linksTo: "bin/ls"},
				{symlink: true, name: "ln", linksTo: "bin/ln"},
			},
			[]create{
				{name: "bin", dir: true},
				{name: "bin/cp"},
				{name: "bin/ls"},
				{name: "bin/ln"},
			},
			"ln -s bin/cp bin/ls bin/ln .",
		},
		{
			// covers usage:
			// ln [OPTIONS]... -t DIRECTORY TARGET...  (4th form) (gnu)
			config{symlink: true, dirtgt: "."},
			[]string{"bin/cp", "bin/ls", "bin/ln"},
			[]result{
				{symlink: true, name: "cp", linksTo: "bin/cp"},
				{symlink: true, name: "ls", linksTo: "bin/ls"},
				{symlink: true, name: "ln", linksTo: "bin/ln"},
			},
			[]create{
				{name: "bin", dir: true},
				{name: "bin/cp"},
				{name: "bin/ls"},
				{name: "bin/ln"},
			},
			"ln -s bin/cp bin/ls bin/ln -t .",
		},
		{
			// covers usage:
			// ln [OPTIONS]... -t DIRECTORY TARGET...  (4th form) (gnu)
			config{symlink: true, dirtgt: "folder", relative: true},
			[]string{"cp", "ls", "ln"},
			[]result{
				{symlink: true, name: "folder/cp", linksTo: "../cp"},
				{symlink: true, name: "folder/ls", linksTo: "../ls"},
				{symlink: true, name: "folder/ln", linksTo: "../ln"},
			},
			[]create{
				{name: "folder", dir: true},
				{name: "cp"},
				{name: "ls"},
				{name: "ln"},
			},
			"ln -s -v -r -t folder cp ls ln",
		},
		{
			// -i -f mutually exclusive (f overwrite evers)
			config{force: true, prompt: true},
			[]string{"a", "overwrite"},
			[]result{
				{name: "overwrite", linksTo: "a"},
			},
			[]create{
				{name: "overwrite"},
				{name: "a"},
			},
			"ln -i -f a overwrite",
		},
	}
}

// testHardLink test if hardlink creation was successful
// 'target' and 'linkName' must exists
// linkName -> target
func testHardLink(linkName, target string, t *testing.T) {
	linkStat, err := os.Stat(linkName)
	if err != nil {
		t.Errorf("stat %q failed: %v", linkName, err)
	}
	targetStat, err := os.Stat(target)
	if err != nil {
		t.Errorf("stat %q failed: %v", target, err)
	}
	if !os.SameFile(linkStat, targetStat) {
		t.Errorf("link %q, %q did not create hard link", linkName, target)
	}
}

// testSymllink test if symlink creation was successful
// 'target' and 'linkName' must exists
// linkName -> target
func testSymlink(linkName, linksTo string, t *testing.T) {
	target := linksTo
	if !filepath.IsAbs(target) {
		target = filepath.Base(target)
	}

	linkStat, err := os.Stat(linkName)
	if err != nil {
		t.Errorf("stat %q failed: %v", linkName, err)
	}
	targetStat, err := os.Stat(target)
	if err != nil {
		t.Errorf("stat %q failed: %v", target, err)
	}
	if !os.SameFile(linkStat, targetStat) {
		t.Errorf("symlink %q, %q did not create symlink", linkName, target)
	}
	targetStat, err = os.Stat(target)
	if err != nil {
		t.Errorf("lstat %q failed: %v", target, err)
	}

	if targetStat.Mode()&os.ModeSymlink == os.ModeSymlink {
		t.Errorf("symlink %q, %q did not create symlink", linkName, target)
	}

	targetStat, err = os.Stat(target)
	if err != nil {
		t.Errorf("stat %q failed: %v", target, err)
	}
	if targetStat.Mode()&os.ModeSymlink != 0 {
		t.Errorf("stat %q did not follow symlink", target)
	}
	s, err := os.Readlink(linkName)
	if err != nil {
		t.Errorf("readlink %q failed: %v", target, err)
	}
	if s != linksTo {
		t.Errorf("after symlink %q != %q", s, target)
	}
	file, err := os.Open(target)
	if err != nil {
		t.Errorf("open %q failed: %v", target, err)
	}
	file.Close()
}

// TestLn make a general tests based on
// tabDriven tests (see loadTests())
func TestLn(t *testing.T) {
	tabDriven := loadTests()
	testDir := t.TempDir()

	// executing ln on isolated testDir
	if err := os.Chdir(testDir); err != nil {
		t.Fatalf("Changing directory for %q fails: %v", testDir, err)
	}
	defer os.Chdir("..") // after defer to go back to the original root

	for caseNum, testCase := range tabDriven {
		d := t.TempDir()
		if err := os.Chdir(d); err != nil {
			t.Fatalf("Changing directory for %q fails: %v", d, err)
		}

		for _, f := range testCase.files {
			t.Logf("Creating: %v (dir: %v)", f.name, f.dir)
			p := filepath.Join(f.name)
			if f.dir {
				if err := os.Mkdir(p, 0o750); err != nil && err == os.ErrExist {
					t.Skipf("Creation of dir %q fails: %v", p, err)
				}
			} else {
				if err := os.WriteFile(p, []byte{'.'}, 0o640); err != nil {
					t.Fatal(err)
				}
			}

		}

		if err := testCase.conf.ln(testCase.args); err != nil {
			t.Errorf("Fails: test [%d]. %v", caseNum+1, err)
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

		// backing to testDir folder
		os.Chdir("..")
	}
}
