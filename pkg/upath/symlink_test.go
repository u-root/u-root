// Copyright 2015-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package upath

import (
	"log"
	"os"
	"path/filepath"
	"testing"
)

func TestSymlink(t *testing.T) {
	td := t.TempDir()
	for _, n := range []string{"bin", "buildbin"} {
		p := filepath.Join(td, n)
		if err := os.Mkdir(p, 0o777); err != nil {
			log.Fatal(err)
		}
	}
	tab := []struct {
		s, t, v string
	}{
		{filepath.Join(td, "bin/ash"), "sh", filepath.Join(td, "buildbin/gosh")},
		{filepath.Join(td, "bin/sh"), "../buildbin/gosh", filepath.Join(td, "buildbin/gosh")},
		{filepath.Join(td, "buildbin/gosh"), "installcommand", filepath.Join(td, "buildbin/gosh")},
	}
	for _, s := range tab {
		if err := os.Symlink(s.t, s.s); err != nil {
			t.Fatalf("symlink %s->%s: got %v, want nil", s.s, s.t, err)
		}
	}
	for _, s := range tab {
		t.Logf("Check %v", s)
		v, err := os.Readlink(s.s)
		t.Logf("Symlink val %v", v)
		if err != nil || v != s.t {
			t.Errorf("readlink %v: got (%v, %v), want (%v, nil)", s.s, v, err, s.t)
		}
		v = ResolveUntilLastSymlink(s.s)
		t.Logf("ResolveUntilLastSymlink val %v", v)
		if v != s.v {
			t.Errorf("ResolveUntilLastSymlink %v: got %v want %v", s.s, v, s.v)
		}
	}
	// test to make sure a plain file gives a reasonable result.
	ic := filepath.Join(td, "x")
	if err := os.WriteFile(ic, nil, 0o666); err != nil {
		t.Fatalf("WriteFile %v: got %v, want nil", ic, err)
	}
	v := ResolveUntilLastSymlink(ic)
	t.Logf("ResolveUntilLastSymlink %v gets %v", ic, v)
	if v != ic {
		t.Errorf("ResolveUntilLastSymlink %v: got %v want %v", ic, v, ic)
	}
}
