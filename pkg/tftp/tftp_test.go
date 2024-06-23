// Copyright 2012-2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tftp

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"os"
	"strings"
	"testing"
)

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
			err:   ErrInvalidTransferMode,
		},
	} {
		t.Run(tt.input, func(t *testing.T) {
			tt := tt
			if _, err := ValidateMode(tt.input); !errors.Is(err, tt.err) {
				t.Errorf("validateMode(): %v, not: %v", err, tt.err)
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

func TestRunInteractive(t *testing.T) {
	for _, tt := range []struct {
		name   string
		f      Flags
		ipPort []string
		input  []string
		expErr error
		expOut string
	}{
		{
			name: "Input_Localhost_Quit",
			input: []string{
				"localhost",
				"q",
			},
		},
		{
			name:   "Localhost_Quit",
			ipPort: []string{"localhost"},
			input:  []string{"q"},
		},
		{
			name:   "Localhost_Quit",
			ipPort: []string{"localhost", "69"},
			input: []string{
				"q",
			},
		},
		{
			name: "Input_Help_Quit",
			input: []string{
				"localhost",
				"help",
				"q",
			},
		},
		{
			name: "Input_ascii_Quit",
			input: []string{
				"localhost",
				"ascii",
				"q",
			},
		},
		{
			name: "Input_binary_Quit",
			input: []string{
				"localhost",
				"binary",
				"q",
			},
		},
		{
			name: "Input_ascii_Quit",
			f: Flags{
				Mode: "ascii",
			},
			input: []string{
				"localhost",
				"mode",
				"q",
			},
			expOut: "Using netascii mode to transfer files.",
		},
		{
			name: "Input_mode-ascii_Quit",
			input: []string{
				"localhost",
				"mode ascii",
				"q",
			},
		},
		{
			name: "Input_mode-binary_Quit",
			input: []string{
				"localhost",
				"mode binary",
				"q",
			},
		},
		{
			name: "Input_mode-error_Quit",
			input: []string{
				"localhost",
				"mode error",
				"q",
			},
			expOut: fmt.Sprintf("%v", ErrInvalidTransferMode),
		},
		{
			name: "Input_mode-binary_Quit",
			input: []string{
				"localhost",
				"connect localhost",
				"q",
			},
		},
		{
			name: "Input_mode-binary_Quit",
			input: []string{
				"localhost",
				"connect localhost 99",
				"q",
			},
		},
		{
			name: "Input_literal_Quit",
			input: []string{
				"localhost",
				"literal",
				"q",
			},
		},
		{
			name: "Input_rexmt_Quit",
			input: []string{
				"localhost",
				"rexmt 10",
				"q",
			},
		},
		{
			name: "Input_status_Quit",
			input: []string{
				"localhost",
				"status",
				"q",
			},
		},
		{
			name: "Input_timeout_Quit",
			input: []string{
				"localhost",
				"timeout 10",
				"q",
			},
		},
		{
			name: "Input_trace_Quit",
			input: []string{
				"localhost",
				"trace",
				"q",
			},
		},
		{
			name: "Input_verbose_Quit",
			input: []string{
				"localhost",
				"verbose",
				"q",
			},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			var inBuf, outBuf bytes.Buffer

			for _, in := range tt.input {
				fmt.Fprintf(&inBuf, "%s\n", in)
			}

			if err := RunInteractive(tt.f, tt.ipPort, &inBuf, &outBuf); !errors.Is(err, tt.expErr) {
				t.Errorf("RunInteractive(): %v, not %v", err, tt.expErr)
			}

			out := outBuf.String()
			if out != "" {
				if !strings.Contains(out, tt.expOut) {
					t.Errorf("unexpected output: %s, not containing: %s", out, tt.expOut)
				}
			}
		})
	}
}

func TestReadHostInteractive(t *testing.T) {
	for _, tt := range []struct {
		input  string
		output string
	}{
		{
			input:  "localhost",
			output: "localhost",
		},
	} {
		t.Run(tt.input, func(t *testing.T) {
			var inBuf, outBuf bytes.Buffer
			fmt.Fprintf(&inBuf, "%s\n", tt.input)
			inScan := bufio.NewScanner(&inBuf)
			ret := readHostInteractive(inScan, &outBuf)
			if ret != tt.output {
				t.Errorf("readHostInteractive(): %s, not %s", ret, tt.output)
			}
		})
	}
}

func TestStatusString(t *testing.T) {
	for _, tt := range []struct {
		mode bool
		exp  string
	}{
		{
			mode: true,
			exp:  "on",
		},
		{
			mode: false,
			exp:  "off",
		},
	} {
		t.Run(tt.exp, func(t *testing.T) {
			ret := statusString(tt.mode)

			if ret != tt.exp {
				t.Errorf("statusString(): %s, not: %s", ret, tt.exp)
			}
		})
	}
}
