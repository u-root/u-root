// Copyright 2012-2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/u-root/u-root/pkg/tftp"
)

func TestParseRun(t *testing.T) {
	for _, tt := range []struct {
		name    string
		cmdline []string
		args    []string
		input   []string
		f       tftp.Flags
		exp     string
		err     error
	}{
		{
			name: "SimpleQuit",
			f: tftp.Flags{
				Mode: "ascii",
			},
			cmdline: []string{},
			args:    []string{},
			input:   []string{"localhost", "q"},
		},
		{
			name: "HelpQuit",
			f: tftp.Flags{
				Mode: "ascii",
			},
			cmdline: []string{},
			args:    []string{},
			input:   []string{"localhost", "h", "q"},
		},
		{
			name: "BinaryQuit",
			f: tftp.Flags{
				Mode: "ascii",
			},
			cmdline: []string{},
			args:    []string{},
			input:   []string{"localhost", "binary", "q"},
		},
		{
			name: "AsciiQuit",
			f: tftp.Flags{
				Mode: "ascii",
			},
			cmdline: []string{},
			args:    []string{},
			input:   []string{"localhost", "ascii", "q"},
		},
		{
			name: "ModeQuit",
			f: tftp.Flags{
				Mode: "ascii",
			},
			cmdline: []string{},
			args:    []string{},
			input:   []string{"localhost", "mode", "q"},
			exp:     "Using netascii mode to transfer files.",
		},
		{
			name: "ModeWithValueQuit",
			f: tftp.Flags{
				Mode: "ascii",
			},
			cmdline: []string{},
			args:    []string{},
			input:   []string{"localhost", "mode binary", "q"},
			exp:     "Using octet mode to transfer files.",
		},
		{
			name: "ModeWithValueQuit",
			f: tftp.Flags{
				Mode: "binary",
			},
			cmdline: []string{},
			args:    []string{},
			input:   []string{"localhost", "mode ascii", "q"},
			exp:     "Using netascii mode to transfer files.",
		},
		{
			name: "ModeWithErrorQuit",
			f: tftp.Flags{
				Mode: "binary",
			},
			cmdline: []string{},
			args:    []string{},
			input:   []string{"localhost", "mode error", "q"},
			exp:     fmt.Sprintf("%v", tftp.ErrInvalidTransferMode),
		},
		{
			name: "LiteralQuit",
			f: tftp.Flags{
				Mode: "ascii",
			},
			cmdline: []string{},
			args:    []string{},
			input:   []string{"localhost", "literal", "q"},
			exp:     "Literal mode is on",
		},
		{
			name: "rexmtQuit",
			f: tftp.Flags{
				Mode: "ascii",
			},
			cmdline: []string{},
			args:    []string{},
			input:   []string{"localhost", "rexmt 20", "q"},
		},
		{
			name: "timeoutQuit",
			f: tftp.Flags{
				Mode: "ascii",
			},
			cmdline: []string{},
			args:    []string{},
			input:   []string{"localhost", "timeout 20", "q"},
		},
		{
			name: "traceQuit",
			f: tftp.Flags{
				Mode: "ascii",
			},
			cmdline: []string{},
			args:    []string{},
			input:   []string{"localhost", "trace", "q"},
			exp:     "Packet tracing on.",
		},
		{
			name: "verboseQuit",
			f: tftp.Flags{
				Mode: "ascii",
			},
			cmdline: []string{},
			args:    []string{},
			input:   []string{"localhost", "verbose", "q"},
			exp:     "Verbose mode on.",
		},
		{
			name: "connectQuit",
			f: tftp.Flags{
				Mode: "ascii",
			},
			cmdline: []string{},
			args:    []string{},
			input:   []string{"localhost", "connect localhost 69", "q"},
		},
		{
			name: "statusQuit",
			f: tftp.Flags{
				Mode: "ascii",
			},
			cmdline: []string{},
			args:    []string{},
			input:   []string{"localhost", "status", "q"},
			exp:     "Connected to localhost\nMode: netascii Verbose: off Tracing: off Literal: of",
		},
		{
			name: "IP supplied",
			f: tftp.Flags{
				Mode: "ascii",
			},
			cmdline: []string{"127.0.0.1"},
			args:    []string{"127.0.0.1"},
			input:   []string{"q"},
		},
		{
			name: "IP/Port supplied",
			f: tftp.Flags{
				Mode: "ascii",
			},
			cmdline: []string{"127.0.0.1", "69"},
			args:    []string{"127.0.0.1", "69"},
			input:   []string{"q"},
		},
		{
			name: "Hostname supplied",
			f: tftp.Flags{
				Mode: "ascii",
			},
			cmdline: []string{"localhost"},
			args:    []string{"localhost"},
			input:   []string{"q"},
		},
		{
			name: "Hostname/Port supplied",
			f: tftp.Flags{
				Mode: "ascii",
			},
			cmdline: []string{"localhost", "69"},
			args:    []string{"localhost", "69"},
			input:   []string{"q"},
		},
		{
			name: "Hostname/Port after options",
			f: tftp.Flags{
				Mode: "ascii",
			},
			cmdline: []string{"-l", "-m", "ascii", "localhost", "69"},
			args:    []string{"localhost", "69"},
			input:   []string{"q"},
		},
		{
			name: "Hostname/Port after options and literal",
			f: tftp.Flags{
				Mode: "ascii",
				Cmd:  "literal",
			},
			cmdline: []string{"-m", "ascii", "localhost", "69", "-c", "literal"},
			args:    []string{"localhost", "69"},
			input:   []string{"q"},
		},
		{
			name: "Hostname/Port after options with get",
			f: tftp.Flags{
				Mode: "ascii",
				Cmd:  "get",
			},
			cmdline: []string{"-l", "-m", "ascii", "localhost", "69", "-c", "get", "hostname:file1", "file2", "file3"},
			args:    []string{"localhost", "69", "hostname:file1", "file2", "file3"},
			input:   []string{"q"},
		},
		{
			name: "NoIPPort get",
			f: tftp.Flags{
				Mode: "ascii",
				Cmd:  "get",
			},
			cmdline: []string{"-l", "-m", "ascii", "-c", "get", "file1", "file2", "file3"},
			args:    []string{"file1", "file2", "file3"},
			input:   []string{"localhost", "q"},
		},
		{
			name: "NoIPPort put",
			f: tftp.Flags{
				Mode: "ascii",
				Cmd:  "put",
			},
			cmdline: []string{"-l", "-m", "ascii", "-c", "put", "file1", "file2", "file3"},
			args:    []string{"file1", "file2", "file3"},
			input:   []string{"localhost", "q"},
		},
		{
			name: "NoIPPort put with no args",
			f: tftp.Flags{
				Mode: "ascii",
				Cmd:  "put",
			},
			cmdline: []string{"-l", "-m", "ascii", "-c", "put"},
			args:    []string{},
			input:   []string{"localhost", "q"},
		},
		{
			name: "NoIPPort get with no args",
			f: tftp.Flags{
				Mode: "ascii",
				Cmd:  "get",
			},
			cmdline: []string{"-l", "-m", "ascii", "-c", "get"},
			args:    []string{},
			input:   []string{"localhost", "q"},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			tt := tt
			var inBuf, outBuf bytes.Buffer

			for _, in := range tt.input {
				fmt.Fprintf(&inBuf, "%s\r\n", in)
			}

			if err := run(tt.f, tt.cmdline, tt.args, &inBuf, &outBuf); !errors.Is(err, tt.err) {
				t.Errorf("run(): %v, expect: %v", err, tt.err)
			}

			if outBuf.Len() > 0 {
				if !strings.Contains(outBuf.String(), tt.exp) {
					t.Errorf("output: %s, not: %s", outBuf.String(), tt.exp)
				}
			}
		})
	}
}

