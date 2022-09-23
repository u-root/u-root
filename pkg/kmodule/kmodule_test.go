// Copyright 2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build darwin || linux

package kmodule

import (
	"bytes"
	"errors"
	"io"
	"os"
	"path/filepath"
	"testing"
)

func TestNewPath(t *testing.T) {
	l, err := NewPath("proc.modules")
	if err != nil {
		t.Fatalf(`NewPath("proc.modules"): %v != nil`, err)
	}

	// At the very minimum, if we have it, we have list.
	var out = &bytes.Buffer{}
	if _, err := io.Copy(out, l); err != nil {
		t.Fatalf("List: got %v, want nil", err)
	}
	t.Logf("List: %s", out.String())
}

func TestNewPathBad(t *testing.T) {
	d := t.TempDir()
	bad := filepath.Join(d, "bad")
	if _, err := NewPath(bad); !errors.Is(err, os.ErrNotExist) {
		t.Errorf("NewPath(%q): %v != %v", bad, err, os.ErrNotExist)
	}
}

// badloader is a loader that is unable to do anything right.
type badloader struct {
	err error
}

func (b *badloader) Init(_ []byte, _ string) error {
	return b.err
}

func (b *badloader) FileInit(_ *os.File, _ string, _ uintptr) error {
	return b.err
}

func (b *badloader) Probe(_, _ string) error {
	return b.err
}

func (b *badloader) Delete(_ string, _ uintptr) error {
	return b.err
}

func (b *badloader) Read(_ []byte) (int, error) {
	// N.B.: io.ReadAll uses the returned byte length
	// before checking to see if there was an error.
	// Returning -1 is a bad idea.
	return 0, b.err
}

var _ Loader = &badloader{}

func TestListBadLoader(t *testing.T) {
	if err := List(&badloader{err: os.ErrInvalid}, nil); !errors.Is(err, os.ErrInvalid) {
		t.Errorf("Calling List with badloader: got %v, want %v", err, os.ErrInvalid)
	}
}

func TestList(t *testing.T) {
	l, err := NewPath("proc.modules")
	if err != nil {
		t.Fatalf(`NewPath("proc.modules"): %v != nil`, err)
	}

	var b bytes.Buffer
	if err := List(l, &b); err != nil {
		t.Errorf("List: %v != nil", err)
	}
}
