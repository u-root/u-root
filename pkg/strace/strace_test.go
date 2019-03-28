// Copyright 2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package strace

import (
	"os/exec"
	"testing"
)

// It's not really easy to write a full up general tester for this.
// Even the simplest commands on Linux have dozens of system calls.
// The Go assembler should in principle let us write a 3 line assembly
// program that just does an exit system call:
// MOVQ $exit, RARG
// SYSCALL
// But that's for someone else to do :-)

func TestNoCommandFail(t *testing.T) {
	Debug = t.Logf
	c, err := New()
	if err != nil {
		t.Fatal(err)
	}
	c.Raw = true
	go c.RunTracerFromCmd(exec.Command("hi", "/etc/hosts"))
	r := <-c.Records
	if r.Err == nil {
		t.Fatalf("Got nil, want a non-nil error")
	}
}

func TestBasicStrace(t *testing.T) {
	Debug = t.Logf
	c, err := New()
	if err != nil {
		t.Fatal(err)
	}

	go c.RunTracerFromCmd(exec.Command("ls", "/etc/hosts"))
	for range c.Records {

	}
}
