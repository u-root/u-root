// Copyright 2012-2023 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"errors"
	"os"
	"testing"
)

func TestTac(t *testing.T) {
	f, err := os.CreateTemp("", "")
	if err != nil {
		t.Fatalf(`os.CreateTemp("", "") = %v, want nil`, err)
	}

	_, err = f.WriteString("hello\nworld\n")
	if err != nil {
		t.Fatalf(`f.WriteString("hello\nworld\n") = %v, want nil`, err)
	}

	stdout := &bytes.Buffer{}

	err = tac(stdout, []string{f.Name(), f.Name()})
	if err != nil {
		t.Fatalf(`tac(stdout, []string{f.Name(), f.Name()}) = %v, want nil`, err)
	}

	expected := "world\nhello\nworld\nhello\n"
	if stdout.String() != expected {
		t.Errorf("expected %s, got %s", expected, stdout.String())
	}
}

func TestTacStdin(t *testing.T) {
	err := tac(nil, nil)
	if !errors.Is(err, errStdin) {
		t.Errorf("expected %v, got %v", errStdin, err)
	}
}
