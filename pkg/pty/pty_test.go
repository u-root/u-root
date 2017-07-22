// Copyright 2015-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pty

import (
	"reflect"
	"testing"
)

func TestNew(t *testing.T) {
	if _, err := New("/bin/bash", "-c", "/bin/date"); err != nil {
		t.Errorf("New pty: want nil, got %v", err)
	}
}

// This test is a nice idea but there's almost certainly no way
// to make it work.
func testStart(t *testing.T) {
	p, err := New("/bin/bash", "-c", "dd count=2 bs=1")
	if err != nil {
		t.Fatalf("TestStart New pty: want nil, got %v", err)
	}
	if err := p.Start(); err != nil {
		t.Fatalf("TestStart Start: want nil, got %v", err)
	}
	var b = [...]byte{'h', 'i'}
	if n, err := p.Ptm.Write(b[:]); n != len(b) || err != nil {
		t.Fatalf("Write to child: want (2, nil) got (%d, %v)", n, err)
	}
	t.Logf("Wrote message")
	if err := p.Wait(); err != nil {
		t.Fatalf("Wait for child: want nil got %v", err)
	}
	if n, err := p.Pts.Read(b[:]); n != 2 || err != nil {
		t.Fatalf("Read from child: want (2, nil) got (%d, %v)", n, err)
	}
}

func TestRunRestoreTTYMode(t *testing.T) {
	p, err := New("echo", "hi")
	if err != nil {
		t.Fatalf("TestStart New pty: want nil, got %v", err)
	}
	if err := p.Run(); err != nil {
		t.Fatalf("TestStart Start: want nil, got %v", err)
	}
	ti, err := p.TTY.Get()
	if err != nil {
		t.Fatalf("TestStart Get: want nil, got %v", err)
	}
	if !reflect.DeepEqual(ti, p.Restorer) {
		t.Errorf("TestStart: want termios from Get %v to be the same as termios from Start (%v) to be the same, they differ", ti, p.Restorer)
	}
}
