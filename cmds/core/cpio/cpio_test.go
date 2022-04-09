// Copyright 2020 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/u-root/u-root/pkg/testutil"
)

type dirEnt struct {
	Name     string
	Type     string
	Content  string
	Target   string
	FileInfo os.FileInfo
}

func TestCpio(t *testing.T) {
	debug = t.Logf
	// Create a temporary directory
	tempDir := t.TempDir()

	targets := []dirEnt{
		{Name: "file1", Type: "file", Content: "Hello World"},
		{Name: "file2", Type: "file", Content: ""},
		{Name: "directory1", Type: "dir"},
	}
	if runtime.GOOS != "plan9" {
		targets = append(targets, []dirEnt{
			{Name: "hardlinked", Type: "hardlink", Target: "file1", Content: "Hello World"},
			{Name: "hardlinkedtofile2", Type: "hardlink", Target: "file2"},
		}...)
	}
	for _, ent := range targets {
		name := filepath.Join(tempDir, ent.Name)
		switch ent.Type {
		case "dir":
			err := os.Mkdir(name, os.FileMode(0o700))
			if err != nil {
				t.Fatalf("cannot create test directory: %v", err)
			}

		case "file":
			f, err := os.Create(name)
			if err != nil {
				t.Fatalf("cannot create test file: %v", err)
			}
			defer f.Close()
			_, err = f.WriteString(ent.Content)
			if err != nil {
				t.Fatal(err)
			}

		case "hardlink":
			target := filepath.Join(tempDir, ent.Target)
			err := os.Link(target, name)
			if err != nil {
				t.Fatalf("cannot create hard link: %v", err)
			}
		}
	}

	// Now that the temporary directory structure is complete, populate
	// the FileInfo for each target. This needs to happen in a second
	// pass because of the link count.
	for key, ent := range targets {
		var err error
		name := filepath.Join(tempDir, ent.Name)
		targets[key].FileInfo, err = os.Stat(name)
		if err != nil {
			t.Fatalf("cannot stat temporary dirent: %v", err)
		}
	}

	c := testutil.Command(t, "-v", "o")
	c.Dir = tempDir

	buffer := bytes.Buffer{}
	for _, ent := range targets {
		buffer.WriteString(ent.Name + "\n")
	}
	c.Stdin = &buffer

	archive, err := c.Output()
	if err != nil {
		t.Fatalf("%s %v", c.Stderr, err)
	}

	// Cpio can't read from a non-seekable input (e.g. a pipe) in input mode.
	// Write the archive to a file instead.
	archiveFile, err := os.CreateTemp("", "archive.cpio")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(archiveFile.Name())
	defer archiveFile.Close()

	if _, err := archiveFile.Write(archive); err != nil {
		t.Fatal(err)
	}

	// Extract to a new directory
	tempExtractDir := t.TempDir()

	c = testutil.Command(t, "-v", "i")
	c.Dir = tempExtractDir
	c.Stdin = archiveFile

	out, err := c.Output()
	if err != nil {
		t.Fatalf("Extraction failed:\n%s\n%s\n%v\n", out, c.Stderr, err)
	}

	for _, ent := range targets {
		name := filepath.Join(tempExtractDir, ent.Name)
		newFileInfo, err := os.Stat(name)
		if err != nil {
			t.Error(err)
			continue
		}
		checkFileInfo(t, &ent, newFileInfo)

		if ent.Type != "dir" {
			content, err := os.ReadFile(name)
			if err != nil {
				t.Error(err)
			}
			if string(content) != ent.Content {
				t.Errorf("File %s has mismatched contents", ent.Name)
			}
		}
	}
}

func TestDirectoryHardLink(t *testing.T) {
	tempDir := t.TempDir()

	// Open an archive containing two directories with the same inode (0).
	// We're trying to test if having the same inode will trigger a hard link.
	archiveFile, err := os.Open("testdata/dir-hard-link.cpio")
	if err != nil {
		t.Fatal(err)
	}
	c := testutil.Command(t, "-v", "i")
	c.Dir = tempDir
	c.Stdin = archiveFile

	out, err := c.Output()
	if err != nil {
		t.Fatalf("Extraction failed:\n%s\n%s\n%v\n", out, c.Stderr, err)
	}
}

func TestMain(m *testing.M) {
	testutil.Run(m, main)
}
