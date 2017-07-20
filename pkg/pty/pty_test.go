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

func TestStart(t *testing.T) {
	p, err := New("/bin/bash", "-c", "cat")
	if err != nil {
		t.Fatalf("TestStart New pty: want nil, got %v", err)
	}
	if err := p.Start(); err != nil {
		t.Fatalf("TestStart Start: want nil, got %v", err)
	}
	var b = [...]byte{'h', 'i'}
	if n, err := p.Pts.Write(b[:]); n != len(b) || err != nil {
		t.Fatalf("Write to child: want (2, nil) got (%d, %v)", n, err)
	}
	if n, err := p.Pts.Write([]byte{4}); n != 1 || err != nil {
		t.Fatalf("Write ^D to child: want (1, nil) got (%d, %v)", n, err)
	}
	if n, err := p.Pts.Read(b[:]); n != 2 || err != nil {
		t.Fatalf("Read from child: want (2, nil) got (%d, %v)", n, err)
	}
}

func TestRun(t *testing.T) {
	p, err := New("/bin/bash", "-c", "/bin/echo", "-n", "hi")
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
	// the process is running. Send it a string and see it comes back.
	if n, err := p.Pts.Write([]byte("hi")); n != 2 || err != nil {
		t.Errorf("Writing to child: want (2, nil) got (%d, %v)", n, err)
	}
}
