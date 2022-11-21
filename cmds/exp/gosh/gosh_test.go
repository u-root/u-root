// Copyright 2021 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// Derived work from Daniel Martí <mvdan@mvdan.cc>

package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"unicode"

	"github.com/u-root/prompt"
	"mvdan.cc/sh/v3/interp"
)

func TestRun(t *testing.T) {
	for _, tt := range []struct {
		name    string
		pairs   []string
		wantErr string
	}{
		{
			name: "echo foo",
			pairs: []string{
				"echo foo",
				"foo",
			},
		},
		{
			name: "quoted echo",
			pairs: []string{
				"echo 'foo\nbar'",
				"foo\nbar",
			},
		},
		{
			name: "exit 1",
			pairs: []string{
				"exit 1; echo foo",
				"",
			},
			wantErr: "exit status 1",
		},
		{
			name: "not parsable",
			pairs: []string{
				"(",
				"",
			},
			wantErr: "not parsable:1:1: reached EOF without matching ( with )",
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			sh := shell{}
			var buf bytes.Buffer
			runner, err := interp.New(interp.StdIO(strings.NewReader(tt.pairs[0]), &buf, &buf))
			if err != nil {
				t.Errorf("Failed creating runner: %v", err)
			}

			if err := sh.run(runner, strings.NewReader(tt.pairs[0]), tt.name); err != nil {
				if err.Error() != tt.wantErr {
					t.Errorf("Failed running command: %v", err)
				}
			}

			if err := readString(&buf, tt.pairs[1]); err != nil {
				t.Errorf("Failed reading string: %v", err)
			}
		})
	}
}

func TestRunAll(t *testing.T) {
	for _, tt := range []struct {
		name string
		narg int
	}{
		{
			name: "no args",
			narg: 0,
		},
		{
			name: "args",
			narg: 1,
		},
		{
			name: "negative args",
			narg: -1,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			sh := shell{}
			if err := sh.runAll(tt.narg); err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
		})
	}
}

func TestRunInteractive(t *testing.T) {
	for _, tt := range []struct {
		name    string
		pairs   []string
		wantErr string
	}{
		{},
		{
			name: "newlines",
			pairs: []string{
				"\n",
				"$ ",
				"\n",
				"$ ",
			},
		},
		{
			name: "echo foo",
			pairs: []string{
				"echo foo\n",
				"foo\n",
			},
		},
		{
			name: "echo foo bar",
			pairs: []string{
				"echo foo\n",
				"foo\n$ ",
				"echo bar\n",
				"bar\n",
			},
		},
		{
			name: "if then",
			pairs: []string{
				"if true\n",
				"> ",
				"then echo bar; fi\n",
				"bar\n",
			},
		},
		{
			name: "quoted echo",
			pairs: []string{
				"echo 'foo\n",
				"> ",
				"bar'\n",
				"foo\nbar\n",
			},
		},
		{
			name: "echo with semicolon",
			pairs: []string{
				"echo foo; echo bar\n",
				"foo\nbar\n",
			},
		},
		{
			name: "echo with semicolon quoted",
			pairs: []string{
				"echo foo; echo 'bar\n",
				"> ",
				"baz'\n",
				"foo\nbar\nbaz\n",
			},
		},
		{
			name: "braces",
			pairs: []string{
				"(\n",
				"> ",
				"echo foo)\n",
				"foo\n",
			},
		},
		{
			name: "double brackets",
			pairs: []string{
				"[[\n",
				"> ",
				"true ]]\n",
				"$ ",
			},
		},
		{
			name: "logic or",
			pairs: []string{
				"echo foo ||\n",
				"> ",
				"echo bar\n",
				"foo\n",
			},
		},
		{
			name: "pipe",
			pairs: []string{
				"echo foo |\n",
				"> ",
				"read var; echo $var\n",
				"foo\n",
			},
		},
		{
			name: "delayed newline",
			pairs: []string{
				"echo foo",
				"",
				" bar\n",
				"foo bar\n",
			},
		},
		{
			name: "escaped newline",
			pairs: []string{
				"echo\\\n",
				"> ",
				" foo\n",
				"foo\n",
			},
		},
		{
			name: "echo foo with escaped newline",
			pairs: []string{
				"echo foo\\\n",
				"> ",
				"bar\n",
				"foobar\n",
			},
		},
		{
			name: "utf8",
			pairs: []string{
				"echo 你好\n",
				"你好\n$ ",
			},
		},
		{
			name: "exit 0",
			pairs: []string{
				"echo foo; exit 0; echo bar\n",
				"foo\n",
				"echo baz\n",
				"",
			},
		},
		{
			name: "exit 1",
			pairs: []string{
				"echo foo; exit 1; echo bar\n",
				"foo\n",
				"echo baz\n",
				"",
			},
			wantErr: "exit status 1",
		},
		{
			name: "no closing brace",
			pairs: []string{
				"(\n",
				"> ",
			},
			wantErr: "1:1: reached EOF without matching ( with )",
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			sh := shell{}
			inReader, inWriter := io.Pipe()
			outReader, outWriter := io.Pipe()
			runner, err := interp.New(interp.StdIO(inReader, outWriter, outWriter))
			if err != nil {
				t.Errorf("Failed creating runner: %v", err)
			}
			errc := make(chan error, 1)
			go func() {
				errc <- sh.runInteractive(runner, inReader, outWriter)
				if _, err := io.Copy(io.Discard, inReader); err != nil {
					t.Errorf("Error discarding IO: %v", err)
				}
			}()

			if err := readString(outReader, "$ "); err != nil {
				t.Errorf("Failed reading string: %v", err)
			}

			for len(tt.pairs) > 0 {
				if _, err := io.WriteString(inWriter, tt.pairs[0]); err != nil {
					t.Errorf("Failed writing string: %v", err)
				}
				if err := readString(outReader, tt.pairs[1]); err != nil {
					t.Errorf("Failed reading string: %v", err)
				}

				tt.pairs = tt.pairs[2:]
			}

			// Close the input pipe, so that the parser can stop
			if err := inWriter.Close(); err != nil {
				t.Errorf("Failed closing input pipe: %v", err)
			}

			// Once the input pipe is closed, close the output pipe
			// so that any remaining prompt writes get discarded
			if err := outReader.Close(); err != nil {
				t.Errorf("Failed closing output pipe: %v", err)
			}

			err = <-errc
			if err != nil && tt.wantErr == "" {
				t.Errorf("Unexpected error: %v", err)
			} else if tt.wantErr != "" && fmt.Sprint(err) != tt.wantErr {
				t.Errorf("Want error %q, got: %v", tt.wantErr, err)
			}
		})
	}
}

