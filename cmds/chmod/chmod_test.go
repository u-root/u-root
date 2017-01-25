// Copyright 2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

var (
	testPath = "."
	// if true removeAll the testPath on the end
	remove = true
)

type test struct {
	args    []string
	expects string
}

func run(c *exec.Cmd) (string, string, error) {
	var o, e bytes.Buffer
	c.Stdout, c.Stderr = &o, &e
	err := c.Run()
	return o.String(), e.String(), err
}

func TestChmodSimple(t *testing.T) {
	tempDir, err := ioutil.TempDir("", "TestChmodSimple")
	if err != nil {
		t.Fatalf("cannot create temporary directory: %v", err)
	}
	if remove {
		defer os.RemoveAll(tempDir)
	}
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

	// Set permissions and check that they are correct
	c := exec.Command(testpath, "0777", f.Name())
	c.Run()
	fileinfo, _ := os.Stat(f.Name())
	fileperm := fileinfo.Mode().Perm()
	if fileperm != 0777 {
		t.Errorf("Wrong file permissions")
	}

	// Change permissions and check that they are correct
	c = exec.Command(testpath, "0644", f.Name())
	c.Run()
	fileinfo, _ = os.Stat(f.Name())
	fileperm = fileinfo.Mode().Perm()
	if fileperm != 0644 {
		t.Errorf("Wrong file permissions")
	}
}

func TestInvocationErrors(t *testing.T) {
	tempDir, err := ioutil.TempDir("", "TestInvocationErrors")
	if err != nil {
		t.Fatalf("cannot create temporary directory: %v", err)
	}
	if remove {
		defer os.RemoveAll(tempDir)
	}
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

	var tests = []test{
		{args: []string{f.Name()}, expects: "usage: chmod mode filepath\n"},
		{args: []string{""}, expects: "usage: chmod mode filepath\n"},
		{args: []string{"01777", f.Name()}, expects: "Invalid octal value. Value larger than 777, was 1777\n"},
		{args: []string{"0abas", f.Name()}, expects: "Unable to decode mode. Please use an octal value. arg was 0abas, err was strconv.ParseUint: parsing \"0abas\": invalid syntax\n"},
		{args: []string{"0777", "blah1234"}, expects: "Unable to chmod, filename was blah1234, err was no such file or directory\n"},
	}

	for _, v := range tests {
		c := exec.Command(testpath, v.args...)
		_, e, err := run(c)
		// Ignore the date and time because we're using Log.Fatalf
		if e[20:] != v.expects {
			t.Errorf("Chmod for '%v' failed: got '%s', want '%s'", v.args, e[20:], v.expects)
		}
		if err == nil {
			t.Errorf("Kill for '%v' failed: got nil want err", v.args)
		}
	}
}
