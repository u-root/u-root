// Copyright 2013-2023 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"errors"
	"os/exec"
	"strings"
	"testing"
)

func TestCommandNotFound(t *testing.T) {
	stdin := strings.NewReader("hello world")
	err := run(stdin, nil, nil, 1, "commandnotfound", "arg1")
	if !errors.Is(err, exec.ErrNotFound) {
		t.Fatalf("expected %v, got %v", exec.ErrNotFound, err)
	}
}

func TestEcho(t *testing.T) {
	stdin := strings.NewReader("hello world")
	stdout := &bytes.Buffer{}
	err := run(stdin, stdout, nil, defaultMaxArgs)
	if err != nil {
		t.Fatalf("expected nil, got %v", err)
	}

	if stdout.String() != "hello world\n" {
		t.Fatalf("expected 'hello world', got %q", stdout.String())
	}
}

func TestEchoWithMaxArgs(t *testing.T) {
	stdin := strings.NewReader("a b c d e f g")
	stdout := &bytes.Buffer{}
	err := run(stdin, stdout, nil, 3)
	if err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
	if stdout.String() != "a b c\nd e f\ng\n" {
		t.Fatalf("expected 'a b c\nd e f\ng\n', got %q", stdout.String())
	}
}
