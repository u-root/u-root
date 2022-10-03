// Copyright 2022 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package efivarfs

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/u-root/u-root/pkg/testutil"
	"golang.org/x/sys/unix"
)

func TestFSGoodFile(t *testing.T) {
	// Temporary folder cleanup can require root.
	testutil.SkipIfNotRoot(t)
	d := t.TempDir()
	f, err := os.Create(filepath.Join(d, "x"))
	if err != nil {
		t.Fatalf("os.Create(%s): %v != nil", filepath.Join(d, "x"), err)
	}
	i, err := getInodeFlags(f)
	if err != nil {
		t.Skipf("Can not getInodeFlags: %v != nil", err)
	}

	if err := setInodeFlags(f, i); err != nil {
		t.Fatalf("setInodeFlags: %v != nil", err)
	}

	restore, err := makeMutable(f)
	if err != nil {
		t.Fatalf("makeMutable: %v != nil", err)
	}
	if restore == nil {
		t.Logf("it was not mutable to start")
	}

	i |= unix.STATX_ATTR_IMMUTABLE
	if err := setInodeFlags(f, i); err != nil {
		t.Skipf("Skipping rest of test, unable to set immutable flag")
	}
}

func TestFSBadFile(t *testing.T) {
	f, err := os.Open("/dev/null")
	if err != nil {
		t.Fatalf("os.Open(/dev/null): %v != nil", err)
	}
	i, err := getInodeFlags(f)
	if err == nil {
		t.Fatalf("getInodeFlags: nil != an error")
	}

	if err := setInodeFlags(f, i); err == nil {
		t.Fatalf("setInodeFlags: nil != an error")
	}

	if _, err := makeMutable(f); err == nil {
		t.Fatalf("makeMutable: nil != some error")
	}
}
