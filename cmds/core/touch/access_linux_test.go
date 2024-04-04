// Copyright 2023 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build linux
// +build linux

package main

import (
	"os"
	"syscall"
	"testing"
	"time"
)

func TestAccess(t *testing.T) {
	tmp := t.TempDir()
	f, err := os.CreateTemp(tmp, "touch_test")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	accessDate, err := time.Parse(time.RFC3339, "2023-01-01T00:00:00Z")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	cmd := command(params{time: accessDate, access: true}, f.Name())
	err = cmd.run()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	mfi, err := f.Stat()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	expected := accessDate.UnixNano()
	at := mfi.Sys().(*syscall.Stat_t).Atim.Nano()
	if at != expected {
		t.Errorf("expected access time %v, got %v", expected, at)
	}
}
