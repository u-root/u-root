// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.


package term

import "testing"

// TestSSH tests the Termios<>SSH term attribute conversions.
func TestSSH(t *testing.T) {
	pty, err := OpenPTY()
	if err != nil {
		t.Fatalf("OpenPTY failed: %v", err)
	}
	defer pty.Close()
	sTerm, err := Attr(pty.Slave)
	if err != nil {
		t.Fatalf("Reading PTY slave attributes failed: %v", err)
	}
	sTerm.Ispeed, sTerm.Ospeed = 19200, 300 // Normally set to 0 so putting something in there for the test.
	var dTerm Termios
	dTerm.FromSSH(sTerm.ToSSH())
	sTerm.Cflag &= CS7 | CS8 | PARENB | PARODD // The only term control modes in the SSH standard.
	if sTerm != dTerm {
		t.Errorf("terminal modes does not match, got: %v want: %v", dTerm, sTerm)
	}
	sTerm.Raw()
	dTerm.FromSSH(sTerm.ToSSH())
	if err := testraw(dTerm, "sshTerm"); err != nil {
		t.Errorf("TestSSH failed: %v", err)
	}
	sTerm.Cook()
	dTerm.FromSSH(sTerm.ToSSH())
	dTerm.Iflag += BRKINT // Not in the SSH standard but set for Linux cooked mode.
	if err := testcook(dTerm, "sshTerm"); err != nil {
		t.Errorf("TestSSH failed: %v", err)
	}
}
