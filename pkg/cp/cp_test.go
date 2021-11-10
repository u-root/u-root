// Copyright 2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cp

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func copyAndTest(t *testing.T, o Options, src, dst string) {
	if err := o.Copy(src, dst); err != nil {
		t.Fatalf("Copy(%q -> %q) = %v, want %v", src, dst, err, nil)
	}
	// if err := cmp.IsEqualTree(o, src, dst); err != nil {
	// 	t.Fatalf("Expected %q and %q to be same, got %v", src, dst, err)
	// }
}

func TestSimpleCopy(t *testing.T) {
	tmpDir := t.TempDir()

	// Copy a directory.
	origd := filepath.Join(tmpDir, "directory")
	if err := os.Mkdir(origd, 0o744); err != nil {
		t.Fatal(err)
	}

	copyAndTest(t, Default, origd, filepath.Join(tmpDir, "directory-copied"))
	copyAndTest(t, NoFollowSymlinks, origd, filepath.Join(tmpDir, "directory-copied-2"))

	// Copy a file.
	origf := filepath.Join(tmpDir, "normal-file")
	if err := os.WriteFile(origf, []byte("F is for fire that burns down the whole town"), 0o766); err != nil {
		t.Fatal(err)
	}

	copyAndTest(t, Default, origf, filepath.Join(tmpDir, "normal-file-copied"))
	copyAndTest(t, NoFollowSymlinks, origf, filepath.Join(tmpDir, "normal-file-copied-2"))

	// Copy a symlink.
	origs := filepath.Join(tmpDir, "foobar")
	// foobar -> normal-file
	if err := os.Symlink(origf, origs); err != nil {
		t.Fatal(err)
	}

	copyAndTest(t, Default, origf, filepath.Join(tmpDir, "foobar-copied"))
	copyAndTest(t, NoFollowSymlinks, origf, filepath.Join(tmpDir, "foobar-copied-just-symlink"))
}

func TestCopyTree(t *testing.T) {
	testfiles := make([]*os.File, 3)
	tmpDir, err := ioutil.TempDir("", "u-root-pkg-cp-")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// Make src directory.
	srcd := filepath.Join(tmpDir, "src")
	if err := os.Mkdir(srcd, 0744); err != nil {
		t.Fatal(err)
	}

	if err := os.Chdir(srcd); err != nil {
		t.Fatal(err)
	}

	// Make some rnd files
	for i := 0; i < 3; i++ {
		testfiles[i], err = os.Create("testfile" + fmt.Sprintf("%d", i))
		if err != nil {
			t.Fatal(err)
		}
	}

	// Make dest directory.
	dest := filepath.Join(tmpDir, "dest")
	if err := os.Mkdir(dest, 0744); err != nil {
		t.Fatal(err)
	}
	// Copy the tree
	if err := CopyTree(srcd, dest); err != nil {
		t.Fatal(err)
	}
	// if err := cmp.IsEqualTree(Default, srcd, dest); err != nil {
	// 	t.Fatal(err)
	// }
}
