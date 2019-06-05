// Copyright 2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/u-root/u-root/pkg/testutil"
	"golang.org/x/sys/unix"
)

// Test major and minor numbers greater then 255.
//
// This is supported since Linux 2.6. The major/minor numbers used for this
// test are (1110, 74616). According to "kdev_t.h":
//
//       mkdev(1110, 74616)
//     = mkdev(0x456, 0x12378)
//     = (0x12378 & 0xff) | (0x456 << 8) | ((0x12378 & ~0xff) << 12)
//     = 0x12345678
func TestLargeDevNumber(t *testing.T) {
	if uid := os.Getuid(); uid != 0 {
		t.Skipf("test requires root, your uid is %d", uid)
	}

	// Make a temporary directory.
	tmpDir, err := ioutil.TempDir("", "ls")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)
	file := filepath.Join(tmpDir, "large_node")

	// Run "mknod large_node b 1110 74616".
	err = testutil.Command(t, file, "b", "1110", "74616").Run()
	if err != nil {
		t.Fatal(err)
	}

	// Check the device number.
	var s unix.Stat_t
	err = unix.Stat(file, &s)
	if err != nil {
		t.Fatal(err)
	}
	if s.Rdev != 0x12345678 {
		t.Fatalf("expected the device number to be 0x12345678, got %#x", s.Rdev)
	}
}
