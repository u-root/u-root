// Copyright 2020 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

type dirEnt struct {
	FileInfo os.FileInfo
	Name     string
	Type     string
	Content  string
	Target   string
}

// prepareTestDir creates a testing directory and returns a list of
// dirEnts and inputFile for archive creation
func prepareTestDir(t *testing.T, tempDir string) ([]dirEnt, *os.File) {
	t.Helper()
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

	inputFile, err := os.CreateTemp(tempDir, "")
	if err != nil {
		t.Fatalf("%v", err)
	}

	for _, ent := range targets {
		name := filepath.Join(tempDir, ent.Name)
		if _, err := fmt.Fprintln(inputFile, name); err != nil {
			t.Fatalf("failed to write file path %v to input file: %v", ent.Name, err)
		}
	}
	inputFile.Seek(0, 0)

	return targets, inputFile
}

func TestCpioList(t *testing.T) {
	tmpDir := t.TempDir()
	targets, inputFile := prepareTestDir(t, tmpDir)

	archive, err := os.CreateTemp(tmpDir, "archive.cpio")
	if err != nil {
		t.Fatalf("failed to create temporary archive file: %v", err)
	}

	err = run([]string{"o"}, inputFile, archive, false, "newc")
	if err != nil {
		t.Fatalf("failed to build archive from filepaths: %v", err)
	}

	stdout := &bytes.Buffer{}
	err = run([]string{"t"}, archive, stdout, false, "newc")
	if err != nil {
		t.Fatalf("failed to list archive: %v", err)
	}

	stdoutStr := stdout.String()
	for _, ent := range targets {
		if !strings.Contains(stdoutStr, ent.Name) {
			t.Errorf("expected to find %q in output", ent.Name)
		}
	}
}

func TestCpio(t *testing.T) {
	debug = t.Logf
	// Create a temporary directory
	tempDir := t.TempDir()
	targets, inputFile := prepareTestDir(t, tempDir)

	archive := &bytes.Buffer{}
	err := run([]string{"o"}, inputFile, archive, true, "newc")
	if err != nil {
		t.Fatalf("failed to build archive from filepaths: %v", err)
	}

	// Cpio can't read from a non-seekable input (e.g. a pipe) in input mode.
	// Write the archive to a file instead.
	archiveFile, err := os.CreateTemp(tempDir, "archive.cpio")
	if err != nil {
		t.Fatal(err)
	}

	if _, err := archiveFile.Write(archive.Bytes()); err != nil {
		t.Fatal(err)
	}

	// Extract to a new directory
	tempExtractDir := t.TempDir()

	out := &bytes.Buffer{}
	// Change directory back afterwards to not interfer with the subsequent tests
	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Could not get current working directory: %v", err)
	}
	defer os.Chdir(wd)

	err = os.Chdir(tempExtractDir)
	if err != nil {
		t.Fatalf("Change to extraction directory %v failed: %#v", tempExtractDir, err)
	}

	err = run([]string{"i"}, archiveFile, out, true, "newc")
	if err != nil {
		t.Fatalf("Extraction failed:\n%#v\n%v\n", out, err)
	}

	for _, ent := range targets {
		name := filepath.Join(tempExtractDir, tempDir, ent.Name)

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
	// Open an archive containing two directories with the same inode (0).
	// We're trying to test if having the same inode will trigger a hard link.
	archiveFile, err := os.Open("testdata/dir-hard-link.cpio")
	if err != nil {
		t.Fatal(err)
	}

	// Change directory back afterwards to not interfer with the subsequent tests
	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Could not get current working directory: %v", err)
	}
	defer os.Chdir(wd)
	tempExtractDir := t.TempDir()
	err = os.Chdir(tempExtractDir)
	if err != nil {
		t.Fatalf("Change to dir %v failed: %v", tempExtractDir, err)
	}

	want := &bytes.Buffer{}
	err = run([]string{"i"}, archiveFile, want, true, "newc")
	if err != nil {
		t.Fatalf("Extraction failed:\n%v\n%v\n", want, err)
	}
}
