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
	err := run(stdin, nil, nil, "commandnotfound", "arg1")
	if !errors.Is(err, exec.ErrNotFound) {
		t.Fatalf("expected %v, got %v", exec.ErrNotFound, err)
	}
}

func TestEcho(t *testing.T) {
	stdin := strings.NewReader("hello world")
	stdout := &bytes.Buffer{}
	err := run(stdin, stdout, nil)
	if err != nil {
		t.Fatalf("expected nil, got %v", err)
	}

	if stdout.String() != "hello world\n" {
		t.Fatalf("expected 'hello world', got %q", stdout.String())
	}
}
