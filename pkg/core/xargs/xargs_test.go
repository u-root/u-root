// Copyright 2013-2023 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package xargs

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
	var stdout, stderr bytes.Buffer
	cmd := New()
	cmd.SetIO(strings.NewReader("hello world"), &stdout, &stderr)
	err := cmd.Run("-n", "1", "commandnotfound", "arg1")
	if !errors.Is(err, exec.ErrNotFound) {
		t.Errorf("expected %v, got %v", exec.ErrNotFound, err)
	}
}

func TestEcho(t *testing.T) {
	var stdout, stderr bytes.Buffer
	cmd := New()
	cmd.SetIO(strings.NewReader("hello world"), &stdout, &stderr)
	err := cmd.Run()
	if err != nil {
		t.Fatalf("expected nil, got %v", err)
	}

	if stdout.String() != "hello world\n" {
		t.Errorf("expected 'hello world', got %q", stdout.String())
	}
}

func TestEchoWithMaxArgs(t *testing.T) {
	var stdout, stderr bytes.Buffer
	cmd := New()
	cmd.SetIO(strings.NewReader("a b c d e f g"), &stdout, &stderr)
	err := cmd.Run("-n", "3", "-t")
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

func TestEchoPrompt(t *testing.T) {
	var stdout, stderr bytes.Buffer

	dir := t.TempDir()
	path := filepath.Join(dir, "tty")
	err := os.WriteFile(path, []byte("yes\nn\ny\n"), 0o644)
	if err != nil {
		t.Fatalf("expected nil, got %v", err)
	}

	cmd := New().(*command)
	SetTTY(cmd, path)
	cmd.SetIO(strings.NewReader("a b c"), &stdout, &stderr)
	err = cmd.Run("-n", "1", "-p")
	if err != nil {
		t.Fatalf("expected nil, got %v", err)
	}

	if stdout.String() != "a\nc\n" {
		t.Errorf("expected 'a\nc\n' got %q", stdout.String())
	}
}

func TestNullDelimiter(t *testing.T) {
	var stdout, stderr bytes.Buffer
	cmd := New()
	cmd.SetIO(strings.NewReader("hello\x00world"), &stdout, &stderr)
	err := cmd.Run("-0")
	if err != nil {
		t.Fatalf("expected nil, got %v", err)
	}

	if stdout.String() != "hello world\n" {
		t.Errorf("expected 'hello world', got %q", stdout.String())
	}
}
