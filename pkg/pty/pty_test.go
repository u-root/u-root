// Copyright 2015-2020 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// This test is flaky AF under the race detector.

//go:build !race

package pty

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
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
	// Restoring the caller's terminal modes needs a caller with a terminal.
	// New no longer refuses when there is none, so this is the skip that says so
	// rather than the ENXIO above, which used to cover both cases.
	if p.TTY == nil {
		t.Skipf("no controlling terminal to restore modes on")
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

// A process without a controlling terminal can still allocate a pty. This is not
// an edge case: u-root's sshd is started from init, has no terminal of its own,
// and must be able to give one to the commands it runs. Before callerTerminal,
// pty.New failed outright with ENXIO on /dev/tty and sshd refused every pty-req.
func TestNewWithoutControllingTerminal(t *testing.T) {
	if os.Getenv("PTY_NO_CTTY_HELPER") == "1" {
		// In the child, which was started with Setsid and so has no controlling
		// terminal: allocating a pty must succeed, with nothing to restore.
		p, err := New()
		if err != nil {
			fmt.Printf("New: %v\n", err)
			os.Exit(1)
		}
		if p.TTY != nil || p.Restorer != nil {
			fmt.Printf("TTY=%v Restorer=%v, want both nil\n", p.TTY, p.Restorer)
			os.Exit(1)
		}
		if _, err := p.Pts.WriteString("ok"); err != nil {
			fmt.Printf("write to Pts: %v\n", err)
			os.Exit(1)
		}
		os.Exit(0)
	}

	cmd := exec.Command(os.Args[0], "-test.run=TestNewWithoutControllingTerminal")
	cmd.Env = append(os.Environ(), "PTY_NO_CTTY_HELPER=1")
	cmd.SysProcAttr = &syscall.SysProcAttr{Setsid: true}
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("a process with no controlling terminal could not allocate a pty: %v: %s", err, out)
	}
}

// callerTerminal reports "no terminal" as nothing to copy, and a real failure as
// an error, so that a pty is not refused for a missing terminal.
func TestCallerTerminalDistinguishesMissingFromBroken(t *testing.T) {
	tty, restorer, err := callerTerminal("/nonexistent/tty")
	if err != nil {
		t.Fatalf("callerTerminal on a missing device: %v, want no error", err)
	}
	if tty != nil || restorer != nil {
		t.Errorf("callerTerminal on a missing device = %v, %v, want nil, nil", tty, restorer)
	}
}
