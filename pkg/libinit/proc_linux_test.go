// Copyright 2014-2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package libinit

import (
	"os"
	"os/exec"
	"reflect"
	"testing"

	"golang.org/x/sys/unix"
)

func TestWithTTYControl(t *testing.T) {
	tests := []struct {
		name string
		ctty bool
		want *unix.SysProcAttr
	}{
		{
			name: "Set controlling TTY",
			ctty: true,
			want: &unix.SysProcAttr{
				Setctty: true,
				Setsid:  true,
			},
		},
		{
			name: "Do not set controlling TTY",
			ctty: false,
			want: &unix.SysProcAttr{
				Setctty: false,
				Setsid:  false,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &exec.Cmd{}
			modifier := WithTTYControl(tt.ctty)
			modifier(cmd)
			if got := cmd.SysProcAttr; !reflect.DeepEqual(got, tt.want) {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithCloneFlags(t *testing.T) {
	tests := []struct {
		name  string
		flags uintptr
		want  *unix.SysProcAttr
	}{
		{
			name:  "Set clone flags",
			flags: unix.CLONE_NEWNS,
			want: &unix.SysProcAttr{
				Cloneflags: unix.CLONE_NEWNS,
			},
		},
		{
			name:  "No clone flags",
			flags: 0,
			want: &unix.SysProcAttr{
				Cloneflags: 0,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &exec.Cmd{}
			modifier := WithCloneFlags(tt.flags)
			modifier(cmd)
			if got := cmd.SysProcAttr; !reflect.DeepEqual(got, tt.want) {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithMultiTTY(t *testing.T) {
	// openFn is part of the modifier signature but is no longer consulted by
	// WithMultiTTY: TTY names are stashed in the command environment and the
	// devices are opened later, in RunCommands. A stub keeps the call valid.
	stubOpen := func([]string) ([]*os.File, error) { return nil, nil }

	tests := []struct {
		name     string
		mtty     bool
		ttyNames []string
		wantEnv  []string
	}{
		{
			name:     "MultiTTY enabled with no TTY names",
			mtty:     true,
			ttyNames: nil,
			wantEnv:  nil,
		},
		{
			name:     "MultiTTY enabled with single TTY",
			mtty:     true,
			ttyNames: []string{"tty1"},
			wantEnv:  []string{"tty0=/dev/tty1"},
		},
		{
			name:     "MultiTTY enabled with multiple TTYs",
			mtty:     true,
			ttyNames: []string{"tty1", "tty2"},
			wantEnv:  []string{"tty0=/dev/tty1", "tty1=/dev/tty2"},
		},
		{
			name:     "MultiTTY disabled",
			mtty:     false,
			ttyNames: []string{"tty1"},
			wantEnv:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := exec.Command("echo", "test")
			cmd.Env = nil
			modifier := WithMultiTTY(tt.mtty, stubOpen, tt.ttyNames)
			modifier(cmd)

			if !reflect.DeepEqual(cmd.Env, tt.wantEnv) {
				t.Errorf("cmd.Env = %v, want %v", cmd.Env, tt.wantEnv)
			}

			// WithMultiTTY must not touch the command's I/O streams; the PTY
			// multiplexing wiring happens later, in RunCommands.
			if cmd.Stdout != nil || cmd.Stderr != nil || cmd.Stdin != nil {
				t.Errorf("expected I/O streams to be untouched, got stdin=%v stdout=%v stderr=%v", cmd.Stdin, cmd.Stdout, cmd.Stderr)
			}
		})
	}
}