func TestRunInteractiveTabCompletion(t *testing.T) {
	for _, tt := range []struct {
		name    string
		pairs   []string
		wantErr string
	}{
		{
			name: "exit shell",
			pairs: []string{
				"exit",
				"",
			},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			sh := shell{
				input: testInputPrompt{
					inputText: tt.pairs[0],
				},
			}
			inReader, inWriter := io.Pipe()
			outReader, outWriter := io.Pipe()
			runner, err := interp.New(interp.StdIO(inReader, outWriter, outWriter))
			if err != nil {
				t.Errorf("Failed creating runner: %v", err)
			}

			if err := sh.runInteractiveTabCompletion(runner, outWriter); err != nil && tt.wantErr == "" {
				t.Errorf("Unexpected error: %v", err)
			} else if tt.wantErr != "" && fmt.Sprint(err) != tt.wantErr {
				t.Errorf("Want error %q, got: %v", tt.wantErr, err)
			}

			if err := readString(outReader, tt.pairs[1]); err != nil {
				t.Errorf("Failed reading string: %v", err)
			}

			// Close the input pipe, so that the parser can stop
			if err := inWriter.Close(); err != nil {
				t.Errorf("Failed closing input pipe: %v", err)
			}

			// Once the input pipe is closed, close the output pipe
			// so that any remaining prompt writes get discarded
			if err := outReader.Close(); err != nil {
				t.Errorf("Failed closing output pipe: %v", err)
			}
		})
	}
}

func readString(r io.Reader, want string) error {
	p := make([]byte, len(want))
	_, err := io.ReadFull(r, p)
	if err != nil {
		return err
	}
	got := string(p)
	if got != want {
		return fmt.Errorf("readString: read %q, wanted %q", got, want)
	}
	return nil
}

type testInputPrompt struct {
	inputText string
}

func (i testInputPrompt) Input(prefix string, completer prompt.Completer, opts ...prompt.Option) string {
	return i.inputText
}

// This test is only intended to run in a sandboxed environment.
// DO NOT run these fuzzing tests on your local system. Executing random commands might mess with your system.
// The fuzzing test might panic, which is checked against a defined array of expected panic messages. Hence only unexpected panics will fail the test.
// Additionally the input space is being stripped from special characters that might hinder terminating the shell in time.
func FuzzRun(f *testing.F) {
	expectedPanics := []string{
		"interface conversion",
		"param expansion",
		"regexp: Compile",
		"runtime error",
		"unexpected arithm expr",
		"unhandled builtin",
		"unhandled command node",
		"unhandled conversion of kind",
		"unhandled redirect op",
		"unhandled shopt flag",
		"unhandled unary test op",
		"unhandled word part",
		"variable name must not be empty",
		"wait with args not handled yet"}
	re := strings.NewReplacer("\x22", "", "\x24", "", "\x26", "", "\x27", "", "\x28", "", "\x29", "", "\x2A", "", "\x3C", "", "\x3E", "", "\x3F", "", "\x5C", "", "\x7C", "")

	sh := shell{}
	var buf bytes.Buffer
	runner, err := interp.New(interp.StdIO(nil, &buf, &buf))
	if err != nil {
		f.Fatalf("failed to initialize runner")
	}

	// get seed corpora
	seeds, err := filepath.Glob("testdata/fuzz/corpora/*.seed")
	if err != nil {
		f.Fatalf("failed to find seed corpora files: %v", err)
	}

	for _, seed := range seeds {
		seedBytes, err := os.ReadFile(seed)
		if err != nil {
			f.Fatalf("failed to read seed corpora from file %v: %v", seed, err)
		}

		f.Add(seedBytes)
	}

	f.Fuzz(func(t *testing.T, data []byte) {
		defer func() {
			if err := recover(); err != nil {
				for _, expPanic := range expectedPanics {
					switch err := err.(type) {
					case string:
						if strings.Contains(err, expPanic) {
							return
						}
					case runtime.Error:
						if strings.Contains(err.Error(), expPanic) {
							return
						}
					case error:
						if strings.Contains(err.Error(), expPanic) {
							return
						}
					}
				}
				t.Fatalf("Unexpected panic: %v", err)
			}
		}()

		if len(data) > 32 {
			return
		}

		// reduce the input space to a set of printable ASCII chars excluding some special characters
		for _, v := range data {
			if v < 0x20 || v > unicode.MaxASCII {
				return
			}
		}

		stringifiedData := re.Replace(string(data))
		if stringifiedData != string(data) {
			return
		}

		if strings.Contains(stringifiedData, "fuzz") {
			return
		}

		dirPath := t.TempDir()
		buf.Reset()
		runner.Reset()
		runner.Dir = dirPath

		sh.run(runner, strings.NewReader(stringifiedData), "fuzz")
	})
}
