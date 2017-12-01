// Copyright 2012 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
//  created by Manoel Vilela (manoel_vilela@engineer.com)

package main

import (
	"io/ioutil"
	"os"
	"path"
	"syscall"
	"testing"

	"github.com/u-root/u-root/pkg/testutil"
)

type file struct {
	name   string
	delete bool
}

type rmTestCase struct {
	name  string
	files []file
	i     bool
	r     bool
	f     bool
	err   func(error) bool
	stdin *testutil.FakeStdin
}

func TestRemove(t *testing.T) {
	var (
		no       = testutil.NewFakeStdin("no")
		fbody    = []byte("Go is cool!")
		tmpFiles = []struct {
			name  string
			mode  os.FileMode
			isdir bool
		}{

			{
				name:  "hi",
				mode:  0755,
				isdir: true,
			},
			{
				name: "hi/one.txt",
				mode: 0666,
			},
			{
				name: "hi/two.txt",
				mode: 0777,
			},
			{
				name: "go.txt",
				mode: 0555,
			},
		}
		nilerr    = func(err error) bool { return err == nil }
		testCases = []rmTestCase{
			{
				name: "no flags",
				files: []file{
					{"hi/one.txt", true},
					{"hi/two.txt", true},
					{"go.txt", true},
				},
				err:   nilerr,
				stdin: no,
			},
			{
				name: "-i",
				files: []file{
					{"hi/one.txt", true},
					{"hi/two.txt", false},
					{"go.txt", true},
				},
				i:     true,
				err:   nilerr,
				stdin: testutil.NewFakeStdin("y", "no", "yes"),
			},
			{
				name: "nonexistent with no flags",
				files: []file{
					{"hi/one.txt", true},
					{"hi/two.doc", true}, // does not exist
					{"go.txt", false},
				},
				err:   os.IsNotExist,
				stdin: no,
			},
			{
				name:  "directory with no flags",
				files: []file{{"hi", false}},
				err:   pathError(syscall.ENOTEMPTY),
				stdin: no,
			},
			{
				name:  "directory with -f",
				files: []file{{"hi", false}},
				f:     true,
				err:   pathError(syscall.ENOTEMPTY),
				stdin: no,
			},
			{
				name:  "directory with -r",
				files: []file{{"hi", true}},
				r:     true,
				err:   nilerr,
				stdin: no,
			},
			{
				name: "directory and file with -r",
				files: []file{
					{"hi", true},
					{"go.txt", true},
				},
				r:     true,
				err:   nilerr,
				stdin: no,
			},
			{
				name: "-f",
				files: []file{
					{"hi/one.doc", true}, // does not exist
					{"hi/two.txt", true},
					{"go.doc", true}, // does  not exist
				},
				f:     true,
				err:   nilerr,
				stdin: no,
			},
			{
				name: "-i -f",
				files: []file{
					{"hi/one.txt", true}, // does not exist
					{"hi/two.txt", true},
					{"go.txt", true}, // does  not exist
				},
				f:     true,
				i:     true,
				err:   nilerr,
				stdin: no,
			},
		}
	)

	for _, tc := range testCases {

		t.Run(tc.name, func(t *testing.T) {
			d, err := ioutil.TempDir(os.TempDir(), "u-root.cmds.rm")
			if err != nil {
				t.Fatal(err)
			}
			defer os.RemoveAll(d)

			for _, f := range tmpFiles {
				var (
					err      error
					filepath = path.Join(d, f.name)
				)
				if f.isdir {
					err = os.Mkdir(filepath, f.mode)
				} else {
					err = ioutil.WriteFile(filepath, fbody, f.mode)
				}
				if err != nil {
					t.Fatal(err)
				}
			}
			testRemove(t, d, tc)
		})
	}
}

func testRemove(t *testing.T, dir string, tc rmTestCase) {
	var files = make([]string, len(tc.files))
	for i, f := range tc.files {
		files[i] = path.Join(dir, f.name)
	}

	flags.v = true
	flags.r = tc.r
	flags.f = tc.f
	flags.i = tc.i

	if err := rm(tc.stdin, files); !tc.err(err) {
		t.Error(err)
	}

	if flags.i && tc.stdin.Count() == 0 {
		t.Error("Expected reading from stdin")
	} else if !flags.i && tc.stdin.Count() > 0 {
		t.Errorf("Did not expect reading %d times from stdin", tc.stdin.Count())
	}
	if tc.stdin.Overflowed() {
		t.Error("Read from stdin too many times")
	}

	for i, f := range tc.files {
		_, err := os.Stat(path.Join(dir, f.name))
		if tc.files[i].delete != os.IsNotExist(err) {
			t.Errorf("File %q deleted: %t, expected: %t",
				f.name, os.IsNotExist(err), tc.files[i].delete)
		}
	}
}

func pathError(errno syscall.Errno) func(error) bool {
	return func(err error) bool {
		pe, ok := err.(*os.PathError)
		if !ok {
			return false
		}
		return pe.Err == errno
	}
}
