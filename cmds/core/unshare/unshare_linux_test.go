// Copyright 2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package main

import (
	"os"
	"slices"
	"syscall"
	"testing"
)

func TestCommand(t *testing.T) {
	tests := []struct {
		args          []string
		expectedArgs  []string
		expectedFlags uintptr
		ipc           bool
		mount         bool
		pid           bool
		net           bool
		uts           bool
		user          bool
	}{
		{
			expectedArgs: []string{"/bin/sh"},
		},
		{
			args:          []string{"echo", "hello"},
			expectedArgs:  []string{"echo", "hello"},
			expectedFlags: syscall.CLONE_NEWNS | syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWNET | syscall.CLONE_NEWIPC | syscall.CLONE_NEWUSER,
			ipc:           true,
			mount:         true,
			pid:           true,
			net:           true,
			uts:           true,
			user:          true,
		},
	}

	for _, test := range tests {
		res := command(test.ipc, test.mount, test.pid, test.net, test.uts, test.user, test.args...)
		if res.SysProcAttr.Cloneflags != test.expectedFlags {
			t.Errorf("expected %v, got %v", test.expectedFlags, res.SysProcAttr.Cloneflags)
		}
		if !slices.Equal(res.Args, test.expectedArgs) {
			t.Errorf("expected equal got, %v != %v", res.Args, test.expectedArgs)
		}
		if res.Stdin != os.Stdin || res.Stdout != os.Stdout || res.Stderr != os.Stderr {
			t.Errorf("expected stdin, stdout, stderr, got %v, %v, %v", res.Stdin, res.Stdout, res.Stderr)
		}
	}
}
