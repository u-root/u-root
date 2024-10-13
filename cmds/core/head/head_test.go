// Copyright 2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestHead(t *testing.T) {
	dir := t.TempDir()
	f1, err := os.CreateTemp(dir, "head")
	if err != nil {
		t.Fatalf("expected nil got %v", err)
	}
	_, err = f1.WriteString("f11\nf12\nf13\nf14\nf15")
	if err != nil {
		t.Fatalf("expected nil got %v", err)
	}
	f2, err := os.CreateTemp(dir, "head")
	if err != nil {
		t.Fatalf("expected nil got %v", err)
	}
	_, err = f2.WriteString("f21\nf22\nf23")
	if err != nil {
		t.Fatalf("expected nil got %v", err)
	}
	f3, err := os.CreateTemp(dir, "head")
	if err != nil {
		t.Fatalf("expected nil got %v", err)
	}

	t.Run("combine error", func(t *testing.T) {
		err := run(nil, nil, nil, 1, 1, f3.Name())
		if !errors.Is(err, errCombine) {
			t.Errorf("expected %v, got %v", errCombine, err)
		}
	})

	t.Run("one file print lines", func(t *testing.T) {
		stdout := &bytes.Buffer{}
		err := run(nil, stdout, nil, 0, 2, f1.Name())
		if err != nil {
			t.Fatalf("expected nil, got %v", err)
		}

		expected := "f11\nf12\n"
		if stdout.String() != expected {
			t.Errorf("%v != %v", expected, stdout.String())
		}
	})

	t.Run("one file default params", func(t *testing.T) {
		stdout := &bytes.Buffer{}
		err := run(nil, stdout, nil, 0, 0, f2.Name())
		if err != nil {
			t.Fatalf("expected nil, got %v", err)
		}

		expected := "f21\nf22\nf23\n"
		if stdout.String() != expected {
			t.Errorf("%v != %v", expected, stdout.String())
		}
	})

	t.Run("two files print bytes", func(t *testing.T) {
		stdout := &bytes.Buffer{}
		err := run(nil, stdout, nil, 3, 0, f1.Name(), f2.Name())
		if err != nil {
			t.Fatalf("expected nil, got %v", err)
		}

		expected := fmt.Sprintf("==> %s <==\nf11\n==> %s <==\nf21",
			f1.Name(), f2.Name())

		if stdout.String() != expected {
			t.Errorf("%v != %v", expected, stdout.String())
		}
	})

	t.Run("request more bytes", func(t *testing.T) {
		stdout := &bytes.Buffer{}
		err := run(nil, stdout, nil, 10000, 0, f1.Name())
		if err != nil {
			t.Fatalf("expected nil, got %v", err)
		}

		expected := "f11\nf12\nf13\nf14\nf15"

		if stdout.String() != expected {
			t.Errorf("%v != %v", expected, stdout.String())
		}
	})

	t.Run("file not exists", func(t *testing.T) {
		stderr := &bytes.Buffer{}
		err := run(nil, nil, stderr, 0, 0, filepath.Join(dir, "filenotexists"))
		if err != nil {
			t.Fatalf("expected nil, got %v", err)
		}

		if !strings.Contains(stderr.String(), "filenotexists") {
			t.Errorf("expected file not exists in stderr, got %v", stderr.String())
		}
	})

	t.Run("stdin bytes", func(t *testing.T) {
		stdin := strings.NewReader("hello\n")
		stdout := &bytes.Buffer{}

		err := run(stdin, stdout, nil, 1, 0)
		if err != nil {
			t.Fatalf("expected nil, got %v", err)
		}

		if stdout.String() != "h" {
			t.Errorf("expected 'h' got %q", stdout.String())
		}
	})

	t.Run("stdin lines", func(t *testing.T) {
		stdin := strings.NewReader("hello\nagain\n")
		stdout := &bytes.Buffer{}

		err := run(stdin, stdout, nil, 0, 1)
		if err != nil {
			t.Fatalf("expected nil, got %v", err)
		}

		if stdout.String() != "hello\n" {
			t.Errorf("expected 'hello\n' got %q", stdout.String())
		}
	})
}
