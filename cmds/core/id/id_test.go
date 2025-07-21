// Copyright 2017-2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build !plan9

package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"testing"
)

func TestBadFiles(t *testing.T) {
	flags := &flags{}

	d := t.TempDir()
	n := filepath.Join(d, "nosuchfile")
	f := filepath.Join(d, "afile")
	if err := os.WriteFile(f, []byte{}, 0o666); err != nil {
		t.Fatalf("writing %q: want nil, got %v", f, err)
	}
	if err := run(io.Discard, "root", flags, n, f); !errors.Is(err, os.ErrNotExist) {
		t.Errorf("Using %q for passwd: want %v, got nil", n, os.ErrNotExist)
	}
	if err := run(io.Discard, "root", flags, f, n); !errors.Is(err, os.ErrNotExist) {
		t.Errorf("Using %q for group: want %v, got nil", n, os.ErrNotExist)
	}
}

func TestGoodFiles(t *testing.T) {
	tests := []struct {
		passwdContent  string
		groupContent   string
		userName       string
		expectedOutput string
		flags          flags
	}{
		{
			passwdContent:  "root:x:0:0:root:/root:/bin/bash",
			groupContent:   "root:x:0:",
			userName:       "root",
			flags:          flags{},
			expectedOutput: "uid=0(root) gid=0(root) groups=0(root)\n",
		},
		{
			passwdContent:  "root:x:0:0:root:/root:/bin/bash",
			groupContent:   "root:x:0:",
			userName:       "root",
			flags:          flags{name: true, user: true},
			expectedOutput: "root\n",
		},
		{
			passwdContent:  "root:x:0:0:root:/root:/bin/bash",
			groupContent:   "root:x:0:",
			userName:       "root",
			flags:          flags{name: true, group: true},
			expectedOutput: "root\n",
		},
		{
			passwdContent:  "root:x:0:0:root:/root:/bin/bash",
			groupContent:   "root:x:0:",
			userName:       "root",
			flags:          flags{group: true},
			expectedOutput: "0\n",
		},
		{
			passwdContent:  "root:x:0:0:root:/root:/bin/bash",
			groupContent:   "root:x:0:\nother:x:1:root",
			userName:       "root",
			flags:          flags{groups: true},
			expectedOutput: "0 1 \n",
		},
		{
			passwdContent:  "root:x:0:0:root:/root:/bin/bash",
			groupContent:   "root:x:0:",
			userName:       "root",
			flags:          flags{user: true},
			expectedOutput: "0\n",
		},
	}

	for _, test := range tests {
		d := t.TempDir()
		passwd, err := os.CreateTemp(d, "passwd")
		if err != nil {
			t.Fatal(err)
		}
		_, err = passwd.WriteString(test.passwdContent)
		if err != nil {
			t.Fatal(err)
		}
		group, err := os.CreateTemp(d, "group")
		if err != nil {
			t.Fatal(err)
		}
		_, err = group.WriteString(test.groupContent)
		if err != nil {
			t.Fatal(err)
		}

		var stdout bytes.Buffer
		if err := run(&stdout, test.userName, &test.flags, passwd.Name(), group.Name()); err != nil {
			t.Fatal(err)
		}

		if got, want := stdout.String(), test.expectedOutput; got != want {
			t.Errorf("want %q, got %q", want, got)
		}
	}
}

type passwd struct {
	name string
	uid  int
	gid  int
}

var passwdShould = []passwd{
	{"root", 0, 0},
	{"bin", 1, 1},
	{"daemon", 2, 2},
	{"lary", 1000, 1000},
	{"curly", 1001, 1001},
	{"moe", 1002, 2002},
}

var passwdShouldnt = []passwd{
	{"adm", 3, 4},
}

var passwdFiles = []string{
	"testdata/passwd-simple.txt",
	"testdata/passwd-comments.txt",
}

type group struct {
	name string
	gid  int
}

var groupShould = []group{
	{"printadmin", 997},
	{"ssh_keys", 996},
	{"rpcuser", 29},
	{"nfsnobody", 65534},
	{"sshd", 74},
	{"wheel", 10},
}

var groupShouldnt = []group{
	{"bad", 314},
	{"wrong", 996},
	{"wheel", 11},
}

var groupFiles = []string{
	"testdata/group-simple.txt",
	"testdata/group-comments.txt",
}

var groupMembers = map[string][]int{
	"larry": {10, 74},
	"curly": {10, 29},
	"moe":   {10},
	"joe":   {},
}

func passwdSame(u *Users, us passwd) error {
	var s string
	var d int
	var err error
	d, err = u.GetUID(us.name)
	if err != nil {
		return fmt.Errorf("failed to GetUID expected user %s: %w", us.name, err)
	}
	if d != us.uid {
		return fmt.Errorf("wrong UID for %s: got %d, expected %d", us.name, d, us.uid)
	}

	d, err = u.GetGID(us.uid)
	if err != nil {
		return fmt.Errorf("failed to GetGID expected uid %d: %w", us.uid, err)
	}
	if d != us.gid {
		return fmt.Errorf("wrong GID for uid %d: got %d, expected %d", us.uid, d, us.gid)
	}

	s, err = u.GetUser(us.uid)
	if err != nil {
		return fmt.Errorf("failed to GetUser expected user %s: %w", us.name, err)
	}
	if s != us.name {
		return fmt.Errorf("wrong username for %d: got %s, expected %s", us.uid, s, us.name)
	}
	return nil
}

