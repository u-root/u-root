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
