// Copyright 2021 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// Derived work from Daniel Martí <mvdan@mvdan.cc>
//go:build !tinygo || tinygo.enable

package main

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"syscall"
	"testing"
	"time"
	"unicode"

	"github.com/Netflix/go-expect"
	"github.com/u-root/gobusybox/src/pkg/golang"
	"mvdan.cc/sh/v3/interp"
	"mvdan.cc/sh/v3/syntax"
)

func TestRun(t *testing.T) {
	echoSc := filepath.Join(t.TempDir(), "b.sh")
	_ = os.WriteFile(echoSc, []byte("echo foo"), 0o777)

	for _, tt := range []struct {
		name    string
		command string
		args    []string
		wantOut string
		wantErr string
	}{
		{
			name: "no args",
		},
		{
			name:    "args",
			args:    []string{echoSc},
			wantOut: "foo\n",
		},
		{
			name:    "cmd",
			command: "echo foo",
			wantOut: "foo\n",
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			// gosh, since 63f3119ec, can no longer safely use a bytes.Buffer
			// for stdin. There may be a better fix, but for now we create an
			// io.pipe and close the write side, passing the read side to run.
			inr, inw := io.Pipe()
			inw.Close()
			var out, err bytes.Buffer
			if err := run(inr, &out, &err, tt.command, tt.args...); err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
			t.Logf("out: %s", out.Bytes())
			t.Logf("err: %s", err.Bytes())
			if gotOut := out.String(); gotOut != tt.wantOut {
				t.Errorf("Stdout = %s, want %s", gotOut, tt.wantOut)
			}
			if gotErr := err.String(); gotErr != tt.wantErr {
				t.Errorf("Stderr = %s, want %s", gotErr, tt.wantErr)
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
			// gosh, since 63f3119ec, can no longer safely use a bytes.Buffer
			// for stdin. There may be a better fix, but for now we create an
			// io.pipe and close the write side, passing the read side to run.
			inr, inw := io.Pipe()
			inw.Close()
			if err := run(inr, &bytes.Buffer{}, &bytes.Buffer{}, "", tt.args...); err == nil {
				t.Errorf("want err, got nil")
			}
		})
	}
}

func TestRunScript(t *testing.T) {
	d := t.TempDir()
	script := filepath.Join(d, "a.sh")
	if err := os.WriteFile(script, []byte("echo hi\n"), 0o666); err != nil {
		t.Fatalf("Writing %q: got %v, want nil", script, err)
	}

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
			err: errors.New("read /: is a directory"),
		},
		{
			name: "bad file",
			pairs: []string{
				"bad file",
				"",
			},
			err: errors.New("open bad file: no such file or directory"),
		},
		{
			name: "echo script",
			pairs: []string{
				script,
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

			if err := runScript(runner, tt.pairs[0]); err != nil {
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
			t.Skip("Currently broken")
			inReader, inWriter := io.Pipe()
			outReader, outWriter := io.Pipe()
			runner, err := interp.New(interp.StdIO(inReader, outWriter, outWriter))
			if err != nil {
				t.Errorf("Failed creating runner: %v", err)
			}

			if err := runInteractive(runner, syntax.NewParser(), outWriter, outWriter); err != nil && tt.wantErr == nil {
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
		"wait with args not handled yet",
	}
	re := strings.NewReplacer("\x22", "", "\x24", "", "\x26", "", "\x27", "", "\x28", "", "\x29", "", "\x2A", "", "\x3C", "", "\x3E", "", "\x3F", "", "\x5C", "", "\x7C", "")

	dirPath := f.TempDir()
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

		buf.Reset()
		runner.Reset()
		runner.Dir = dirPath

		runReader(runner, strings.NewReader(stringifiedData), "fuzz")
	})
}

type consoleAction func(*expect.Console) error

func expectOpts(o ...expect.ExpectOpt) consoleAction {
	return func(c *expect.Console) error {
		_, err := c.Expect(o...)
		return err
	}
}

func expectString(s string) consoleAction {
	return func(c *expect.Console) error {
		_, err := c.ExpectString(s)
		if err != nil {
			return fmt.Errorf("failed to expect %q: %w", s, err)
		}
		return nil
	}
}

func send(s string) consoleAction {
	return func(c *expect.Console) error {
		_, err := c.Send(s)
		return err
	}
}

