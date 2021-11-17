// Copyright 2021 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"os"
	"os/exec"
	"testing"

	"src.elv.sh/pkg/prog"
)

func TestElvish(t *testing.T) {
	for _, tt := range []struct {
		name   string
		status int
	}{
		{
			name:   "success",
			status: 0,
		},
		{
			name:   "failure",
			status: 1,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			if os.Getenv("TEST_ELVISH") == "1" {
				sh = func(fds [3]*os.File, args []string, programs ...prog.Program) int {
					return tt.status
				}
				main()
				return
			}
			cmd := exec.Command(os.Args[0], "-test.run=TestElvish")
			cmd.Env = append(cmd.Env, "TEST_ELVISH=1")
			err := cmd.Run()
			e, ok := err.(*exec.ExitError)
			if ok && !e.Success() && tt.name == "success" {
				t.Errorf("expected exit code 0, got %d", e.ExitCode())
			}
			if ok && e.Success() && tt.name == "failure" {
				t.Error("expected exit code !0 but got 0")
			}
		})
	}
}

func TestShouldRun(t *testing.T) {
	for _, tt := range []struct {
		name   string
		daemon bool
	}{
		{
			name:   "nodaemon",
			daemon: false,
		},
		{
			name:   "daemon",
			daemon: true,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			f := &prog.Flags{
				Daemon: tt.daemon,
			}
			d := daemonStub{}
			result := d.ShouldRun(f)
			if tt.name == "nodaemon" && result {
				t.Error("expected false, got true")
			}
			if tt.name == "daemon" && !result {
				t.Error("expected true, got false")
			}
		})
	}
}

func TestRun(t *testing.T) {
	d := daemonStub{}
	err := d.Run([3]*os.File{}, &prog.Flags{}, []string{})
	if err == nil {
		t.Error("expected error, got nil")
	}
}
