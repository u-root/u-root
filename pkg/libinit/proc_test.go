// Copyright 2014-2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package libinit

import (
	"bytes"
	"os/exec"
	"reflect"
	"strings"
	"testing"
)

func TestWithArguments(t *testing.T) {
	tests := []struct {
		name string
		args []string
		want []string
	}{
		{
			name: "No arguments",
			args: []string{},
			want: nil,
		},
		{
			name: "Single argument",
			args: []string{"arg1"},
			want: []string{"arg1"},
		},
		{
			name: "Multiple arguments",
			args: []string{"arg1", "arg2"},
			want: []string{"arg1", "arg2"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &exec.Cmd{}
			modifier := WithArguments(tt.args...)
			modifier(cmd)
			if !reflect.DeepEqual(cmd.Args, tt.want) {
				t.Errorf("got %v, want %v", cmd.Args, tt.want)
			}
		})
	}
}

func TestWithStdin(t *testing.T) {
	cmd := &exec.Cmd{}
	input := "input data"
	r := strings.NewReader(input)
	modifier := WithStdin(r)
	modifier(cmd)
	if cmd.Stdin != r {
		t.Errorf("got %v, want %v", cmd.Stdin, r)
	}
}

func TestWithStdout(t *testing.T) {
	cmd := &exec.Cmd{}
	var w bytes.Buffer
	modifier := WithStdout(&w)
	modifier(cmd)
	if cmd.Stdout != &w {
		t.Errorf("got %v, want %v", cmd.Stdout, &w)
	}
}

func TestWithStderr(t *testing.T) {
	cmd := &exec.Cmd{}
	var w bytes.Buffer
	modifier := WithStderr(&w)
	modifier(cmd)
	if cmd.Stderr != &w {
		t.Errorf("got %v, want %v", cmd.Stderr, &w)
	}
}
