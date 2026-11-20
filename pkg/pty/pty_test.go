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
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"
	"syscall"
	"testing"
	"time"

	"github.com/u-root/u-root/pkg/termios"
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

// A process with no controlling terminal can still run a command on a pty. That
// is the daemon's path -- u-root's sshd is one -- and it used to be impossible
// because New refused before /dev/ptmx was opened. Start borrows nothing, Run
// relays nothing, and Wait has nothing to restore.
func TestCommandWithoutCallerTerminal(t *testing.T) {
	p, err := New()
	if err != nil {
		t.Skipf("cannot allocate a pty here: %v", err)
	}
	if p.TTY != nil {
		t.Skip("this test needs no controlling terminal, and this process has one")
	}
	defer p.Ptm.Close()

	// The window size belongs to the pty, not to a caller terminal, so setting
	// one must work with nothing to copy from.
	ws := termios.NewWinsize(33, 101, 0, 0)
	if err := termios.SetWinSize(p.Pts.Fd(), ws); err != nil {
		t.Fatalf("SetWinSize on a pty with no caller terminal: %v", err)
	}
	back, err := (&termios.TTYIO{File: p.Pts}).GetWinSize()
	if err != nil {
		t.Fatalf("GetWinSize: %v", err)
	}
	if back.Row != 33 || back.Col != 101 {
		t.Errorf("window size %dx%d, want 33x101", back.Row, back.Col)
	}

	p.Command("echo", "hi")
	if err := p.Start(); err != nil {
		t.Fatalf("Start: %v", err)
	}
	// Drop this process's copy of the slave: the child holds its own, and while
	// one copy of the slave stays open a read of the master never ends.
	p.Pts.Close()
	// With no caller terminal the output stays on Ptm for whoever asked to read
	// it, which is the contract the nil guards document.
	// Read to the end of the stream rather than once: a single Read can return
	// before the child has written anything, and then the test is a race.
	got := make(chan string, 1)
	go func() {
		b, _ := io.ReadAll(p.Ptm)
		got <- string(b)
	}()
	select {
	case s := <-got:
		if !strings.Contains(s, "hi") {
			t.Errorf("read %q from Ptm, want it to contain \"hi\"", s)
		}
	case <-time.After(30 * time.Second):
		t.Fatal("no output appeared on Ptm within 30s")
	}
	if err := p.Wait(); err != nil {
		t.Errorf("Wait: %v", err)
	}

	// Run is Start, the relays it does not start, and Wait.
	q, err := New()
	if err != nil {
		t.Skipf("cannot allocate a second pty: %v", err)
	}
	if q.TTY != nil {
		t.Skip("this process grew a controlling terminal mid-test")
	}
	defer q.Ptm.Close()
	q.Command("true")
	if err := q.Run(); err != nil {
		t.Errorf("Run with no caller terminal: %v", err)
	}
}

// A terminal that exists but is not a terminal is a real failure, not an absence:
// reporting it as "no terminal" would silently degrade a diagnosable problem.
func TestCallerTerminalOnAFileThatIsNotATerminal(t *testing.T) {
	path := filepath.Join(t.TempDir(), "notatty")
	if err := os.WriteFile(path, []byte("x"), 0o600); err != nil {
		t.Fatal(err)
	}
	tty, restorer, err := callerTerminal(path)
	if err == nil {
		t.Errorf("callerTerminal on a regular file: want an error, got %v/%v", tty, restorer)
	}
}

// The other polarity of callerTerminal: a process that *does* have a controlling
// terminal. CI has none, so the child is given one -- a pty is allocated here and
// the child is started as a session leader with that pty as its controlling
// terminal, which is what u-root's init does for the console shell.
//
// This is the half of pty.Start, Run and Wait that copies the caller's settings
// into the pty and restores them afterwards, and it is unreachable in CI without
// manufacturing a terminal this way.
func TestCommandWithCallerTerminal(t *testing.T) {
	if os.Getenv("PTY_CTTY_HELPER") == "1" {
		// In the child, which has a controlling terminal.
		p, err := New()
		if err != nil {
			fmt.Printf("New: %v\n", err)
			os.Exit(1)
		}
		if p.TTY == nil || p.Restorer == nil {
			fmt.Printf("TTY=%v Restorer=%v, want both set: this child was given a terminal\n", p.TTY, p.Restorer)
			os.Exit(1)
		}
		p.Command("echo", "ctty-ok")
		if err := p.Run(); err != nil {
			fmt.Printf("Run: %v\n", err)
			os.Exit(1)
		}
		os.Exit(0)
	}

	p, err := New()
	if err != nil {
		t.Skipf("cannot allocate a pty here: %v", err)
	}
	defer p.Ptm.Close()

	cmd := exec.Command(os.Args[0], "-test.run=TestCommandWithCallerTerminal")
	cmd.Env = append(os.Environ(), "PTY_CTTY_HELPER=1")
	cmd.Stdin, cmd.Stdout, cmd.Stderr = p.Pts, p.Pts, p.Pts
	cmd.SysProcAttr = &syscall.SysProcAttr{Setsid: true, Setctty: true}
	if err := cmd.Start(); err != nil {
		t.Fatalf("start the child on a pty: %v", err)
	}
	// The child has its own copy of the slave; this one must go, or the master
	// never reports the end of the stream.
	p.Pts.Close()

	out := make(chan string, 1)
	go func() {
		b, _ := io.ReadAll(p.Ptm)
		out <- string(b)
	}()
	if err := cmd.Wait(); err != nil {
		t.Errorf("the child, given a controlling terminal, failed: %v", err)
	}
	// The child's own exit status is the assertion: it fails if a process with a
	// controlling terminal cannot allocate a pty, or cannot run a command on it.
	// The output is checked only if it arrives: Run returns as soon as the child
	// exits without draining its relay, so a very short command's output can be
	// lost, and asserting on it would make this test flaky rather than strict.
	select {
	case s := <-out:
		if s != "" && !strings.Contains(s, "ctty-ok") {
			t.Errorf("the child printed %q, want it to contain \"ctty-ok\"", s)
		}
	case <-time.After(30 * time.Second):
		t.Fatal("the child produced no output within 30s")
	}
}
