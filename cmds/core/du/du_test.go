// Copyright 2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package main

import (
	"bytes"
	"os"
	"regexp"
	"testing"
)

func TestDU(t *testing.T) {
	t.Run("8K file", func(t *testing.T) {
		dir := t.TempDir()
		f, err := os.CreateTemp(dir, "")
		if err != nil {
			t.Fatal(err)
		}
		f.Write(make([]byte, 8096))

		blocks, err := du(f.Name())
		if err != nil {
			t.Fatalf("expected nil got %v", err)
		}

		if blocks != 16 {
			t.Errorf("expected 16 blocks, got %d", blocks)
		}
	})
	t.Run("empty file", func(t *testing.T) {
		dir := t.TempDir()
		f, err := os.CreateTemp(dir, "")
		if err != nil {
			t.Fatal(err)
		}

		blocks, err := du(f.Name())
		if err != nil {
			t.Fatalf("expected nil got %v", err)
		}

		if blocks != 0 {
			t.Errorf("expected 0 blocks, got %d", blocks)
		}
	})
	t.Run("one bit file", func(t *testing.T) {
		dir := t.TempDir()
		f, err := os.CreateTemp(dir, "")
		if err != nil {
			t.Fatal(err)
		}
		f.Write(make([]byte, 1))

		blocks, err := du(f.Name())
		if err != nil {
			t.Fatalf("expected nil got %v", err)
		}

		if blocks != 8 {
			t.Errorf("expected 8 blocks, got %d", blocks)
		}
	})
}

func TestRun(t *testing.T) {
	t.Run("empty folder", func(t *testing.T) {
		dir := t.TempDir()
		err := os.Chdir(dir)
		if err != nil {
			t.Fatal(err)
		}

		stdout := &bytes.Buffer{}
		err = run(stdout)
		if err != nil {
			t.Fatalf("expected nil got %v", err)
		}

		r := regexp.MustCompile(`^\d\t\.\n$`)
		if !r.MatchString(stdout.String()) {
			t.Error("expected number tab dot new-line")
		}
	})
}
