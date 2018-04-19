// Copyright 2016-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func runChmodBinary(f flags, args []string) error {
	var a []string
	if f.recursive {
		a = append(a, "-R")
	}
	if len(f.reference) > 0 {
		a = append(a, "--reference", f.reference)
	}
	a = append(a, args...)
	c := exec.Command(os.Getenv("EXECPATH"), a...)
	c.Stdout, c.Stderr = os.Stdout, os.Stderr
	return c.Run()
}

func TestChmod(t *testing.T) {
	chmod := chmodMain
	// If EXECPATH is set, test this implementation against the binary
	// given by EXECPATH.
	if len(os.Getenv("EXECPATH")) != 0 {
		chmod = runChmodBinary
	}

	tempDir, err := ioutil.TempDir("", "chmod")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	// Set up single file for simple test.
	f, err := ioutil.TempFile(tempDir, "chmod-tmp-test")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	// Set up single file for simple test.
	f2, err := ioutil.TempFile(tempDir, "chmod-tmp-test")
	if err != nil {
		t.Fatal(err)
	}
	defer f2.Close()

	file0544, err := ioutil.TempFile(tempDir, "file0544")
	if err != nil {
		t.Fatal(err)
	}
	defer file0544.Close()
	if err := os.Chmod(file0544.Name(), 0544); err != nil {
		t.Fatal(err)
	}

	// Set up complicated directory structure for recursive test.
	recDir, err := ioutil.TempDir(tempDir, "recursive")
	if err != nil {
		t.Fatal(err)
	}

	recursives := []string{recDir}
	for _, dir := range []string{
		"L1_A",
		"L1_B",
		"L1_C",
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
		dir = filepath.Join(recDir, dir)
		if err := os.MkdirAll(dir, os.FileMode(0700)); err != nil {
			t.Fatalf("cannot create test directory: %v", err)
		}
		recursives = append(recursives, dir)
	}

	for _, tt := range []struct {
		flags flags
		args  []string
		want  error

		// fileList is the list of files that wantMode should be set on
		// after calling chmod.
		fileList []string
		wantMode os.FileMode
	}{
		{
			args:     []string{"0777", f.Name()},
			want:     nil,
			fileList: []string{f.Name()},
			wantMode: 0777,
		},
		{
			args:     []string{"0644", f.Name()},
			want:     nil,
			fileList: []string{f.Name()},
			wantMode: 0644,
		},
		{
			args: []string{"0707", recDir},
			want: nil,
			flags: flags{
				recursive: true,
			},
			fileList: recursives,
			wantMode: 0707,
		},
		{
			args: []string{"0770", recDir},
			want: nil,
			flags: flags{
				recursive: true,
			},
			fileList: recursives,
			wantMode: 0770,
		},
		{
			args: []string{f2.Name()},
			want: nil,
			flags: flags{
				reference: file0544.Name(),
			},
			fileList: []string{f2.Name()},
			wantMode: 0544,
		},
		{
			args: []string{"01777", f.Name()},
			want: fmt.Errorf("invalid octal value 1777: value should be less than or equal to 0777"),
		},
		{
			args: []string{"0abas", f.Name()},
			want: fmt.Errorf("unable to decode mode \"0abas\": must use an octal value: strconv.ParseUint: parsing \"0abas\": invalid syntax"),
		},
		{
			args: []string{"0777", "blah1234"},
			want: fmt.Errorf("chmod blah1234: no such file or directory"),
		},
	} {
		if err := chmod(tt.flags, tt.args); err != tt.want && (err == nil || tt.want == nil || err.Error() != tt.want.Error()) {
			t.Errorf("chmod(%#v, %v) = %v, want %v", tt.flags, tt.args, err, tt.want)
		}

		for _, file := range tt.fileList {
			checkPath(t, file, tt.wantMode)
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
