// Copyright 2015-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package find

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/u-root/u-root/pkg/uroot/util"
)

type tests struct {
	name  string
	opts  func(*Finder) error
	names string
	errs  string
}

func noop(f *Finder) error {
	return nil
}

// TODO: I don't now where this subtesting stuff originated, I just copied it,
// but it's bad practice as you can not pick individual tests.
// Break this out into individual tests.
func TestSimple(t *testing.T) {
	var testCases = []tests{
		{name: "basic find",
			opts: func(_ *Finder) error { return nil },
			names: `

/root
/root/ab
/root/ab/c
/root/ab/c/d
/root/ab/c/d/e
/root/ab/c/d/e/f
/root/ab/c/d/e/f/ghij
/root/ab/c/d/e/f/ghij/k
/root/ab/c/d/e/f/ghij/k/l
/root/ab/c/d/e/f/ghij/k/l/m
/root/ab/c/d/e/f/ghij/k/l/m/n
/root/ab/c/d/e/f/ghij/k/l/m/n/o
/root/ab/c/d/e/f/ghij/k/l/m/n/o/p
/root/ab/c/d/e/f/ghij/k/l/m/n/o/p/q
/root/ab/c/d/e/f/ghij/k/l/m/n/o/p/q/r
/root/ab/c/d/e/f/ghij/k/l/m/n/o/p/q/r/s
/root/ab/c/d/e/f/ghij/k/l/m/n/o/p/q/r/s/t
/root/ab/c/d/e/f/ghij/k/l/m/n/o/p/q/r/s/t/u
/root/ab/c/d/e/f/ghij/k/l/m/n/o/p/q/r/s/t/u/v
/root/ab/c/d/e/f/ghij/k/l/m/n/o/p/q/r/s/t/u/v/w
/root/ab/c/d/e/f/ghij/k/l/m/n/o/p/q/r/s/t/u/v/w/xyz
/root/ab/c/d/e/f/ghij/k/l/m/n/o/p/q/r/s/t/u/v/w/xyz/0777
/root/ab/c/d/e/f/ghij/k/l/m/n/o/p/q/r/s/t/u/v/w/xyz/file`, errs: ""},
		{name: "just a dir",
			opts: func(f *Finder) error {
				f.Mode = os.ModeDir
				f.ModeMask = os.ModeDir
				return nil
			}, names: `

/root
/root/ab
/root/ab/c
/root/ab/c/d
/root/ab/c/d/e
/root/ab/c/d/e/f
/root/ab/c/d/e/f/ghij
/root/ab/c/d/e/f/ghij/k
/root/ab/c/d/e/f/ghij/k/l
/root/ab/c/d/e/f/ghij/k/l/m
/root/ab/c/d/e/f/ghij/k/l/m/n
/root/ab/c/d/e/f/ghij/k/l/m/n/o
/root/ab/c/d/e/f/ghij/k/l/m/n/o/p
/root/ab/c/d/e/f/ghij/k/l/m/n/o/p/q
/root/ab/c/d/e/f/ghij/k/l/m/n/o/p/q/r
/root/ab/c/d/e/f/ghij/k/l/m/n/o/p/q/r/s
/root/ab/c/d/e/f/ghij/k/l/m/n/o/p/q/r/s/t
/root/ab/c/d/e/f/ghij/k/l/m/n/o/p/q/r/s/t/u
/root/ab/c/d/e/f/ghij/k/l/m/n/o/p/q/r/s/t/u/v
/root/ab/c/d/e/f/ghij/k/l/m/n/o/p/q/r/s/t/u/v/w
/root/ab/c/d/e/f/ghij/k/l/m/n/o/p/q/r/s/t/u/v/w/xyz`, errs: ""},
		{name: "just a file",
			opts: func(f *Finder) error {
				f.Mode = 0
				f.ModeMask = os.ModeType
				return nil
			}, names: `
/root/ab/c/d/e/f/ghij/k/l/m/n/o/p/q/r/s/t/u/v/w/xyz/0777
/root/ab/c/d/e/f/ghij/k/l/m/n/o/p/q/r/s/t/u/v/w/xyz/file`, errs: ""},
		{name: "file by mode",
			opts: func(f *Finder) error {
				f.Mode = 0444
				f.ModeMask = os.ModePerm
				f.Debug = t.Logf
				return nil
			}, names: `
/root/ab/c/d/e/f/ghij/k/l/m/n/o/p/q/r/s/t/u/v/w/xyz/0777`, errs: ""},
		{name: "file by name",
			opts: func(f *Finder) error {
				f.Debug = t.Logf
				f.Pattern = "*file"
				return nil
			}, names: `
/root/ab/c/d/e/f/ghij/k/l/m/n/o/p/q/r/s/t/u/v/w/xyz/file`, errs: ""},
	}
	d, err := ioutil.TempDir(os.TempDir(), "u-root.cmds.find")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(d)

	var namespace = []util.Creator{
		util.Dir{Name: filepath.Join(d, "root/ab/c/d/e/f/ghij/k/l/m/n/o/p/q/r/s/t/u/v/w/xyz"), Mode: 0775},
		util.File{Name: filepath.Join(d, "root//ab/c/d/e/f/ghij/k/l/m/n/o/p/q/r/s/t/u/v/w/xyz/file"), Mode: 0664},
		util.File{Name: filepath.Join(d, "root//ab/c/d/e/f/ghij/k/l/m/n/o/p/q/r/s/t/u/v/w/xyz/0777"), Mode: 0444},
	}
	for _, c := range namespace {
		if err := c.Create(); err != nil {
			t.Fatalf("Error creating %s: %v", c, err)
		}
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			f, err := New(func(f *Finder) error {
				f.Root = d
				return nil
			}, tc.opts)
			if err != nil {
				t.Fatal(err)
			}
			go f.Find()

			var names, errs string
			for o := range f.Names {
				if o.Err != nil {
					t.Errorf("%v: got %v, want nil", o.Name, o.Err)
				}
				// get rid of the prefix, as it changes for every test.
				names = fmt.Sprintf("%s\n%s", names, o.Name[len(d):])
			}

			t.Logf("names %v", names)
			if len(names) != len(tc.names) {
				t.Errorf("Find output: got %d bytes, want %d bytes", len(names), len(tc.names))
			}
			if names != tc.names {
				t.Errorf("Find output: got %v, want %v", names, tc.names)
			}
			t.Logf("errs %v", errs)
			if len(errs) != len(tc.errs) {
				t.Errorf("Find output: got %d bytes, want %d bytes", len(errs), len(tc.errs))
			}
			if errs != tc.errs {
				t.Errorf("Find errors: got %v, want %v", errs, tc.errs)
			}
		})
	}
}
