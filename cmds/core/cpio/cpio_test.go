// Copyright 2020 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"syscall"
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
	// Create a temporary directory
	tempDir, err := ioutil.TempDir("", "TestCpio")
	if err != nil {
		t.Fatalf("cannot create temporary directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	var targets = []dirEnt{
		{Name: "file1", Type: "file", Content: "Hello World"},
		{Name: "file2", Type: "file", Content: ""},
		{Name: "hardlinked", Type: "hardlink", Target: "file1"},
		{Name: "hardlinkedtofile2", Type: "hardlink", Target: "file2"},
		{Name: "directory1", Type: "dir"},
	}
	for _, ent := range targets {
		var name = filepath.Join(tempDir, ent.Name)
		switch ent.Type {
		case "dir":
			err := os.Mkdir(name, os.FileMode(0700))
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
		var name = filepath.Join(tempDir, ent.Name)
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
	archiveFile, err := ioutil.TempFile("", "archive.cpio")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(archiveFile.Name())
	defer archiveFile.Close()

	if _, err := archiveFile.Write(archive); err != nil {
		t.Fatal(err)
	}

	// Extract to a new directory
	tempExtractDir, err := ioutil.TempDir(tempDir, "extract")
	if err != nil {
		t.Fatalf("cannot create temporary directory: %v", err)
	}

	c = testutil.Command(t, "-v", "i")
	c.Dir = tempExtractDir
	c.Stdin = archiveFile

	out, err := c.Output()
	if err != nil {
		t.Fatalf("Extraction failed:\n%s\n%s\n%v\n", out, c.Stderr, err)
	}

	for _, ent := range targets {
		var name = filepath.Join(tempExtractDir, ent.Name)
		newFileInfo, err := os.Stat(name)
		if err != nil {
			t.Fatal(err)
		}
		newStatT := newFileInfo.Sys().(*syscall.Stat_t)
		statT := ent.FileInfo.Sys().(*syscall.Stat_t)
		if ent.FileInfo.Name() != newFileInfo.Name() ||
			ent.FileInfo.Size() != newFileInfo.Size() ||
			ent.FileInfo.Mode() != newFileInfo.Mode() ||
			ent.FileInfo.IsDir() != newFileInfo.IsDir() ||
			statT.Mode != newStatT.Mode ||
			statT.Uid != newStatT.Uid ||
			statT.Gid != newStatT.Gid ||
			statT.Nlink != newStatT.Nlink {
			msg := "File has mismatched attributes:\n"
			msg += "Property |   Original |  Extracted\n"
			msg += fmt.Sprintf("Name:    | %10s | %10s\n", ent.FileInfo.Name(), newFileInfo.Name())
			msg += fmt.Sprintf("Size:    | %10d | %10d\n", ent.FileInfo.Size(), newFileInfo.Size())
			msg += fmt.Sprintf("Mode:    | %10d | %10d\n", ent.FileInfo.Mode(), newFileInfo.Mode())
			msg += fmt.Sprintf("IsDir:   | %10t | %10t\n", ent.FileInfo.IsDir(), newFileInfo.IsDir())
			msg += fmt.Sprintf("FI Mode: | %10d | %10d\n", statT.Mode, newStatT.Mode)
			msg += fmt.Sprintf("Uid:     | %10d | %10d\n", statT.Uid, newStatT.Uid)
			msg += fmt.Sprintf("Gid:     | %10d | %10d\n", statT.Gid, newStatT.Gid)
			msg += fmt.Sprintf("NLink:   | %10d | %10d\n", statT.Nlink, newStatT.Nlink)
			t.Fatal(msg)
		}

		if ent.Type == "file" {
			content, err := ioutil.ReadFile(name)
			if err != nil {
				t.Fatal(err)
			}
			if string(content) != ent.Content {
				t.Fatalf("File %s has mismatched contents", ent.Name)
			}
		}
	}
}

func TestMain(m *testing.M) {
	testutil.Run(m, main)
}
