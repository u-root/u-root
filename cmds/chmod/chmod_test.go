// Copyright 2016-2017 the u-root Authors. All rights reserved
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
	"testing"
)

var (
	testPath = "."
)

func run(c *exec.Cmd) (string, string, error) {
	var o, e bytes.Buffer
	c.Stdout, c.Stderr = &o, &e
	err := c.Run()
	return o.String(), e.String(), err
}

func TestChmodSimple(t *testing.T) {
	// Temporary directories.
	tempDir, err := ioutil.TempDir("", "TestChmodSimple")
	if err != nil {
		t.Fatalf("cannot create temporary directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	f, err := ioutil.TempFile(tempDir, "BLAH1")
	if err != nil {
		t.Fatalf("cannot create temporary file: %v", err)
	}
	defer f.Close()

	// Build chmod binary.
	testpath := filepath.Join(tempDir, "testchmod.exe")
	out, err := exec.Command("go", "build", "-o", testpath, ".").CombinedOutput()
	if err != nil {
		t.Fatalf("go build -o %v cmds/chmod: %v\n%s", testpath, err, string(out))
	}

	for _, perm := range []os.FileMode{0777, 0644} {
		// Set permissions using chmod.
		c := exec.Command(testpath, fmt.Sprintf("%0o", perm), f.Name())
		c.Run()

		// Check that it worked.
		info, err := os.Stat(f.Name())
		if err != nil {
			t.Fatalf("stat(%q) failed: %v", f.Name(), err)
		}
		if got := info.Mode().Perm(); got != perm {
			t.Errorf("Wrong file permissions on %q: got %0o, want %0o", f.Name(), got, perm)
		}
	}
}

func checkPath(t *testing.T, path string, want os.FileMode) {
	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("stat(%q) failed: %v", path, err)
	}
	if got := info.Mode().Perm(); got != want {
		t.Fatalf("Wrong file permissions on file %q: got %0o, want %0o",
			path, got, want)
	}
}
func TestChmodRecursive(t *testing.T) {
	// Temporary directories.
	tempDir, err := ioutil.TempDir("", "TestChmodRecursive")
	if err != nil {
		t.Fatalf("cannot create temporary directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	var targetFiles []string
	var targetDirectories []string
	for _, dir := range []string{"L1_A", "L1_B", "L1_C",
		filepath.Join("L1_A", "L2_A"),
		filepath.Join("L1_A", "L2_B"),
		filepath.Join("L1_A", "L2_C"),
		filepath.Join("L1_B", "L2_A"),
		filepath.Join("L1_B", "L2_B"),
		filepath.Join("L1_B", "L2_C"),
		filepath.Join("L1_C", "L2_A"),
		filepath.Join("L1_C", "L2_B"),
		filepath.Join("L1_C", "L2_C"),
	} {
		dir = filepath.Join(tempDir, dir)
		err := os.Mkdir(dir, os.FileMode(0700))
		if err != nil {
			t.Fatalf("cannot create test directory: %v", err)
		}
		targetDirectories = append(targetDirectories, dir)
		targetFile, err := os.Create(filepath.Join(dir, "X"))
		if err != nil {
			t.Fatalf("cannot create temporary file: %v", err)
		}
		targetFiles = append(targetFiles, targetFile.Name())

	}

	// Build chmod binary.
	testpath := filepath.Join(tempDir, "testchmod.exe")
	out, err := exec.Command("go", "build", "-o", testpath, ".").CombinedOutput()
	if err != nil {
		t.Fatalf("go build -o %v cmds/chmod: %v\n%s", testpath, err, string(out))
	}

	for _, perm := range []os.FileMode{0707, 0770} {
		// Set target file permissions using chmod.
		c := exec.Command(testpath,
			"-R",
			fmt.Sprintf("%0o", perm),
			tempDir)
		c.Run()

		// Check that it worked.
		for _, dir := range targetDirectories {
			checkPath(t, dir, perm)
		}
	}
}

func TestChmodReference(t *testing.T) {
	// Temporary directories.
	tempDir, err := ioutil.TempDir("", "TestChmodReference")
	if err != nil {
		t.Fatalf("cannot create temporary directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	sourceFile, err := ioutil.TempFile(tempDir, "BLAH1")
	if err != nil {
		t.Fatalf("cannot create temporary file: %v", err)
	}
	defer sourceFile.Close()

	targetFile, err := ioutil.TempFile(tempDir, "BLAH2")
	if err != nil {
		t.Fatalf("cannot create temporary file: %v", err)
	}
	defer targetFile.Close()

	// Build chmod binary.
	testpath := filepath.Join(tempDir, "testchmod.exe")
	out, err := exec.Command("go", "build", "-o", testpath, ".").CombinedOutput()
	if err != nil {
		t.Fatalf("go build -o %v cmds/chmod: %v\n%s", testpath, err, string(out))
	}

	for _, perm := range []os.FileMode{0777, 0644} {
		os.Chmod(sourceFile.Name(), perm)

		// Set target file permissions using chmod.
		c := exec.Command(testpath,
			"--reference",
			sourceFile.Name(),
			targetFile.Name())
		c.Run()

		// Check that it worked.
		info, err := os.Stat(targetFile.Name())
		if err != nil {
			t.Fatalf("stat(%q) failed: %v", targetFile.Name(), err)
		}
		if got := info.Mode().Perm(); got != perm {
			t.Fatalf("Wrong file permissions on file %q: got %0o, want %0o",
				targetFile.Name(), got, perm)
		}
	}
}

func TestInvocationErrors(t *testing.T) {
	tempDir, err := ioutil.TempDir("", "TestInvocationErrors")
	if err != nil {
		t.Fatalf("cannot create temporary directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	f, err := ioutil.TempFile(tempDir, "BLAH1")
	if err != nil {
		t.Fatalf("cannot create temporary file: %v", err)
	}
	defer f.Close()

	testpath := filepath.Join(tempDir, "testchmod.exe")
	out, err := exec.Command("go", "build", "-o", testpath, ".").CombinedOutput()
	if err != nil {
		t.Fatalf("go build -o %v cmds/chmod: %v\n%s", testpath, err, string(out))
	}

	for _, v := range []struct {
		args     []string
		want     string
		skipTo   int
		skipFrom int
	}{

		{
			args:     []string{f.Name()},
			want:     "Usage",
			skipTo:   0,
			skipFrom: len("Usage"),
		},
		{
			args:     []string{""},
			want:     "Usage",
			skipTo:   0,
			skipFrom: len("Usage"),
		},
		{
			args:     []string{"01777", f.Name()},
			want:     "Invalid octal value 1777. Value should be less than or equal to 0777.\n",
			skipTo:   20,
			skipFrom: -1,
		},
		{
			args:     []string{"0abas", f.Name()},
			want:     "Unable to decode mode \"0abas\". Please use an octal value: strconv.ParseUint: parsing \"0abas\": invalid syntax\n",
			skipTo:   20,
			skipFrom: -1,
		},
		{
			args:     []string{"0777", "blah1234"},
			want:     "chmod blah1234: no such file or directory\n",
			skipTo:   20,
			skipFrom: -1,
		},
	} {
		cmd := exec.Command(testpath, v.args...)
		_, stderr, err := run(cmd)
		if v.skipFrom == -1 {
			v.skipFrom = len(stderr)
		}
		// Ignore the date and time because we're using Log.Fatalf
		if got := stderr[v.skipTo:v.skipFrom]; got != v.want {
			t.Errorf("Chmod for %q failed: got %q, want %q", v.args, got, v.want)
		}
		if err == nil {
			t.Errorf("Chmod for %q failed: got nil want err", v.args)
		}
	}
}
