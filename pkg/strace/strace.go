// Copyright 2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package strace supports tracing programs.
// The basic control of tracing is via a Tracer, which returns raw
// TraceRecords via a chan. The easiest way to create a Tracer is via
// RunTracerFromCommand, which uses a filled out exec.Cmd to start a
// process and produce trace records.
// Forking and nice printing are not yet supported.
package strace

import (
	"fmt"
	"os/exec"
	"runtime"
	"syscall"

	"golang.org/x/sys/unix"
)

var Debug = func(string, ...interface{}) {}

// Tracer has information to trace one process. It can be created by
// starting a command, or attaching. Attaching is not supported yet.
type Tracer struct {
	Pid     int
	Records chan *TraceRecord
	Count   int
}

func New() (*Tracer, error) {
	return &Tracer{Pid: -1, Records: make(chan *TraceRecord)}, nil
}

// StartTracerFromCommand runs a Tracer given an exec.Cmd.
func (t *Tracer) RunTracerFromCmd(c *exec.Cmd) {
	defer close(t.Records)
	if c.SysProcAttr == nil {
		c.SysProcAttr = &syscall.SysProcAttr{}
	}
	c.SysProcAttr.Ptrace = true
	// Because the go runtime forks traced processes with PTRACE_TRACEME
	// we need to maintain the parent-child relationship for ptrace to work.
	// We've learned this the hard way. So we lock down this thread to
	// this proc, and start the command here.
	// Note this function will block; if you want it to be nonblocking you
	// need to use go etc.
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	if err := c.Start(); err != nil {
		Debug("Start gets err %v", err)
		t.Records <- &TraceRecord{Err: err}
		return
	}
	Debug("Start gets pid %v", c.Process.Pid)
	if err := c.Wait(); err != nil {
		fmt.Printf("Wait returned: %v\n", err)
		t.Records <- &TraceRecord{Err: err}
	}
	t.Pid = c.Process.Pid
	t.Run()
}

func (t *Tracer) Run() error {
	x := Enter
	for {
		t.Count++
		r := &TraceRecord{Serial: t.Count, EX: x, Pid: t.Pid}
		if err := unix.PtraceGetRegs(t.Pid, &r.Regs); err != nil {
			Debug("ptracegetregs for %d gets %v", t.Pid, err)
			r.Err = err
			t.Records <- r
			break
		}
		t.Records <- r
		switch x {
		case Enter:
			x = Exit
		default:
			x = Enter
		}
		if err := unix.PtraceSyscall(t.Pid, 0); err != nil {
			r := &TraceRecord{Serial: t.Count, EX: x, Pid: t.Pid}
			Debug("ptracesyscall for %d gets %v", t.Pid, err)
			r.Err = err
			t.Records <- r
			break
		}
		if w, err := Wait(t.Pid); err != nil {
			r := &TraceRecord{Serial: t.Count, EX: x, Pid: t.Pid}
			Debug("wait4 for %d gets %v, %v", t.Pid, w, err)
			r.Err = err
			t.Records <- r
			break
		}

	}
	Debug("Pushed %d records", t.Count)
	return nil
}

// EventType describes whether a record is system call Entry or Exit
type EventType string

const (
	Enter EventType = "E"
	Exit            = "X"
)

// TraceRecord has information about a ptrace event.
type TraceRecord struct {
	EX     EventType
	Regs   unix.PtraceRegs
	Serial int
	Pid    int
	Err    error
}

// String is a stringer for TraceRecords
// TODO: stringer for Regs.
func (t *TraceRecord) String() string {
	pre := fmt.Sprintf("%s %d#%d:", t.EX, t.Pid, t.Serial)
	if t.Err != nil {
		return fmt.Sprintf("%s(%v)", pre, t.Err)
	}
	return fmt.Sprintf("%s %v", pre, t.Regs)
}
