// Copyright 2016-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/u-root/u-root/pkg/testutil"
)

type fileModeTrans struct {
	before os.FileMode
	after  os.FileMode
}

func run(c *exec.Cmd) (string, string, error) {
	var o, e bytes.Buffer
	c.Stdout, c.Stderr = &o, &e
	err := c.Run()
	return o.String(), e.String(), err
}

func TestChmodSimple(t *testing.T) {
	// Temporary directories.
	f, err := os.CreateTemp(t.TempDir(), "BLAH1")
	if err != nil {
		t.Fatalf("cannot create temporary file: %v", err)
	}
	defer f.Close()

	for k, v := range map[string]fileModeTrans{
		"0777":       {before: 0o000, after: 0o777},
		"0644":       {before: 0o777, after: 0o644},
		"u-rwx":      {before: 0o777, after: 0o077},
		"g-rx":       {before: 0o777, after: 0o727},
		"a-xr":       {before: 0o222, after: 0o222},
		"a-xw":       {before: 0o666, after: 0o444},
		"u-xw":       {before: 0o666, after: 0o466},
		"a=":         {before: 0o777, after: 0o000},
		"u=":         {before: 0o777, after: 0o077},
		"u-":         {before: 0o777, after: 0o777},
		"o+":         {before: 0o700, after: 0o700},
		"g=rx":       {before: 0o777, after: 0o757},
		"u=rx":       {before: 0o077, after: 0o577},
		"o=rx":       {before: 0o077, after: 0o075},
		"u=xw":       {before: 0o742, after: 0o342},
		"a-rwx":      {before: 0o777, after: 0o000},
		"a-rx":       {before: 0o777, after: 0o222},
		"a-x":        {before: 0o777, after: 0o666},
		"o+rwx":      {before: 0o000, after: 0o007},
		"a+rwx":      {before: 0o000, after: 0o777},
		"a+xrw":      {before: 0o000, after: 0o777},
		"a+xxxxxxxx": {before: 0o000, after: 0o111},
		"o+xxxxx":    {before: 0o000, after: 0o001},
		"a+rx":       {before: 0o000, after: 0o555},
		"a+r":        {before: 0o111, after: 0o555},
		"a=rwx":      {before: 0o000, after: 0o777},
		"a=rx":       {before: 0o000, after: 0o555},
	} {
		// Set up the 'before' state
		err := os.Chmod(f.Name(), v.before)
		if err != nil {
			t.Fatalf("chmod(%q) failed: %v", f.Name(), err)
		}

		// Set permissions using chmod.
		c := testutil.Command(t, k, f.Name())
		err = c.Run()
		if err != nil {
			t.Fatalf("setting permissions failed: %v", err)
		}

		// Check that it worked.
		checkPath(t, f.Name(), k, v)
	}
}

func checkPath(t *testing.T, path string, instruction string, v fileModeTrans) {
	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("stat(%q) failed: %v", path, err)
	}
	if got := info.Mode().Perm(); got != v.after {
		t.Errorf("Wrong file permissions on %q: executed %s, before %0o, got %0o, want %0o", path, instruction, v.before, got, v.after)
	}
}

func TestChmodRecursive(t *testing.T) {
	// Temporary directories.
	tempDir := t.TempDir()

	var targetDirectories []string
	for _, dir := range []string{
		"L1_A", "L1_B", "L1_C",
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
		err := os.Mkdir(dir, os.FileMode(0o700))
		if err != nil {
			t.Fatalf("cannot create test directory: %v", err)
		}
		targetDirectories = append(targetDirectories, dir)
	}

	for k, v := range map[string]fileModeTrans{
		"0707":      {before: 0o755, after: 0o707},
		"0770":      {before: 0o755, after: 0o770},
		"o-rwx":     {before: 0o777, after: 0o770},
		"g-rx":      {before: 0o777, after: 0o727},
		"a=rrrrrwx": {before: 0o777, after: 0o777},
		"a+w":       {before: 0o700, after: 0o722},
		"g+xr":      {before: 0o700, after: 0o750},
		// This last step is causing a flake in the testing package, somehow.
		// testing.go:967: TempDir RemoveAll cleanup: unlinkat /tmp/TestChmodRecursive2303678599/001/L1_A/L2_A: permission denied
		// It's not always, but i is frequent enough to break things.
		// Just disable it; recursion is tested here, and a=rx elsewhere.
		// "a=rx":      {before: 0o777, after: 0o555},
	} {

		// Set up the 'before' state
		for _, dir := range targetDirectories {
			err := os.Chmod(dir, v.before)
			if err != nil {
				t.Fatalf("chmod(%q) failed: %v", dir, err)
			}
		}

		// Set permissions using chmod.
		c := testutil.Command(t, "-R", k, tempDir)
		if err := c.Run(); err != nil {
			t.Fatalf("setting permissions failed: %v", err)
		}

		// Check that it worked.
		for _, dir := range targetDirectories {
			checkPath(t, dir, k, v)
		}
	}
}

func TestChmodReference(t *testing.T) {
	// Temporary directories.
	tempDir := t.TempDir()

	sourceFile, err := os.CreateTemp(tempDir, "BLAH1")
	if err != nil {
		t.Fatalf("cannot create temporary file: %v", err)
	}
	defer sourceFile.Close()

	targetFile, err := os.CreateTemp(tempDir, "BLAH2")
	if err != nil {
		t.Fatalf("cannot create temporary file: %v", err)
	}
	defer targetFile.Close()

	for _, perm := range []os.FileMode{0o777, 0o644} {
		err = os.Chmod(sourceFile.Name(), perm)
		if err != nil {
			t.Fatalf("chmod(%q) failed: %v", sourceFile.Name(), err)
		}

		// Set target file permissions using chmod.
		c := testutil.Command(t,
			"--reference",
			sourceFile.Name(),
			targetFile.Name())
		err = c.Run()
		if err != nil {
			t.Fatalf("setting permissions failed: %v", err)
		}

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
	f, err := os.CreateTemp(t.TempDir(), "BLAH1")
	if err != nil {
		t.Fatalf("cannot create temporary file: %v", err)
	}
	defer f.Close()

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
			want:     "Unable to decode mode \"0abas\". Please use an octal value or a valid mode string.\n",
			skipTo:   20,
			skipFrom: -1,
		},
		{
			args:     []string{"0777", "blah1234"},
			want:     "chmod blah1234: no such file or directory\n",
			skipTo:   20,
			skipFrom: -1,
		},
		{
			args:     []string{"a=9rwx", f.Name()},
			want:     "Unable to decode mode \"a=9rwx\". Please use an octal value or a valid mode string.\n",
			skipTo:   20,
			skipFrom: -1,
		},
		{
			args:     []string{"+r", f.Name()},
			want:     "Unable to decode mode \"+r\". Please use an octal value or a valid mode string.\n",
			skipTo:   20,
			skipFrom: -1,
		},
		{
			args:     []string{"a%rwx", f.Name()},
			want:     "Unable to decode mode \"a%rwx\". Please use an octal value or a valid mode string.\n",
			skipTo:   20,
			skipFrom: -1,
		},
		{
			args:     []string{"b=rwx", f.Name()},
			want:     "Unable to decode mode \"b=rwx\". Please use an octal value or a valid mode string.\n",
			skipTo:   20,
			skipFrom: -1,
		},
	} {
		cmd := testutil.Command(t, v.args...)
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

func TestMain(m *testing.M) {
	testutil.Run(m, main)
}
