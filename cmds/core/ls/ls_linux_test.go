// Copyright 2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"path/filepath"
	"strings"
	"testing"

	"github.com/hugelgupf/vmtest/guest"
	"github.com/u-root/u-root/pkg/ls"
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

// Test listName func
func TestListNameLinux(t *testing.T) {
	guest.SkipIfNotInVM(t)

	// Create a directory
	d := t.TempDir()
	if err := unix.Mknod(filepath.Join(d, "large_node"), 0o660|unix.S_IFBLK, 0x12345678); err != nil {
		t.Fatalf("err in unix.Mknod: %v", err)
	}

	var c cmd
	// Setting the flags
	c.long = true
	// Running the tests
	// Write output in buffer.
	var buf bytes.Buffer
	var s ls.Stringer = ls.NameStringer{}
	if c.quoted {
		s = ls.QuotedStringer{}
	}
	if c.long {
		s = ls.LongStringer{Human: c.human, Name: s}
	}
	c.w = &buf
	_ = c.listName(s, d, false)
	if !strings.Contains(buf.String(), "1110, 74616") {
		t.Errorf("Expected value: %s, got: %s", "1110, 74616", buf.String())
	}
}
