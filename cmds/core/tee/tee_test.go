// Copyright 2022 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//go:build !tinygo || tinygo.enable

package main

import (
	"bytes"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func TestTee(t *testing.T) {
	tests := []struct {
		name          string
		input         string
		args          []string
		append        bool
		appendContent string
	}{
		{
			name:  "default tee",
			input: "hello",
			args:  []string{"a1", "a2"},
		},
		{
			name:          "with append flag",
			input:         "a\nb\n",
			args:          []string{"b1"},
			append:        true,
			appendContent: "hello",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			tempDir := t.TempDir()

			for i := 0; i < len(test.args); i++ {
				test.args[i] = filepath.Join(tempDir, test.args[i])
			}

			if test.append {
				for _, arg := range test.args {
					err := os.WriteFile(arg, []byte(test.appendContent), 0o666)
					if err != nil {
						t.Fatal(err)
					}
				}
			}

			var stdout bytes.Buffer
			var stderr bytes.Buffer
			// cmd := command(test.append, false, test.args)
			cmd := &cmd{
				stdin:  strings.NewReader(test.input),
				stdout: &stdout,
				stderr: &stderr,
				args:   test.args,
				cat:    test.append,
			}
			if err := cmd.run(); err != nil {
				t.Fatal(err)
			}

			// test if stdin match stdout
			if test.input != stdout.String() {
				t.Errorf("wanted: %q, got: %q", test.input, stdout.String())
			}

			for _, name := range test.args {
				b, err := os.ReadFile(name)
				if err != nil {
					t.Error(err)
				}
				res := string(b)
				expectedContent := test.input

				if test.append {
					expectedContent = test.appendContent + expectedContent
				}

				if res != expectedContent {
					t.Errorf("wanted: %q, got %q", expectedContent, res)
				}
			}
		})
	}
}

func TestCommand(t *testing.T) {
	stdin := os.Stdin
	stdout := os.Stdout
	stderr := os.Stderr

	tests := []struct {
		name    string
		args    []string
		wantCmd *cmd
	}{
		{
			name: "Append mode short flag",
			args: []string{"tee", "-a"},
			wantCmd: &cmd{
				stdin:  stdin,
				stdout: stdout,
				stderr: stderr,
				cat:    true,
				ignore: false,
				args:   []string{},
			},
		},
		{
			name: "Ignore interrupts long flag",
			args: []string{"tee", "--ignore-interrupts"},
			wantCmd: &cmd{
				stdin:  stdin,
				stdout: stdout,
				stderr: stderr,
				cat:    false,
				ignore: true,
				args:   []string{},
			},
		},
		{
			name: "Append and ignore interrupts",
			args: []string{"tee", "-a", "-i"},
			wantCmd: &cmd{
				stdin:  stdin,
				stdout: stdout,
				stderr: stderr,
				cat:    true,
				ignore: true,
				args:   []string{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotCmd := command(tt.args)
			if !reflect.DeepEqual(gotCmd, tt.wantCmd) {
				t.Errorf("%s: command() = %+v, want %+v", tt.name, gotCmd, tt.wantCmd)
			}
		})
	}
}
