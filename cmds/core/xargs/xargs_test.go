// Copyright 2013-2023 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestCommandNotFound(t *testing.T) {
	stdin := strings.NewReader("hello world")
	p := params{maxArgs: 1, trace: false}
	c := command(stdin, nil, nil, p)
	err := c.run("commandnotfound", "arg1")
	if !errors.Is(err, exec.ErrNotFound) {
		t.Errorf("expected %v, got %v", exec.ErrNotFound, err)
	}
}

func TestEcho(t *testing.T) {
	stdin := strings.NewReader("hello world")
	stdout := &bytes.Buffer{}
	c := command(stdin, stdout, nil, params{maxArgs: defaultMaxArgs})
	err := c.run()
	if err != nil {
		t.Fatalf("expected nil, got %v", err)
	}

	if stdout.String() != "hello world\n" {
		t.Errorf("expected 'hello world', got %q", stdout.String())
	}
}

func TestEchoWithMaxArgs(t *testing.T) {
	stdin := strings.NewReader("a b c d e f g")
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	c := command(stdin, stdout, stderr, params{maxArgs: 3, trace: true})
	err := c.run()
	if err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
	if stdout.String() != "a b c\nd e f\ng\n" {
		t.Errorf("expected 'a b c\nd e f\ng\n', got %q", stdout.String())
	}
	expectedStderr := "echo a b c\necho d e f\necho g\n"
	if stderr.String() != expectedStderr {
		t.Errorf("expected %q, got %q", expectedStderr, stderr.String())
	}
}

func TestEchoPromt(t *testing.T) {
	stdin := strings.NewReader("a b c")
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}

	dir := t.TempDir()
	path := filepath.Join(dir, "tty")
	err := os.WriteFile(path, []byte("yes\nn\ny\n"), 0644)
	if err != nil {
		t.Fatalf("expected nil, got %v", err)
	}

	p := params{maxArgs: 1, prompt: true, trace: true}
	c := command(stdin, stdout, stderr, p)
	c.tty = path
	err = c.run()
	if err != nil {
		t.Fatalf("expected nil, got %v", err)
	}

	if stdout.String() != "a\nc\n" {
		t.Errorf("expected 'a\nc\n' got %q", stdout.String())
	}
}

func TestDefaultParams(t *testing.T) {
	p := parseParams()
	if p.maxArgs != defaultMaxArgs {
		t.Errorf("expected %d, got %d", defaultMaxArgs, p.maxArgs)
	}
	if p.trace {
		t.Errorf("expected %t, got %t", false, p.trace)
	}
	if p.prompt {
		t.Errorf("expected %t, got %t", false, p.prompt)
	}
}
