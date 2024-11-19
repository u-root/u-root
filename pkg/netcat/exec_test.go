// Copyright 2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package netcat

import (
	"bytes"
	"testing"
)

func TestParseCommands(t *testing.T) {
	tests := []struct {
		name        string
		execs       []Exec
		expected    Exec
		expectError bool
	}{
		{
			name: "Single valid command",
			execs: []Exec{
				{Type: EXEC_TYPE_SHELL, Command: "echo 'Hello, world!'"},
			},
			expected:    Exec{Type: EXEC_TYPE_SHELL, Command: "echo 'Hello, world!'"},
			expectError: false,
		},
		{
			name: "Multiple commands with one valid",
			execs: []Exec{
				{Type: EXEC_TYPE_SHELL, Command: ""},
				{Type: EXEC_TYPE_LUA, Command: "echo 'Second command'"},
			},
			expected:    Exec{Type: EXEC_TYPE_LUA, Command: "echo 'Second command'"},
			expectError: false,
		},
		{
			name: "Multiple valid commands",
			execs: []Exec{
				{Type: EXEC_TYPE_SHELL, Command: "echo 'First command'"},
				{Type: EXEC_TYPE_LUA, Command: "echo 'Second command'"},
			},
			expected:    Exec{},
			expectError: true,
		},
		{
			name:        "No commands",
			execs:       []Exec{{EXEC_TYPE_NONE, ""}},
			expected:    Exec{EXEC_TYPE_NONE, ""},
			expectError: false,
		},
		{
			name: "All empty commands",
			execs: []Exec{
				{Type: EXEC_TYPE_SHELL, Command: ""},
				{Type: EXEC_TYPE_LUA, Command: ""},
			},
			expected:    Exec{EXEC_TYPE_NONE, ""},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseCommands(tt.execs...)
			if (err != nil) != tt.expectError {
				t.Errorf("ParseCommands() error = %v, wantErr %v", err, tt.expectError)
				return
			}
			if !tt.expectError && (got.Type != tt.expected.Type || got.Command != tt.expected.Command) {
				t.Errorf("ParseCommands() got = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestExec_Execute(t *testing.T) {
	tests := []struct {
		name    string
		exec    Exec
		stdin   string
		eol     []byte
		wantErr bool
	}{
		{
			name: "empty command",
			exec: Exec{
				Type:    EXEC_TYPE_NATIVE,
				Command: "",
			},
			stdin:   "",
			eol:     []byte("\n"),
			wantErr: true,
		},
		{
			name: "Execute shell command successfully",
			exec: Exec{
				Type:    EXEC_TYPE_SHELL,
				Command: "echo 'Hello, Shell!'",
			},
			stdin:   "",
			eol:     []byte("\n"),
			wantErr: false,
		},
		{
			name: "Lua execution not implemented",
			exec: Exec{
				Type:    EXEC_TYPE_LUA,
				Command: "lua cmd",
			},
			stdin:   "",
			eol:     []byte("\n"),
			wantErr: true,
		},
		{
			name: "Invalid exec type",
			exec: Exec{
				Type:    999, // Invalid type
				Command: "echo 'Invalid Type!'",
			},
			stdin:   "",
			eol:     []byte("\n"),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stdout := &bytes.Buffer{}
			stderr := &bytes.Buffer{}

			err := tt.exec.Execute(stdout, stderr, tt.eol)
			if (err != nil) != tt.wantErr {
				t.Errorf("Exec.Execute() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
