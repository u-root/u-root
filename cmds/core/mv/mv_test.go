// Copyright 2015-2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"os"
	"path/filepath"
	"syscall"
	"testing"
	"time"

	"github.com/u-root/u-root/pkg/testutil"
)

type makeIt struct {
	n string      // name
	m os.FileMode // mode
	c []byte      // content
}

var hiFileContent = []byte("hi")

var old = makeIt{
	n: "old.txt",
	m: 0o777,
	c: []byte("old"),
}

var new = makeIt{
	n: "new.txt",
	m: 0o777,
	c: []byte("new"),
}

var tests = []makeIt{
	{
		n: "hi1.txt",
		m: 0o666,
		c: hiFileContent,
	},
	{
		n: "hi2.txt",
		m: 0o777,
		c: hiFileContent,
	},
	old,
	new,
}

func setup(t *testing.T) (string, error) {
	d := t.TempDir()

	tmpdir := filepath.Join(d, "hi.sub.dir")
	if err := os.Mkdir(tmpdir, 0o777); err != nil {
		return "", err
	}

	for _, tt := range tests {
		if err := os.WriteFile(filepath.Join(d, tt.n), tt.c, tt.m); err != nil {
			return "", err
		}
	}

	return d, nil
}

func getInode(file string) (uint64, error) {
	var stat syscall.Stat_t
	if err := syscall.Stat(file, &stat); err != nil {
		return 0, err
	}
	return stat.Ino, nil
}

func TestMv(t *testing.T) {
	d, err := setup(t)
	if err != nil {
		t.Fatal("err")
	}
	defer os.RemoveAll(d)

	t.Logf("Renaming file...")
	{
		originalInode, err := getInode(filepath.Join(d, "hi1.txt"))
		if err != nil {
			t.Error(err)
		}

		files := []string{filepath.Join(d, "hi1.txt"), filepath.Join(d, "hi4.txt")}
		res := testutil.Command(t, files...)
		_, err = res.CombinedOutput()
		if err = testutil.IsExitCode(err, 0); err != nil {
			t.Error(err)
		}

		t.Logf("Verify renamed file integrity...")
		{
			content, err := os.ReadFile(filepath.Join(d, "hi4.txt"))
			if err != nil {
				t.Error(err)
			}

			if !bytes.Equal(hiFileContent, content) {
				t.Errorf("Expected file content to equal %s, got %s", hiFileContent, content)
			}

			movedInode, err := getInode(filepath.Join(d, "hi4.txt"))
			if err != nil {
				t.Error(err)
			}

			if originalInode != movedInode {
				t.Errorf("Expected inode to equal. Expected %d, got %d", originalInode, movedInode)
			}
		}
	}

	dsub := filepath.Join(d, "hi.sub.dir")

	t.Logf("Moving files to directory...")
	{
		originalInode, err := getInode(filepath.Join(d, "hi2.txt"))
		if err != nil {
			t.Error(err)
		}

		originalInodeFour, err := getInode(filepath.Join(d, "hi4.txt"))
		if err != nil {
			t.Error(err)
		}

		files := []string{filepath.Join(d, "hi2.txt"), filepath.Join(d, "hi4.txt"), dsub}
		res := testutil.Command(t, files...)
		_, err = res.CombinedOutput()
		if err = testutil.IsExitCode(err, 0); err != nil {
			t.Error(err)
		}

		t.Logf("Verify moved files into directory file integrity...")
		{
			content, err := os.ReadFile(filepath.Join(dsub, "hi4.txt"))
			if err != nil {
				t.Error(err)
			}

			if !bytes.Equal(hiFileContent, content) {
				t.Errorf("Expected file content to equal %s, got %s", hiFileContent, content)
			}

			movedInode, err := getInode(filepath.Join(dsub, "hi2.txt"))
			if err != nil {
				t.Error(err)
			}

			movedInodeFour, err := getInode(filepath.Join(dsub, "hi4.txt"))
			if err != nil {
				t.Error(err)
			}

			if originalInode != movedInode {
				t.Errorf("Expected inode to equal. Expected %d, got %d", originalInode, movedInode)
			}

			if originalInodeFour != movedInodeFour {
				t.Errorf("Expected inode to equal. Expected %d, got %d", originalInodeFour, movedInodeFour)
			}
		}
	}
}

