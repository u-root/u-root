// Copyright 2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package loop

import (
	"os"
	"path/filepath"
	"syscall"
	"testing"

	"github.com/hugelgupf/vmtest/guest"
	"github.com/u-root/u-root/pkg/core/cp"
	"golang.org/x/sys/unix"
)

const (
	_LOOP_MAJOR = 7
)

func TestFindDevice(t *testing.T) {
	guest.SkipIfNotInVM(t)

	loopdev, err := FindDevice()
	if err != nil {
		t.Fatalf("Failed to find loop device: %v", err)
	}

	s, err := os.Stat(loopdev)
	if err != nil {
		t.Fatalf("Could not stat loop device: %v", err)
	}

	st := s.Sys().(*syscall.Stat_t)
	if m := unix.Major(st.Rdev); m != _LOOP_MAJOR {
		t.Fatalf("Device %s is not a loop device: got major no %d, want %d", loopdev, m, _LOOP_MAJOR)
	}
}

func TestSetFile(t *testing.T) {
	guest.SkipIfNotInVM(t)

	tmpDir := t.TempDir()
	testdisk := filepath.Join(tmpDir, "testdisk")
	if err := cp.Copy("./testdata/pristine-vfat-disk", testdisk); err != nil {
		t.Fatal(err)
	}

	loopdev, err := New(testdisk, "vfat", "")
	if err != nil {
		t.Fatal(err)
	}

	if err := os.MkdirAll("/tmp/disk", 0o755); err != nil && !os.IsExist(err) {
		t.Fatalf("Could not create /tmp/disk: %v", err)
	}

	mp, err := loopdev.Mount("/tmp/disk", 0)
	if err != nil {
		t.Fatalf("Failed to mount /tmp/disk: %v", err)
	}
	defer mp.Unmount(0) //nolint:errcheck

	if err := os.WriteFile("/tmp/disk/foobar", []byte("Are you feeling it now Mr Krabs"), 0o755); err != nil {
		t.Fatal(err)
	}
}