func TestSplitArgs(t *testing.T) {
	for _, tt := range []struct {
		name       string
		cmdline    []string
		args       []string
		expHost    []string
		expCmdArgs []string
	}{
		{
			name:       "Hostname/Port/Get/3Files",
			cmdline:    []string{"localhost", "69", "-m", "ascii", "-c", "get", "hostname:file1", "file2", "file3"},
			args:       []string{"127.0.0.1", "69", "hostname:file1", "file2", "file3"},
			expCmdArgs: []string{"hostname:file1", "file2", "file3"},
			expHost:    []string{"127.0.0.1", "69"},
		},
		{
			name:       "NoIPPort/Get/NoFile",
			cmdline:    []string{"-m", "ascii", "-c", "get"},
			args:       []string{},
			expCmdArgs: []string{},
			expHost:    []string{},
		},
		{
			name:       "IP/Port/Get/1File",
			cmdline:    []string{"127.0.0.1", "69", "-m", "ascii", "-c", "get", "hostname:file1"},
			args:       []string{"127.0.0.1", "69", "hostname:file1"},
			expCmdArgs: []string{"hostname:file1"},
			expHost:    []string{"127.0.0.1", "69"},
		},
		{
			name:       "IP/Port/Get/2Files",
			cmdline:    []string{"127.0.0.1", "69", "-m", "ascii", "-c", "get", "hostname:file1", "file2"},
			args:       []string{"127.0.0.1", "69", "hostname:file1", "file2"},
			expCmdArgs: []string{"hostname:file1", "file2"},
			expHost:    []string{"127.0.0.1", "69"},
		},
		{
			name:       "IP/Port/Get/3Files",
			cmdline:    []string{"127.0.0.1", "69", "-m", "ascii", "-c", "get", "hostname:file1", "file2", "file3"},
			args:       []string{"127.0.0.1", "69", "hostname:file1", "file2", "file3"},
			expCmdArgs: []string{"hostname:file1", "file2", "file3"},
			expHost:    []string{"127.0.0.1", "69"},
		},
		{
			name:       "noIP/connect/with_args",
			cmdline:    []string{"-m", "ascii", "-c", "conenct", "127.0.0.1", "69"},
			args:       []string{"127.0.0.1", "69"},
			expCmdArgs: []string{"127.0.0.1", "69"},
			expHost:    []string{},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			tt := tt
			cmdArgs, host := splitArgs(tt.cmdline, tt.args)

			if !eqStringSlice(cmdArgs, tt.expCmdArgs) {
				t.Errorf("cmdsArgs: %s not equal expeted: %s", cmdArgs, tt.expCmdArgs)
			}

			if !eqStringSlice(host, tt.expHost) {
				t.Errorf("host: %s not equal expeted: %s", host, tt.expHost)
			}
		})
	}
}

func eqStringSlice(s1, s2 []string) bool {
	if len(s1) != len(s2) {
		return false
	}

	for iterator, item := range s1 {
		if item != s2[iterator] {
			return false
		}
	}

	return true
}
