// Copyright 2016 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"errors"
	"fmt"
	"os/exec"
	"testing"
)

func getUnusedPID() string {
	cmd := exec.Command("true")
	cmd.Start()
	pid := cmd.Process.Pid
	cmd.Wait()
	return fmt.Sprintf("%d", pid)
}

func TestKillProcess(t *testing.T) {
	for _, tt := range []struct {
		name string
		args []string
		want string
		err  error
	}{
		{
			name: "not enough args",
			args: []string{"kill"},
			err:  errUsage,
		},
		{
			name: "list signal names with too much args",
			args: []string{"kill", "-l", "10"},
			err:  errUsage,
		},
		{
			name: "list signal names",
			args: []string{"kill", "-l"},
			want: fmt.Sprintf("%s\n", siglist()),
		},
		{
			name: "kill signal without signal and pid",
			args: []string{"kill", "-s"},
			err:  errUsage,
		},
		{
			name: "kill signal with signal but without pid",
			args: []string{"kill", "--signal", "9"},
			err:  errUsage,
		},
		{
			name: "kill signal with signal and wrong pid",
			args: []string{"kill", "--signal", "9", getUnusedPID()},
			err:  errCannotKill,
		},
		{
			name: "signal is invalid",
			args: []string{"kill", "--signal", "a"},
			err:  errInvalidSignal,
		},
		{
			name: "signal is invalid",
			args: []string{"kill", "-1", "a"},
			err:  errCannotKill,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			buf := &bytes.Buffer{}
			if err := killProcess(buf, tt.args...); err != nil {
				if !errors.Is(err, tt.err) {
					t.Errorf("killProcess() = %v, want %v", err, tt.err)
				}
			} else {
				if buf.String() != tt.want {
					t.Errorf("killProcess() = %q, want: %q", buf.String(), tt.want)
				}
			}
		})
	}
}

func TestSignalList(t *testing.T) {
	if len(signames)*2 != len(signums) {
		t.Error("len(signames) * 2 != len(signums)")
	}
}
