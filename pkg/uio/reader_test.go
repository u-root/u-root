// Copyright 2021 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package uio

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func readAndCheck(t *testing.T, want, tmpfileP string) {
	t.Helper()
	r := strings.NewReader(want)
	if err := ReadIntoFile(r, tmpfileP); err != nil {
		t.Errorf("ReadIntoFile(%v, %s) = %v, want no error", r, tmpfileP, err)
	}

	got, err := os.ReadFile(tmpfileP)
	if err != nil {
		t.Fatalf("os.ReadFile(%s) = %v, want no error", tmpfileP, err)
	}
	if want != string(got) {
		t.Errorf("got: %v, want %s", string(got), want)
	}
}

func TestReadIntoFile(t *testing.T) {
	want := "I am the wanted"

	dir := t.TempDir()

	// Write to a file already exist.
	p := filepath.Join(dir, "uio-out")
	// Expect net effect to be creating a new empty file: "uio-out".
	f, err := os.OpenFile(p, os.O_RDONLY|os.O_CREATE|os.O_TRUNC, 0o755)
	if err != nil {
		t.Fatal(err)
	}
	if err := f.Close(); err != nil {
		t.Fatal(err)
	}
	readAndCheck(t, want, f.Name())

	// Write to a file that does not exist.
	p = filepath.Join(dir, "uio-out-not-existing")
	readAndCheck(t, want, p)

	// Write to an existing file that has pre-existing content.
	p = filepath.Join(dir, "uio-out-prexist-content")
	f, err = os.OpenFile(p, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o755)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := f.Write([]byte("temporary file's content")); err != nil {
		t.Fatal(err)
	}
	if err := f.Close(); err != nil {
		t.Fatal(err)
	}
	readAndCheck(t, want, p)
}
