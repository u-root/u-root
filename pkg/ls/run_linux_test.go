// Copyright 2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ls

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/hugelgupf/vmtest/guest"
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
	var s Stringer = NameStringer{}
	if c.quoted {
		s = QuotedStringer{}
	}
	if c.long {
		s = LongStringer{Human: c.human, Name: s}
	}
	c.w = &buf
	_ = c.listName(s, d, false)
	if !strings.Contains(buf.String(), "1110, 74616") {
		t.Errorf("Expected value: %s, got: %s", "1110, 74616", buf.String())
	}
}

func TestNotExist(t *testing.T) {
	d := t.TempDir()
	b := &bytes.Buffer{}
	var c = cmd{w: b}
	if err := c.listName(NameStringer{}, filepath.Join(d, "b"), false); err != nil {
		t.Fatalf("listName(NameString{}, %q/b, w, false): nil != %v", d, err)
	}
	// yeesh.
	// errors not consistent and ... the error has this gratuitous 'lstat ' in front
	// of the filename ...
	eexist := fmt.Sprintf("%s:%v", filepath.Join(d, "b"), os.ErrNotExist)
	enoent := fmt.Sprintf("%s: %v", filepath.Join(d, "b"), unix.ENOENT)
	if !strings.Contains(b.String(), eexist) && !strings.Contains(b.String(), enoent) {
		t.Fatalf("ls of bad name: %q does not contain %q or %q", b.String(), eexist, enoent)
	}
}
