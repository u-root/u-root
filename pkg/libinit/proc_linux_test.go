// Copyright 2014-2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package libinit

import (
	"bytes"
	"errors"
	"io"
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
	tests := []struct {
		name      string
		mtty      bool
		openFn    func([]string) ([]io.Writer, error)
		ttyNames  []string
		expectErr bool
	}{
		{
			name: "MultiTTY enabled with no writers",
			mtty: true,
			openFn: func([]string) ([]io.Writer, error) {
				return []io.Writer{}, nil
			},
			ttyNames:  nil,
			expectErr: true,
		},
		{
			name: "MultiTTY enabled with single writer",
			mtty: true,
			openFn: func([]string) ([]io.Writer, error) {
				return []io.Writer{&bytes.Buffer{}}, nil
			},
			ttyNames:  []string{"tty1"},
			expectErr: false,
		},
		{
			name: "MultiTTY enabled with multiple writers",
			mtty: true,
			openFn: func([]string) ([]io.Writer, error) {
				return []io.Writer{&bytes.Buffer{}, &bytes.Buffer{}}, nil
			},
			ttyNames:  []string{"tty1", "tty2"},
			expectErr: false,
		},
		{
			name: "MultiTTY enabled with openFn returning error",
			mtty: true,
			openFn: func([]string) ([]io.Writer, error) {
				return nil, errors.New("failed to open TTY devices")
			},
			ttyNames:  []string{"tty1"},
			expectErr: true,
		},
		{
			name: "MultiTTY disabled",
			mtty: false,
			openFn: func([]string) ([]io.Writer, error) {
				return []io.Writer{&bytes.Buffer{}}, nil
			},
			ttyNames:  []string{"tty1"},
			expectErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := exec.Command("echo", "test")
			modifier := WithMultiTTY(tt.mtty, tt.openFn, tt.ttyNames)
			modifier(cmd)

			if tt.mtty {
				if tt.expectErr {
					if cmd.Stdout != nil || cmd.Stderr != nil {
						t.Errorf("expected fallback to default stdout and stderr, got %v and %v", cmd.Stdout, cmd.Stderr)
					}
				} else {
					if cmd.Stdout == nil || cmd.Stderr == nil {
						t.Errorf("expected stdout and stderr to be set, got %v and %v", cmd.Stdout, cmd.Stderr)
					}
				}
			} else {
				if cmd.Stdout != nil || cmd.Stderr != nil {
					t.Errorf("expected no writers to be set, got %v and %v", cmd.Stdout, cmd.Stderr)
				}
			}
		})
	}
}
