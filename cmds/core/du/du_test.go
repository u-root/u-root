// Copyright 2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package main

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"testing"

	"github.com/u-root/u-root/pkg/testutil"
)

func TestDU(t *testing.T) {
	t.Run("8K file", func(t *testing.T) {
		dir := t.TempDir()
		f, err := os.CreateTemp(dir, "")
		if err != nil {
			t.Fatal(err)
		}
		f.Write(make([]byte, 8096))

		cmd := &cmd{stdout: io.Discard}
		blocks, err := cmd.du(f.Name())
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

		cmd := &cmd{stdout: io.Discard}
		blocks, err := cmd.du(f.Name())
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

		cmd := &cmd{stdout: io.Discard}
		blocks, err := cmd.du(f.Name())
		if err != nil {
			t.Fatalf("expected nil got %v", err)
		}

		if blocks != 8 {
			t.Errorf("expected 8 blocks, got %d", blocks)
		}
	})
	t.Run("follow symlink", func(t *testing.T) {
		d1 := t.TempDir()
		f, err := os.CreateTemp(d1, "")
		if err != nil {
			t.Fatal(err)
		}
		f.Write(make([]byte, 4096))

		d2 := t.TempDir()
		sl := filepath.Join(d2, "symlink")
		err = os.Symlink(f.Name(), sl)
		if err != nil {
			t.Fatal(err)
		}

		cmd := &cmd{stdout: io.Discard, followSymlinks: true}
		blocks, err := cmd.du(sl)
		if err != nil {
			t.Fatal(err)
		}

		if blocks != 8 {
			t.Errorf("expected 8 blocks, got %d", blocks)
		}
	})
}

func TestRun(t *testing.T) {
	t.Run("empty folder", func(t *testing.T) {
		currDir, err := os.Getwd()
		if err != nil {
			t.Fatal(err)
		}

		dir := t.TempDir()
		err = os.Chdir(dir)
		if err != nil {
			t.Fatal(err)
		}
		defer os.Chdir(currDir)

		stdout := &bytes.Buffer{}
		err = command(stdout, []string{"du"}).run()
		if err != nil {
			t.Fatalf("expected nil got %v", err)
		}

		r := regexp.MustCompile(`^\d\t\.\n$`)
		if !r.MatchString(stdout.String()) {
			t.Error("expected number tab dot new-line")
		}
	})
	t.Run("report all files", func(t *testing.T) {
		dir := prepareDir(t)
		stdout := &bytes.Buffer{}
		err := command(stdout, []string{"du", "-a", dir}).run()
		if err != nil {
			t.Fatalf("expected nil got %v", err)
		}
		lines := strings.Split(strings.TrimSpace(stdout.String()), "\n")

		// should print dir, subdir and two files
		if len(lines) != 4 {
			t.Errorf("expected file1, file2 and temp dir, but got %d lines", len(lines))
		}
	})
	t.Run("with -k flag", func(t *testing.T) {
		dir := t.TempDir()
		f, err := os.CreateTemp(dir, "")
		if err != nil {
			t.Fatal(err)
		}
		f.Write(make([]byte, 4096))

		stdout := &bytes.Buffer{}
		err = command(stdout, []string{"du", "-k", f.Name()}).run()
		if err != nil {
			t.Fatalf("expected nil got %v", err)
		}

		if stdout.String()[0] != '4' {
			t.Errorf("expected 4 blocks with -k, got %q", stdout.String())
		}
	})
	t.Run("total sum", func(t *testing.T) {
		dir := prepareDir(t)
		stdout := &bytes.Buffer{}
		err := command(stdout, []string{"du", "-s", dir}).run()
		if err != nil {
			t.Fatalf("expected nil got %v", err)
		}
		lines := strings.Split(strings.TrimSpace(stdout.String()), "\n")
		if len(lines) != 1 {
			t.Errorf("expected one line per file with -s flag, got %d", len(lines))
		}
	})
	t.Run("both -s and -a", func(t *testing.T) {
		err := command(io.Discard, []string{"du", "-sa"}).run()
		if err == nil {
			t.Errorf("expected %v, got %v", errUsage, err)
		}
	})
	t.Run("both -H and -L", func(t *testing.T) {
		err := command(io.Discard, []string{"du", "-HL"}).run()
		if err == nil {
			t.Errorf("expected %v, got %v", errUsage, err)
		}
	})
	t.Run("symlink with -h", func(t *testing.T) {
		dir := prepareDir(t)
		slDir := filepath.Join(t.TempDir(), "sl")
		err := os.Symlink(dir, slDir)
		if err != nil {
			t.Fatal(err)
		}

		stdout := &bytes.Buffer{}
		err = command(stdout, []string{"du", slDir}).run()
		if err != nil {
			t.Fatalf("expected nil got %v", err)
		}
		r := parseBlocks(t, stdout)
		if r != 0 {
			t.Errorf("expected 0 got %d", r)
		}

		stdout.Reset()
		err = command(stdout, []string{"du", "-H", slDir}).run()
		if err != nil {
			t.Fatalf("expected nil got %v", err)
		}
		if parseBlocks(t, stdout) == 0 {
			t.Error("expected du to follow symlink but value was 0")
		}
	})
	t.Run("relative symlinks", func(t *testing.T) {
		testutil.SkipIfInVMTest(t)

		stdout := &bytes.Buffer{}
		err := command(stdout, []string{"du", "-s", "testdata/rel-symlinks"}).run()
		if err != nil {
			t.Fatalf("expected nil got %v", err)
		}

		sizeWithoutL := parseBlocks(t, stdout)

		stdout.Reset()
		err = command(stdout, []string{"du", "-sL", "testdata/rel-symlinks"}).run()
		if err != nil {
			t.Fatalf("expected nil got %v", err)
		}
		sizeWithL := parseBlocks(t, stdout)

		if sizeWithL <= sizeWithoutL {
			t.Errorf("expected size with -L to be bigger, because to relative symlinks, got -s: %d, -sL: %d", sizeWithoutL, sizeWithL)
		}
	})
}

// parse first line of the output
func parseBlocks(t *testing.T, stdout *bytes.Buffer) int64 {
	t.Helper()
	sp := strings.Fields(stdout.String())
	r, err := strconv.ParseInt(sp[0], 10, 64)
	if err != nil {
		t.Fatal(err)
	}

	return r
}

func prepareDir(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()

	f1, err := os.Create(filepath.Join(dir, "file1"))
	if err != nil {
		t.Fatal(err)
	}
	f1.Write(make([]byte, 4096))
	dir1 := filepath.Join(dir, "dir1")
	err = os.Mkdir(dir1, 0o722)
	if err != nil {
		t.Fatal(err)
	}
	f2, err := os.Create(filepath.Join(dir1, "file2"))
	if err != nil {
		t.Fatal(err)
	}
	f2.Write(make([]byte, 8012))
	if err != nil {
		t.Fatal(err)
	}

	return dir
}