func TestUsers(t *testing.T) {
	t.Run("non-existent passwd file", func(t *testing.T) {
		f := "testdata/does-not-exist"
		u, e := NewUsers(f)
		if e == nil {
			t.Errorf("NewUser on non-existant file should return an error")
		}
		if u == nil {
			t.Errorf("NewUser should return a valid Users object, even on error")
		}
	})
	t.Run("empty passwd file", func(t *testing.T) {
		f := "testdata/passwd-empty.txt"
		u, e := NewUsers(f)
		if e != nil {
			t.Errorf("NewUser should not report error for empty passwd file")
		}
		if u == nil {
			t.Errorf("NewUser should return a valid Users object even if passwd file is empty")
		}
	})
	t.Run("almost empty passwd file", func(t *testing.T) {
		f := "testdata/passwd-newline.txt"
		u, e := NewUsers(f)
		if e != nil {
			t.Errorf("NewUser should not report error for empty passwd file")
		}
		if u == nil {
			t.Errorf("NewUser should return a valid Users object even if passwd file is empty")
		}
	})
	for _, f := range passwdFiles {
		t.Run(f, func(t *testing.T) {
			u, e := NewUsers(f)
			if e != nil {
				t.Errorf("NewUser should not return an error on valid file")
			}
			if u == nil {
				t.Errorf("NewUser should return a valid Users object on valid file")
			}
			for _, us := range passwdShould {
				if e := passwdSame(u, us); e != nil {
					t.Errorf("%v", e)
				}
			}
			for _, us := range passwdShouldnt {
				if e := passwdSame(u, us); e == nil {
					t.Errorf("user %s matched when it shouldn't", us.name)
				}
			}
		})
	}
}

func groupSame(g *Groups, gs group) error {
	var s string
	var d int
	var err error

	d, err = g.GetGID(gs.name)
	if err != nil {
		return fmt.Errorf("failed to GetGID expected group %s: %w", gs.name, err)
	}
	if d != gs.gid {
		return fmt.Errorf("wrong GID for %s: got %d, expected %d", gs.name, d, gs.gid)
	}

	s, err = g.GetGroup(gs.gid)
	if err != nil {
		return fmt.Errorf("failed to GetGroup expected group %s: %w", gs.name, err)
	}
	if s != gs.name {
		return fmt.Errorf("wrong groupname for %d: got %s, expected %s", gs.gid, s, gs.name)
	}
	return nil
}

func TestGroups(t *testing.T) {
	t.Run("non-existent group file", func(t *testing.T) {
		f := "testdata/does-not-exist"
		g, e := NewGroups(f)
		if e == nil {
			t.Errorf("NewGroups jnon-existant file should return an error")
		}
		if g == nil {
			t.Errorf("NewGroups should return a valid Groups object, even on error")
		}
	})
	t.Run("empty group file", func(t *testing.T) {
		f := "testdata/group-empty.txt"
		g, e := NewGroups(f)
		if e != nil {
			t.Errorf("NewGroups should not report error for empty passwd file")
		}
		if g == nil {
			t.Errorf("NewGroups should return a valid Users object even if passwd file is empty")
		}
	})
	t.Run("almost empty group file", func(t *testing.T) {
		f := "testdata/group-newline.txt"
		g, e := NewGroups(f)
		if e != nil {
			t.Errorf("NewGroups should not report error for empty passwd file")
		}
		if g == nil {
			t.Errorf("NewGroups should return a valid Users object even if passwd file is empty")
		}
	})
	for _, f := range groupFiles {
		t.Run(f, func(t *testing.T) {
			g, e := NewGroups(f)
			if e != nil {
				t.Errorf("NewGroups should not return an error on valid file")
			}
			if g == nil {
				t.Errorf("NewGroups should return a valid Users object on valid file")
			}
			for _, gs := range groupShould {
				if e := groupSame(g, gs); e != nil {
					t.Errorf("%v", e)
				}
			}
			for _, gs := range groupShouldnt {
				if e := groupSame(g, gs); e == nil {
					t.Errorf("group %s matched when it shouldn't", gs.name)
				}
			}
			for u, is := range groupMembers {
				js := g.UserGetGIDs(u)
				if len(js) != len(is) {
					t.Errorf("unequal gid lists for %s: %v vs %v", u, is, js)
				} else {
					sort.Ints(is)
					sort.Ints(js)
					for i := range is {
						if is[i] != js[i] {
							t.Errorf("unequal gid lists for %s: %v vs %v", u, is, js)
						}
					}
				}
			}
		})
	}
}

func TestFlags(t *testing.T) {
	if !correctFlags(true, false, false) {
		t.Errorf("correctFLags(true, false, false): got ! ok, want ok")
	}
	if correctFlags(true, true, false) {
		t.Errorf("correctFLags(true, true, false): got ! ok, want ok")
	}

	tests := []struct {
		testFlags *flags
		wantError error
	}{
		{
			testFlags: &flags{name: true},
			wantError: errNotOnlyNames,
		},
		{
			testFlags: &flags{real: true},
			wantError: errNotOnlyNamesOrIDs,
		},
		{
			testFlags: &flags{group: true, groups: true},
			wantError: errOnlyOneChoice,
		},
		{
			testFlags: &flags{groups: true, user: true},
			wantError: errOnlyOneChoice,
		},
		{
			testFlags: &flags{group: true, user: true},
			wantError: errOnlyOneChoice,
		},
		{
			testFlags: &flags{group: true, groups: true, user: true},
			wantError: errOnlyOneChoice,
		},
	}

	for _, tt := range tests {
		if err := run(io.Discard, "r", tt.testFlags, "", ""); !errors.Is(err, tt.wantError) {
			t.Errorf(`run(%v, "", ""): got %v, want %v`, tt.testFlags, err, tt.wantError)
		}
	}
}
