// Copyright 2015-2020 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// This test is flaky AF under the race detector.

//go:build !race

package pty

import (
	"encoding/json"
	"errors"
	"os"
	"reflect"
	"syscall"
	"testing"
)

func TestNew(t *testing.T) {
	if _, err := New(); os.IsNotExist(err) || errors.Is(err, syscall.ENXIO) {
		t.Skipf("Failed allocate /dev/pts device")
	} else if err != nil {
		t.Errorf("New pty: want nil, got %v", err)
	}
}

func TestRunRestoreTTYMode(t *testing.T) {
	p, err := New()
	if os.IsNotExist(err) || errors.Is(err, syscall.ENXIO) {
		t.Skipf("Failed to allocate /dev/pts device")
	} else if err != nil {
		t.Fatalf("TestStart New pty: want nil, got %v", err)
	}

	p.Command("echo", "hi")
	if err := p.Start(); err != nil {
		t.Fatalf("TestStart Start: want nil, got %v", err)
	}
	if err := p.Wait(); err != nil {
		t.Error(err)
	}
	ti, err := p.TTY.Get()
	if err != nil {
		t.Fatalf("TestStart Get: want nil, got %v", err)
	}
	if !reflect.DeepEqual(ti, p.Restorer) {
		tt, err := json.Marshal(ti)
		if err != nil {
			t.Fatalf("Can't marshall %v: %v", ti, err)
		}
		r, err := json.Marshal(p.Restorer)
		if err != nil {
			t.Fatalf("Can't marshall %v: %v", p.Restorer, err)
		}
		t.Errorf("TestStart: want termios from Get %s to be the same as termios from Start (%s) to be the same, they differ", tt, r)
	}
	b := make([]byte, 1024)
	n, err := p.Ptm.Read(b)
	t.Logf("ptm read is %d bytes, b is %q", n, b[:n])
	if err != nil {
		t.Fatalf("Error reading from process: %v", err)
	}
	if n != 4 {
		t.Errorf("Bogus returned amount: got %d, want 4", n)
	}
	if string(b[:n]) != "hi\r\n" {
		t.Errorf("bogus returned data: got %q, want %q", string(b[:n]), "hi\r\n")
	}
}
