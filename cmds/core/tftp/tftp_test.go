// Copyright 2012-2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"os"
	"strings"
	"testing"
)

func TestParseRun(t *testing.T) {
	for _, tt := range []struct {
		name    string
		cmdline []string
		args    []string
		input   []string
		f       Flags
		exp     string
		err     error
	}{
		{
			name: "SimpleQuit",
			f: Flags{
				mode: "ascii",
			},
			cmdline: []string{},
			args:    []string{},
			input:   []string{"localhost", "q"},
		},
		{
			name: "HelpQuit",
			f: Flags{
				mode: "ascii",
			},
			cmdline: []string{},
			args:    []string{},
			input:   []string{"localhost", "h", "q"},
		},
		{
			name: "BinaryQuit",
			f: Flags{
				mode: "ascii",
			},
			cmdline: []string{},
			args:    []string{},
			input:   []string{"localhost", "binary", "q"},
		},
		{
			name: "AsciiQuit",
			f: Flags{
				mode: "ascii",
			},
			cmdline: []string{},
			args:    []string{},
			input:   []string{"localhost", "ascii", "q"},
		},
		{
			name: "ModeQuit",
			f: Flags{
				mode: "ascii",
			},
			cmdline: []string{},
			args:    []string{},
			input:   []string{"localhost", "mode", "q"},
			exp:     "Using netascii mode to transfer files.",
		},
		{
			name: "ModeWithValueQuit",
			f: Flags{
				mode: "ascii",
			},
			cmdline: []string{},
			args:    []string{},
			input:   []string{"localhost", "mode binary", "q"},
			exp:     "Using octet mode to transfer files.",
		},
		{
			name: "ModeWithValueQuit",
			f: Flags{
				mode: "binary",
			},
			cmdline: []string{},
			args:    []string{},
			input:   []string{"localhost", "mode ascii", "q"},
			exp:     "Using netascii mode to transfer files.",
		},
		{
			name: "ModeWithErrorQuit",
			f: Flags{
				mode: "binary",
			},
			cmdline: []string{},
			args:    []string{},
			input:   []string{"localhost", "mode error", "q"},
			exp:     fmt.Sprintf("%v", errInvalidTransferMode),
		},
		{
			name: "LiteralQuit",
			f: Flags{
				mode: "ascii",
			},
			cmdline: []string{},
			args:    []string{},
			input:   []string{"localhost", "literal", "q"},
			exp:     "Literal mode is on",
		},
		{
			name: "rexmtQuit",
			f: Flags{
				mode: "ascii",
			},
			cmdline: []string{},
			args:    []string{},
			input:   []string{"localhost", "rexmt 20", "q"},
		},
		{
			name: "timeoutQuit",
			f: Flags{
				mode: "ascii",
			},
			cmdline: []string{},
			args:    []string{},
			input:   []string{"localhost", "timeout 20", "q"},
		},
		{
			name: "traceQuit",
			f: Flags{
				mode: "ascii",
			},
			cmdline: []string{},
			args:    []string{},
			input:   []string{"localhost", "trace", "q"},
			exp:     "Packet tracing on.",
		},
		{
			name: "verboseQuit",
			f: Flags{
				mode: "ascii",
			},
			cmdline: []string{},
			args:    []string{},
			input:   []string{"localhost", "verbose", "q"},
			exp:     "Verbose mode on.",
		},
		{
			name: "connectQuit",
			f: Flags{
				mode: "ascii",
			},
			cmdline: []string{},
			args:    []string{},
			input:   []string{"localhost", "connect localhost 69", "q"},
		},
		{
			name: "statusQuit",
			f: Flags{
				mode: "ascii",
			},
			cmdline: []string{},
			args:    []string{},
			input:   []string{"localhost", "status", "q"},
			exp:     "Connected to localhost\nMode: netascii Verbose: off Tracing: off Literal: of",
		},
		{
			name: "IP supplied",
			f: Flags{
				mode: "ascii",
			},
			cmdline: []string{"127.0.0.1"},
			args:    []string{"127.0.0.1"},
			input:   []string{"q"},
		},
		{
			name: "IP/Port supplied",
			f: Flags{
				mode: "ascii",
			},
			cmdline: []string{"127.0.0.1", "69"},
			args:    []string{"127.0.0.1", "69"},
			input:   []string{"q"},
		},
		{
			name: "Hostname supplied",
			f: Flags{
				mode: "ascii",
			},
			cmdline: []string{"localhost"},
			args:    []string{"localhost"},
			input:   []string{"q"},
		},
		{
			name: "Hostname/Port supplied",
			f: Flags{
				mode: "ascii",
			},
			cmdline: []string{"localhost", "69"},
			args:    []string{"localhost", "69"},
			input:   []string{"q"},
		},
		{
			name: "Hostname/Port after options",
			f: Flags{
				mode: "ascii",
			},
			cmdline: []string{"-l", "-m", "ascii", "localhost", "69"},
			args:    []string{"localhost", "69"},
			input:   []string{"q"},
		},
		{
			name: "Hostname/Port after options and literal",
			f: Flags{
				mode: "ascii",
				cmd:  "literal",
			},
			cmdline: []string{"-m", "ascii", "localhost", "69", "-c", "literal"},
			args:    []string{"localhost", "69"},
			input:   []string{"q"},
		},
		/*
			{
				name: "Hostname/Port after options with get",
				f: Flags{
					mode: "ascii",
					cmd:  "get",
				},
				cmdline: []string{"-l", "-m", "ascii", "localhost", "69", "-c", "get", "hostname:file1", "file2", "file3"},
				args:    []string{"localhost", "69", "hostname:file1", "file2", "file3"},
			},
			{
				name: "NoIPPort get with no args",
				f: Flags{
					mode: "ascii",
					cmd:  "get",
				},
				cmdline: []string{"-l", "-m", "ascii", "-c", "get", "file1", "file2", "file3"},
				args:    []string{"file1", "file2", "file3"},
			},
			{
				name: "NoIPPort put with no args",
				f: Flags{
					mode: "ascii",
					cmd:  "put",
				},
				cmdline: []string{"-l", "-m", "ascii", "-c", "put", "file1", "file2", "file3"},
				args:    []string{"file1", "file2", "file3"},
			},
			{
				name: "NoIPPort put with no args",
				f: Flags{
					mode: "ascii",
					cmd:  "put",
				},
				cmdline: []string{"-l", "-m", "ascii", "localhost", "69", "-c", "put"},
				args:    []string{"localhost", "69"},
			},
		*/
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

func TestReadInteractiveInput(t *testing.T) {
	for _, tt := range []struct {
		name  string
		input []string
		exp   []string
	}{
		{
			name:  "Hello_World",
			input: []string{"Hello World"},
			exp:   []string{"Hello", "World"},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			tt := tt
			var inbuf, outBuf bytes.Buffer

			inScan := bufio.NewScanner(&inbuf)

			for _, input := range tt.input {
				fmt.Fprintf(&inbuf, "%s\n", input)
				ret := readInputInteractive(inScan, &outBuf)
				for iterator, r := range ret {
					if r != tt.exp[iterator] {
						t.Errorf("%s != %s", ret[iterator], input)
					}
				}

			}

		})
	}
}

func TestValidateMode(t *testing.T) {
	for _, tt := range []struct {
		input string
		err   error
	}{
		{
			input: "ascii",
		},
		{
			input: "binary",
		},
		{
			input: "garbage",
			err:   errInvalidTransferMode,
		},
	} {
		t.Run(tt.input, func(t *testing.T) {
			tt := tt
			if _, err := validateMode(tt.input); !errors.Is(err, tt.err) {
				t.Errorf("validateMode(): %v, not: %v", err, tt.err)
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

func TestConstructURL(t *testing.T) {
	for _, tt := range []struct {
		name string
		host string
		port string
		dir  string
		file string
		exp  string
	}{
		{
			name: "SimpleHostPortFile",
			host: "localhost",
			port: "69",
			dir:  "",
			file: "abc.file",
			exp:  "tftp://localhost:69/abc.file",
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			url := constructURL(tt.host, tt.port, tt.dir, tt.file)
			if url != tt.exp {
				t.Errorf("constructURL() failed:%s, want: %s", url, tt.exp)
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

func TestExecuteGetPut(t *testing.T) {
	for _, tt := range []struct {
		name   string
		client ClientIf
		host   string
		port   string
		dir    string
		files  []string
	}{
		{
			name:   "GetPut1File",
			client: &ClientMock{},
			host:   "localhost",
			port:   "69",
			dir:    "",
			files:  []string{"abc.file"},
		},
		{
			name:   "GetPut2File",
			client: &ClientMock{},
			host:   "localhost",
			port:   "69",
			dir:    "",
			files:  []string{"abc.file", "cde.file"},
		},
		{
			name:   "GetPut3File",
			client: &ClientMock{},
			host:   "localhost",
			port:   "69",
			dir:    "",
			files:  []string{"abc.file", "cde.file", "fgh.file"},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			files := make([]string, 0)
			for _, file := range tt.files {
				tf, err := os.CreateTemp("", file)
				if err != nil {
					t.Error(err)
				}
				files = append(files, tf.Name())
				tf.Close()
			}
			if err := executeGet(tt.client, tt.host, tt.port, files); err != nil {
				t.Error(err)
			}

			for _, file := range files {
				if err := os.Remove(file); err != nil {
					t.Error(err)
				}
			}
		})

		t.Run(tt.name, func(t *testing.T) {
			files := make([]string, 0)
			for _, file := range tt.files {
				tf, err := os.CreateTemp("", file)
				if err != nil {
					t.Error(err)
				}
				files = append(files, tf.Name())
				tf.Close()
			}

			if err := executePut(tt.client, tt.host, tt.port, files); err != nil {
				t.Error(err)
			}

			for _, file := range files {
				if err := os.Remove(file); err != nil {
					t.Error(err)
				}
			}
		})

	}

}
