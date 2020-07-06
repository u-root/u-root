// Copyright 2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package lockfile

import (
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func testProcess(t *testing.T) *os.Process {
	p := exec.Command("sleep", "1000")
	if err := p.Start(); err != nil {
		t.Fatal(err)
	}

	go p.Wait()
	return p.Process
}

func TestTryLock(t *testing.T) {
	p1 := testProcess(t)
	defer p1.Kill()
	p2 := testProcess(t)
	defer p2.Kill()

	dir, err := ioutil.TempDir("", "lockfile-")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	l := &Lockfile{
		path: filepath.Join(dir, "test1"),
		pid:  p1.Pid,
	}
	if err := l.TryLock(); err != nil {
		t.Fatalf("l(%d).TryLock() = %v, want nil", l.pid, err)
	}

	l2 := &Lockfile{
		path: l.path,
		pid:  p2.Pid,
	}
	if err := l2.TryLock(); err != ErrBusy {
		t.Fatalf("l(%d).TryLock() = %v, want ErrBusy", l2.pid, err)
	}

	if err := l.Unlock(); err != nil {
		t.Fatal(err)
	}

	if err := l.Unlock(); err != ErrRogueDeletion {
		t.Fatalf("2nd Unlock() = %v, want ErrRogueDeletion", err)
	}

	if err := l2.TryLock(); err != nil {
		t.Fatalf("l(%d).TryLock() = %v, want nil", l2.pid, err)
	}
}

func TestLockFileRemoval(t *testing.T) {
	p := testProcess(t)
	defer p.Kill()

	dir, err := ioutil.TempDir("", "lockfile-")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	l := &Lockfile{
		path: filepath.Join(dir, "test2"),
		pid:  p.Pid,
	}
	if err := l.TryLock(); err != nil {
		t.Fatalf("l(%d).TryLock() = %v, want nil", l.pid, err)
	}

	// Evil actor deletes the lockfile.
	if err := os.Remove(l.path); err != nil {
		t.Fatalf("remove(%v) = %v, want nil", l.path, err)
	}

	if err := l.Unlock(); err != ErrRogueDeletion {
		t.Fatalf("l(%d).Unlock() = %v, want ErrRogueDeletion", l.pid, err)
	}

	if err := l.TryLock(); err != nil {
		t.Fatalf("l(%d).TryLock() = %v, want nil", l.pid, err)
	}
}

func TestDeadProcess(t *testing.T) {
	p1 := testProcess(t)
	defer p1.Kill()
	p2 := testProcess(t)
	defer p2.Kill()

	dir, err := ioutil.TempDir("", "lockfile-")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	l := &Lockfile{
		path: filepath.Join(dir, "test3"),
		pid:  p1.Pid,
	}
	if err := l.TryLock(); err != nil {
		t.Fatalf("l(%d).TryLock() = %v, want nil", l.pid, err)
	}
	if err := p1.Kill(); err != nil {
		t.Fatalf("Kill() = %v, want nil", err)
	}
	p1.Wait()

	l2 := &Lockfile{
		path: l.path,
		pid:  p2.Pid,
	}
	if err := l2.TryLock(); err != nil {
		t.Fatalf("l(%d).TryLock() = %v, want nil", l2.pid, err)
	}
	if err := l2.Unlock(); err != nil {
		t.Fatalf("l(%d).Unlock() = %v, want nil", l2.pid, err)
	}
}
