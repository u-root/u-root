// Copyright 2025 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build !race

package integration

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestCPUAMD64 tests both general and specific things. The specific parts are
// the io and cmos commands.
func TestCPUAMD64(t *testing.T) {
	t.Parallel()

	g := bootGuest(t, "amd64")

	// Files can be reached from inside and outside the guest without adding
	// them to the initramfs or reading them off the console.
	if err := os.WriteFile(filepath.Join(g.Dir(), "a"), []byte("hi"), 0o644); err != nil {
		t.Fatal(err)
	}

	for _, tt := range []struct {
		cmd  string
		args []string
		ok   bool
	}{
		{cmd: "/bbin/dd", args: []string{"if=/dev/x"}, ok: false},
		{cmd: "/bbin/dd", args: []string{"if=/dev/null"}, ok: true},
		{cmd: "/bbin/dd", args: []string{"if=/tmp/cpu/a", "of=/tmp/cpu/b"}, ok: true},
	} {
		b, err := g.Run(tt.cmd, tt.args...)
		if (err == nil) != tt.ok {
			t.Errorf("%s %s: got %v, want ok=%v", tt.cmd, tt.args, err, tt.ok)
		}
		t.Logf("%q", string(b))
	}

	b, err := os.ReadFile(filepath.Join(g.Dir(), "b"))
	if err != nil {
		t.Fatalf("reading b: got %v, want nil", err)
	}
	if string(b) != "hi" {
		t.Fatalf("file b: got %q, want %q", b, "hi")
	}

	for _, tt := range []struct {
		args string
		out  string
	}{
		{args: "cw 14 1", out: ""},
		{args: "cr 14", out: "0x01\n"},
		{args: "cw 14 0", out: ""},
		{args: "cr 14", out: "0x00\n"},
	} {
		b, err := g.Run("/bbin/io", strings.Split(tt.args, " ")...)
		if err != nil {
			t.Errorf("io %s: got %v, want nil", tt.args, err)
		}
		if string(b) != tt.out {
			t.Errorf("io %s: got %v, want %v", tt.args, string(b), tt.out)
		}
		t.Logf("io %s = %q", tt.args, string(b))
	}

	// The io integration tests include writing to 3f8. There's no need to do
	// that, the cmos write tests all that needs testing, as it uses inb and
	// outb, and uart hardware is fickle. The test above is enough for the IO
	// package.
}