func TestInteractiveBubbline(t *testing.T) {
	dir := t.TempDir()
	execPath := filepath.Join(dir, "gosh")

	var opts *golang.BuildOpts
	// Setting -cover without GOCOVERDIR adds extra warning output, which changes the result of the test.
	if os.Getenv("GOCOVERDIR") != "" {
		opts = &golang.BuildOpts{ExtraArgs: []string{"-covermode=atomic"}}
	}
	// Build the stuff.
	if err := golang.Default(golang.DisableCGO()).BuildDir("", execPath, opts); err != nil {
		t.Fatal(err)
	}

	for _, tt := range []struct {
		name   string
		expect []consoleAction
	}{
		{
			name: "exit shell",
			expect: []consoleAction{
				expectString("> "),
				send("echo hi\x0D"),
				expectString("> "),
				send("exit\x0D"),
			},
		},
		{
			name: "source script",
			expect: []consoleAction{
				expectString("> "),
				send("source ./testdata/setenv.sh && echo $FOO\x0D"),
				expectString("hi"),
				expectString("hahaha"),
				expectString("> "),
				send("exit\x0D"),
			},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			con, err := expect.NewTestConsole(t, expect.WithDefaultTimeout(2*time.Second), expect.WithStdout(os.Stdout))
			if err != nil {
				t.Fatal(err)
			}

			cmd := exec.CommandContext(context.Background(), execPath)
			cmd.Stdin, cmd.Stdout, cmd.Stderr = con.Tty(), con.Tty(), con.Tty()
			if err := cmd.Start(); err != nil {
				t.Fatal(err)
			}
			// Close our end of child's tty.
			con.Tty().Close()

			for i, a := range tt.expect {
				if err := a(con); err != nil {
					t.Errorf("Action %d: %v", i, err)
				}
			}

			if err := cmd.Wait(); err != nil {
				t.Error(err)
			}
			if err := con.Close(); err != nil {
				t.Error(err)
			}
		})
	}
}

func TestInteractiveLiner(t *testing.T) {
	dir := t.TempDir()
	execPath := filepath.Join(dir, "gosh")

	var opts *golang.BuildOpts
	// Setting -cover without GOCOVERDIR adds extra warning output, which changes the result of the test.
	if os.Getenv("GOCOVERDIR") != "" {
		opts = &golang.BuildOpts{ExtraArgs: []string{"-covermode=atomic"}}
	}
	// Build the stuff.
	if err := golang.Default(golang.DisableCGO(), golang.WithBuildTag("goshliner")).BuildDir("", execPath, opts); err != nil {
		t.Fatal(err)
	}

	for _, tt := range []struct {
		name   string
		expect []consoleAction
	}{
		{
			name: "exit shell",
			expect: []consoleAction{
				expectString("$ "),
				send("echo hi\x0D"),
				expectString("hi"),
				expectString("$ "),
				send("exit\x0D"),
			},
		},
		{
			name: "source script",
			expect: []consoleAction{
				expectString("$ "),
				send("source ./testdata/setenv.sh && echo $FOO\x0D"),
				expectString("hi"),
				expectString("hahaha"),
				expectString("$ "),
				send("exit\x0D"),
			},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			con, err := expect.NewTestConsole(t, expect.WithDefaultTimeout(2*time.Second))
			if err != nil {
				t.Fatal(err)
			}

			cmd := exec.CommandContext(context.Background(), execPath)
			cmd.Stdin, cmd.Stdout, cmd.Stderr = con.Tty(), con.Tty(), con.Tty()
			if err := cmd.Start(); err != nil {
				t.Fatal(err)
			}
			// Close our end of child's tty.
			con.Tty().Close()

			for i, a := range tt.expect {
				if err := a(con); err != nil {
					t.Errorf("Action %d: %v", i, err)
				}
			}

			if err := cmd.Wait(); err != nil {
				t.Error(err)
			}
			if err := con.Close(); err != nil {
				t.Error(err)
			}
		})
	}
}

