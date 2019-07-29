// Copyright 2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cp_test

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/u-root/u-root/pkg/cp"
	"github.com/u-root/u-root/pkg/cp/cmp"
)

func copyAndTest(t *testing.T, o cp.Options, src, dst string) {
	if err := o.Copy(src, dst); err != nil {
		t.Fatalf("Copy(%q -> %q) = %v, want %v", src, dst, err, nil)
	}
	if err := cmp.IsEqualTree(o, src, dst); err != nil {
		t.Fatalf("Expected %q and %q to be same, got %v", src, dst, err)
	}
}

func TestSimpleCopy(t *testing.T) {
	tmpDir, err := ioutil.TempDir("", "u-root-pkg-cp-")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// Copy a directory.
	origd := filepath.Join(tmpDir, "directory")
	if err := os.Mkdir(origd, 0744); err != nil {
		t.Fatal(err)
	}

	copyAndTest(t, cp.Default, origd, filepath.Join(tmpDir, "directory-copied"))
	copyAndTest(t, cp.NoFollowSymlinks, origd, filepath.Join(tmpDir, "directory-copied-2"))

	// Copy a file.
	origf := filepath.Join(tmpDir, "normal-file")
	if err := ioutil.WriteFile(origf, []byte("F is for fire that burns down the whole town"), 0766); err != nil {
		t.Fatal(err)
	}

	copyAndTest(t, cp.Default, origf, filepath.Join(tmpDir, "normal-file-copied"))
	copyAndTest(t, cp.NoFollowSymlinks, origf, filepath.Join(tmpDir, "normal-file-copied-2"))

	// Copy a symlink.
	origs := filepath.Join(tmpDir, "foobar")
	// foobar -> normal-file
	if err := os.Symlink(origf, origs); err != nil {
		t.Fatal(err)
	}

	copyAndTest(t, cp.Default, origf, filepath.Join(tmpDir, "foobar-copied"))
	copyAndTest(t, cp.NoFollowSymlinks, origf, filepath.Join(tmpDir, "foobar-copied-just-symlink"))
}
