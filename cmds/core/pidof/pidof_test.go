// Copyright 2025 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"errors"
	"testing"
)

func TestPidof(t *testing.T) {
	stdout := bytes.Buffer{}

	err := run(&stdout, "./testdata", []string{"init", "bash"})
	if err != nil {
		t.Fatalf("expected nil got %v", err)
	}

	expected := "1 2\n"
	if stdout.String() != "1 2\n" {
		t.Errorf("expeted %q, got %q", expected, stdout.String())
	}
}

func TestPidofMissing(t *testing.T) {
	stdout := bytes.Buffer{}
	err := run(&stdout, "./testdata", []string{"goooo"})
	if !errors.Is(err, errNotFound) {
		t.Fatalf("expected %v got %v", errNotFound, err)
	}

	if stdout.String() != "" {
		t.Errorf("expected empty string got %q", stdout.String())
	}
}