func TestMvUpdate(t *testing.T) {
	*update = true
	d, err := setup(t)
	if err != nil {
		t.Error(err)
	}
	defer os.RemoveAll(d)
	t.Logf("Testing mv -u...")

	// Ensure that the newer file actually has a newer timestamp
	currentTime := time.Now().Local()
	oldTime := currentTime.Add(-10 * time.Second)
	err = os.Chtimes(filepath.Join(d, old.n), oldTime, oldTime)
	if err != nil {
		t.Error(err)
	}
	err = os.Chtimes(filepath.Join(d, new.n), currentTime, currentTime)
	if err != nil {
		t.Error(err)
	}

	// Check that it doesn't downgrade files with -u switch
	{
		files := []string{"-u", filepath.Join(d, old.n), filepath.Join(d, new.n)}
		res := testutil.Command(t, files...)
		_, err = res.CombinedOutput()
		if err = testutil.IsExitCode(err, 0); err != nil {
			t.Error(err)
		}
		newContent, err := os.ReadFile(filepath.Join(d, new.n))
		if err != nil {
			t.Error(err)
		}
		if bytes.Equal(newContent, old.c) {
			t.Error("Newer file was overwritten by older file. Should not happen with -u.")
		}
	}

	// Check that it does update files with -u switch
	{
		files := []string{"-u", filepath.Join(d, new.n), filepath.Join(d, old.n)}
		res := testutil.Command(t, files...)
		_, err = res.CombinedOutput()
		if err = testutil.IsExitCode(err, 0); err != nil {
			t.Error(err)
		}
		newContent, err := os.ReadFile(filepath.Join(d, old.n))
		if err != nil {
			t.Error(err)
		}
		if !bytes.Equal(newContent, new.c) {
			t.Error("Older file was not overwritten by newer file. Should happen with -u.")
		}
		if _, err := os.Lstat(filepath.Join(d, old.n)); err != nil {
			t.Error("The new file shouldn't be there anymore.")
		}
	}
}

func TestMvNoClobber(t *testing.T) {
	*noClobber = true
	d, err := setup(t)
	if err != nil {
		t.Error(err)
	}
	defer os.RemoveAll(d)
	t.Logf("Testing mv -n...")

	// Check that it doesn't override files with -n switch
	{
		files := []string{"-n", filepath.Join(d, old.n), filepath.Join(d, new.n)}
		res := testutil.Command(t, files...)
		_, err = res.CombinedOutput()
		if err = testutil.IsExitCode(err, 0); err != nil {
			t.Error(err)
		}
		newContent, err := os.ReadFile(filepath.Join(d, new.n))
		if err != nil {
			t.Error(err)
		}
		if bytes.Equal(newContent, old.c) {
			t.Error("File was overwritten. Should not happen with -u.")
		}
	}

	// Check that it does mv files with -u switch
	{
		files := []string{"-n", filepath.Join(d, new.n), filepath.Join(d, "hi3.txt")}
		res := testutil.Command(t, files...)
		_, err = res.CombinedOutput()
		if err = testutil.IsExitCode(err, 0); err != nil {
			t.Error(err)
		}
		newContent, err := os.ReadFile(filepath.Join(d, "hi3.txt"))
		if err != nil {
			t.Error(err)
		}
		if !bytes.Equal(newContent, new.c) {
			t.Error("File was not moved. Should happen with -u.")
		}
		if _, err := os.Lstat(filepath.Join(d, old.n)); err != nil {
			t.Error("File was copied but not moved.")
		}
	}
}

func TestMain(m *testing.M) {
	testutil.Run(m, main)
}
