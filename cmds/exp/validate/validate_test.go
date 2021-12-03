// Copyright 2016 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"io"
	"os"
	"path/filepath"
	"syscall"
	"testing"

	"github.com/u-root/u-root/pkg/testutil"
)

type file struct {
	name string
	a    string
	val  []byte
	o    string
	e    string
	x    int // XXX wrong for Plan 9 and Harvey
}

func TestValidate(t *testing.T) {
	data := []byte(`127.0.0.1	localhost
127.0.1.1	akaros
192.168.28.16	ak
192.168.28.131	uroot

# The following lines are desirable for IPv6 capable hosts
::1     localhost ip6-localhost ip6-loopback
ff02::1 ip6-allnodes
ff02::2 ip6-allrouters
`)
	tests := []file{
		{name: "hosts.sha1", val: []byte("3f397a3b3a7450075da91b078afa35b794cf6088  hosts"), o: "SHA1\n"},
	}

	tmpDir := t.TempDir()
	if err := os.WriteFile(filepath.Join(tmpDir, "hosts"), data, 0o444); err != nil {
		t.Fatalf("Can't set up data file: %v", err)
	}

	for _, v := range tests {
		if err := os.WriteFile(filepath.Join(tmpDir, v.name), v.val, 0o444); err != nil {
			t.Fatalf("Can't set up hash file: %v", err)
		}

		c := testutil.Command(t, filepath.Join(tmpDir, v.name), filepath.Join(tmpDir, "hosts"))
		ep, err := c.StderrPipe()
		if err != nil {
			t.Fatalf("Can't start StderrPipe: %v", err)
		}
		op, err := c.StdoutPipe()
		if err != nil {
			t.Fatalf("Can't start StdoutPipe: %v", err)
		}

		if err := c.Start(); err != nil {
			t.Fatalf("Can't start %v: %v", c, err)
		}
		e, err := io.ReadAll(ep)
		if err != nil {
			t.Fatalf("Can't get stderr of %v: %v", c, err)
		}
		o, err := io.ReadAll(op)
		if err != nil {
			t.Fatalf("Can't get sdout of %v: %v", c, err)
		}

		if err = c.Wait(); err != nil {
			t.Fatalf("Can's Wait %v: %v", c, err)
		}

		// TODO: fix this for Plan 9/Harvey
		s := c.ProcessState.Sys().(syscall.WaitStatus).ExitStatus()

		if s != v.x {
			t.Errorf("Validate %v hosts %v (%v): want (exit: %v), got (exit %v), output %v", v.a, v.name, string(v.val), v.x, s, string(o))
			continue
		}

		if err != nil && string(e) != v.e {
			t.Errorf("Validate %v hosts %v (%v): want stderr: %v, got %v)", v.a, v.name, string(v.val), v.e, string(o))
			continue
		}

		if string(o) != v.o {
			t.Errorf("Validate %v hosts %v (%v): want stdout: %v, got %v)", v.a, v.name, string(v.val), v.o, string(o))
			continue
		}

		t.Logf("Validate %v hosts %v: %v", v.a, v.name, string(o))
	}
}

func TestMain(m *testing.M) {
	testutil.Run(m, main)
}
