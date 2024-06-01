// Copyright 2016 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
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
	}{
		{
			name: "not enough args",
			args: []string{"kill"},
			want: fmt.Sprintf("%s\n", eUsage),
		},
		{
			name: "list signal names with too much args",
			args: []string{"kill", "-l", "10"},
			want: fmt.Sprintf("%s\n", eUsage),
		},
		{
			name: "list signal names",
			args: []string{"kill", "-l"},
			want: fmt.Sprintf("%s\n", siglist()),
		},
		{
			name: "kill pid 2",
			args: []string{"kill", "2"},
			want: "",
		},
		{
			name: "kill signal without signal and pid",
			args: []string{"kill", "-s"},
			want: fmt.Sprintf("%s\n", eUsage),
		},
		{
			name: "kill signal with signal but without pid",
			args: []string{"kill", "--signal", "9"},
			want: fmt.Sprintf("%s\n", eUsage),
		},
		{
			name: "kill signal with signal and pid",
			args: []string{"kill", "--signal", "50", "2"},
			want: "",
		},
		{
			name: "kill signal with signal and wrong pid",
			args: []string{"kill", "--signal", "9", getUnusedPID()},
			want: "some processes could not be killed",
		},
		{
			name: "signal is invalid",
			args: []string{"kill", "--signal", "a"},
			want: "is not a valid signal",
		},
		{
			name: "signal is invalid",
			args: []string{"kill", "-1", "a"},
			want: "some processes could not be killed: a: arguments must be process or job IDS",
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			buf := &bytes.Buffer{}
			if err := killProcess(buf, tt.args...); err != nil {
				if !strings.Contains(err.Error(), tt.want) {
					t.Errorf("killProcess() = %q, want to contain: %q", err.Error(), tt.want)
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
