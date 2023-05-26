// Copyright 2021 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// Derived work from Daniel Martí <mvdan@mvdan.cc>

package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"unicode"

	"mvdan.cc/sh/v3/interp"
	"mvdan.cc/sh/v3/syntax"
)

func TestRun(t *testing.T) {
	for _, tt := range []struct {
		name string
		args []string
	}{
		{
			name: "no args",
		},
		{
			name: "args",
			args: []string{"echo"},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			if err := run(&bytes.Buffer{}, &bytes.Buffer{}, &bytes.Buffer{}, false, tt.args...); err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
		})
	}
}

func TestRunFail(t *testing.T) {
	for _, tt := range []struct {
		name string
		args []string
	}{
		{
			name: "run a bad file",
			args: []string{"/"},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			if err := run(&bytes.Buffer{}, &bytes.Buffer{}, &bytes.Buffer{}, false, tt.args...); err == nil {
				t.Errorf("want err, got nil")
			}
		})
	}
}

func TestRunScript(t *testing.T) {
	for _, tt := range []struct {
		name  string
		pairs []string
		err   error
	}{
		{
			name: "bad file",
			pairs: []string{
				"/",
				"",
			},
			err: errors.New("open bad file: no such file or directory"),
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			runner, err := interp.New(interp.StdIO(nil, &buf, &buf))
			if err != nil {
				t.Errorf("Failed creating runner: %v", err)
			}

			parser := syntax.NewParser()

			if err := runScript(runner, parser, tt.name); err != nil {
				// can't use errors.Is: please ask mvdan to fix that.
				if fmt.Sprintf("%v", err) != fmt.Sprintf("%v", tt.err) {
					t.Errorf("got '%v', want '%v'", err, tt.err)
				}
			}

			if err := readString(&buf, tt.pairs[1]); err != nil {
				t.Errorf("Failed reading string: %v", err)
			}
		})
	}
}

func TestRunInteractive(t *testing.T) {
	for _, tt := range []struct {
		name    string
		pairs   []string
		wantErr error
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
			t.Skip("Find out how to instrument bubbline")
			inReader, inWriter := io.Pipe()
			outReader, outWriter := io.Pipe()
			runner, err := interp.New(interp.StdIO(inReader, outWriter, outWriter))
			if err != nil {
				t.Errorf("Failed creating runner: %v", err)
			}

			parser := syntax.NewParser()

			if err := runInteractive(runner, parser, outWriter, outWriter); err != nil && tt.wantErr == nil {
				t.Errorf("Unexpected error: %v", err)
			} else if tt.wantErr != nil && fmt.Sprint(err) != tt.wantErr.Error() {
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

	dirPath := f.TempDir()
	var buf bytes.Buffer
	runner, err := interp.New(interp.StdIO(nil, &buf, &buf))
	if err != nil {
		f.Fatalf("failed to initialize runner")
	}

	parser := syntax.NewParser()

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

		buf.Reset()
		runner.Reset()
		runner.Dir = dirPath

		runCmd(runner, parser, strings.NewReader(stringifiedData), "fuzz")
	})
}
