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
		args []string
		want string
	}{

		{
			args: []string{f.Name()},
			want: "usage: chmod mode filepath\n",
		},
		{
			args: []string{""},
			want: "usage: chmod mode filepath\n",
		},
		{
			args: []string{"01777", f.Name()},
			want: "Invalid octal value 1777. Value should be less than or equal to 0777.\n",
		},
		{
			args: []string{"0abas", f.Name()},
			want: "Unable to decode mode \"0abas\". Please use an octal value: strconv.ParseUint: parsing \"0abas\": invalid syntax\n",
		},
		{
			args: []string{"0777", "blah1234"},
			want: "chmod blah1234: no such file or directory\n",
		},
	} {
		cmd := exec.Command(testpath, v.args...)
		_, stderr, err := run(cmd)
		// Ignore the date and time because we're using Log.Fatalf
		if got := stderr[20:]; got != v.want {
			t.Errorf("Chmod for %q failed: got %q, want %q", v.args, got, v.want)
		}
		if err == nil {
			t.Errorf("Chmod for %q failed: got nil want err", v.args)
		}
	}
}