func TestGoshInvocation(t *testing.T) {
	dir := t.TempDir()
	execPath := filepath.Join(dir, "gosh")
	var opts *golang.BuildOpts
	// Setting -cover without GOCOVERDIR adds extra warning output, which changes the result of the test.
	if os.Getenv("GOCOVERDIR") != "" {
		opts = &golang.BuildOpts{ExtraArgs: []string{"-covermode=atomic"}}
	}
	// Build the stuff.
	if err := golang.Default(golang.DisableCGO()).BuildDir("", execPath, opts); err != nil {
		t.Fatal(err)
	}

	for _, tt := range []struct {
		name       string
		args       []string
		stdout     string
		stdin      string
		exitStatus int
	}{
		{
			name: "echo",
			args: []string{
				"-c", "echo foo",
			},
			stdout: "foo\n",
		},
		{
			name: "cmdline cmd wins",
			args: []string{
				"-c", "echo foo", "./testdata/setenv.sh",
			},
			stdout: "foo\n",
		},
		{
			name: "cmd",
			args: []string{
				"./testdata/setenv.sh",
			},
			stdout: "hi\n",
		},
		{
			name:   "stdin",
			stdin:  "echo hi",
			stdout: "hi\n",
		},
		{
			name:       "exit status",
			stdin:      "false",
			exitStatus: 1,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			cmd := exec.CommandContext(context.Background(), execPath, tt.args...)
			var b bytes.Buffer
			if tt.stdin == "" {
				// This means os.Stdin would register as a terminal.
				con, err := expect.NewTestConsole(t, expect.WithDefaultTimeout(2*time.Second))
				if err != nil {
					t.Fatal(err)
				}
				cmd.Stdin = con.Tty()
				t.Cleanup(func() { con.Close() })
			} else {
				cmd.Stdin = strings.NewReader(tt.stdin)
			}
			cmd.Stdout, cmd.Stderr = &b, &b
			if err := cmd.Start(); err != nil {
				t.Fatal(err)
			}

			var e *exec.ExitError
			if err := cmd.Wait(); errors.As(err, &e) {
				if got := e.Sys().(syscall.WaitStatus).ExitStatus(); got != tt.exitStatus {
					t.Errorf("Exit status = %v, want %d", err, tt.exitStatus)
				}
			} else if err != nil {
				t.Error(err)
			}
			if got := b.String(); got != tt.stdout {
				t.Errorf("gosh = %s, want %s", got, tt.stdout)
			}
		})
	}
}

func TestInteractiveSimple(t *testing.T) {
	dir := t.TempDir()
	execPath := filepath.Join(dir, "gosh")

	var opts *golang.BuildOpts
	// Setting -cover without GOCOVERDIR adds extra warning output, which changes the result of the test.
	if os.Getenv("GOCOVERDIR") != "" {
		opts = &golang.BuildOpts{ExtraArgs: []string{"-covermode=atomic"}}
	}
	// Build the stuff.
	if err := golang.Default(golang.DisableCGO(), golang.WithBuildTag("goshsmall")).BuildDir("", execPath, opts); err != nil {
		t.Fatal(err)
	}

	for _, tt := range []struct {
		name   string
		expect []consoleAction
	}{
		{
			name: "exit shell",
			expect: []consoleAction{
				expectString("$ "),
				send("echo hi\x0D"),
				expectString("hi"),
				expectString("$ "),
				send("exit\x0D"),
			},
		},
		{
			name: "source script",
			expect: []consoleAction{
				expectString("$ "),
				send("source ./testdata/setenv.sh && echo $FOO\x0D"),
				expectString("hi"),
				expectString("hahaha"),
				expectString("$ "),
				send("exit\x0D"),
			},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			con, err := expect.NewTestConsole(t, expect.WithDefaultTimeout(2*time.Second))
			if err != nil {
				t.Fatal(err)
			}

			cmd := exec.CommandContext(context.Background(), execPath)
			cmd.Stdin, cmd.Stdout, cmd.Stderr = con.Tty(), con.Tty(), con.Tty()
			if err := cmd.Start(); err != nil {
				t.Fatal(err)
			}
			// Close our end of child's tty.
			con.Tty().Close()

			for i, a := range tt.expect {
				if err := a(con); err != nil {
					t.Errorf("Action %d: %v", i, err)
				}
			}

			if err := cmd.Wait(); err != nil {
				t.Error(err)
			}
			if err := con.Close(); err != nil {
				t.Error(err)
			}
		})
	}
}
