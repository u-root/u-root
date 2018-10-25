// Copyright 2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package strace supports tracing programs.
// The basic control of tracing is via a Tracer, which returns raw
// TraceRecords via a chan. The easiest way to create a Tracer is via
// RunTracerFromCommand, which uses a filled out exec.Cmd to start a
// process and produce trace records.
// Forking is not yet supported.
package strace

import (
	"encoding/binary"
	"fmt"
	"io"
	"os/exec"
	"runtime"
	"syscall"
	"time"

	"golang.org/x/sys/unix"
)

var Debug = func(string, ...interface{}) {}

// TraceRecord has information about a ptrace event.
type TraceRecord struct {
	EX     EventType
	Regs   unix.PtraceRegs
	Serial int
	Pid    int
	Err    error
	Errno  int
	Args   SyscallArguments
	Ret    [2]SyscallArgument
	Sysno  int
	Time   time.Duration
	Out    string
}

// Tracer has information to trace one process. It can be created by
// starting a command, or attaching. Attaching is not supported yet.
type Tracer struct {
	Pid     int
	Records chan *TraceRecord
	Count   int
	Raw     bool // Set by the user, it disables pretty printing
	Name    string
	Printer func(t *Tracer, r *TraceRecord)
	// We save the output from the previous Enter so Exit handling
	// can both use and adjust it.
	output []string
}

func New() (*Tracer, error) {
	return &Tracer{Pid: -1, Records: make(chan *TraceRecord), Printer: SysCall}, nil
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
	t.Name = fmt.Sprintf("%s(%d)", c.Args[0], t.Pid)
	t.Run()
}

func (t *Tracer) Step(e EventType) error {
	if err := unix.PtraceSyscall(t.Pid, 0); err != nil {
		r := &TraceRecord{Serial: t.Count, EX: e, Pid: t.Pid}
		Debug("ptracesyscall for %d gets %v", t.Pid, err)
		r.Err = err
		t.Records <- r
		return err
	}
	if w, err := Wait(t.Pid); err != nil {
		r := &TraceRecord{Serial: t.Count, EX: e, Pid: t.Pid}
		Debug("wait4 for %d gets %v, %v", t.Pid, w, err)
		r.Err = err
		t.Records <- r
		return err
	}
	return nil
}

func (t *Tracer) Run() error {
	var tm time.Time
	var a SyscallArguments
	var sysno = syscall.SYS_EXECVE
	for {
		t.Count++
		x := &TraceRecord{Serial: t.Count, EX: Exit, Pid: t.Pid, Args: a}
		if err := unix.PtraceGetRegs(t.Pid, &x.Regs); err != nil {
			Debug("ptracegetregs for %d gets %v", t.Pid, err)
			x.Err = err
			t.Records <- x
			break
		}
		x.FillArgs()
		x.FillRet()
		x.Sysno = sysno
		x.Time = time.Since(tm)

		// in the original version, we did the various probes of memory
		// in the strace command. Oops: on Unix, you can't look at process
		// memory when it's not stopped (that's not true on all systems,
		// but it's the Unix model). So we do this potentially expensive step
		// here, not knowing if we will not need it.
		if !t.Raw {
			SysCall(t, x)
		}
		t.Records <- x
		if err := t.Step(Enter); err != nil {
			return err
		}
		e := &TraceRecord{Serial: t.Count, EX: Enter, Pid: t.Pid}
		if err := unix.PtraceGetRegs(t.Pid, &e.Regs); err != nil {
			Debug("ptracegetregs for %d gets %v", t.Pid, err)
			e.Err = err
			t.Records <- e
			break
		}
		tm = time.Now()
		e.FillArgs()
		if !t.Raw {
			SysCall(t, e)
		}
		a = e.Args
		sysno = e.Sysno
		t.Records <- e

		if err := t.Step(Exit); err != nil {
			return err
		}

		t.Count++
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

// String is a stringer for TraceRecords
// TODO: stringer for Regs.
func (t *TraceRecord) String() string {
	pre := fmt.Sprintf("%s %d#%d:", t.EX, t.Pid, t.Serial)
	if t.Err != nil {
		return fmt.Sprintf("%s(%v)", pre, t.Err)
	}
	return fmt.Sprintf("%s %v", pre, t.Regs)
}

type ProcIO struct {
	pid   int
	addr  uintptr
	bytes int
}

func NewProcReader(pid int, addr uintptr) io.Reader {
	return &ProcIO{pid: pid, addr: addr}
}

func (p *ProcIO) Read(b []byte) (int, error) {
	n, err := unix.PtracePeekData(p.pid, p.addr, b)
	if err != nil {
		return n, err
	}
	p.addr += uintptr(n)
	p.bytes += n
	return n, nil
}

func NewProcWriter(pid int, addr uintptr) io.Writer {
	return &ProcIO{pid: pid, addr: addr}
}

func (p *ProcIO) Write(b []byte) (int, error) {
	n, err := unix.PtracePokeData(p.pid, p.addr, b)
	if err != nil {
		return n, err
	}
	p.addr += uintptr(n)
	p.bytes += n
	return n, nil
}

func (t *Tracer) Read(addr Addr, v interface{}) (int, error) {
	p := NewProcReader(t.Pid, uintptr(addr))
	err := binary.Read(p, binary.LittleEndian, v)
	return p.(*ProcIO).bytes, err
}

func (t *Tracer) ReadString(addr Addr, max int) (string, error) {
	if addr == 0 {
		return "<nil>", nil
	}
	var s string
	var b [1]byte
	for len(s) < max {
		if _, err := t.Read(addr, b[:]); err != nil {
			return "", err
		}
		if b[0] == 0 {
			break
		}
		s = s + string(b[:])
		addr++
	}
	return s, nil
}

func (t *Tracer) ReadStringVector(addr Addr, maxsize, maxno int) ([]string, error) {
	var v []Addr
	if addr == 0 {
		return []string{}, nil
	}

	fmt.Printf("read vec at %#x", addr)
	// Read in a maximum of maxno addresses
	for len(v) < maxno {
		var a uint64
		n, err := t.Read(addr, &a)
		if err != nil {
			fmt.Printf("Could not read vec elemtn at %v", addr)
			return nil, err
		}
		if a == 0 {
			break
		}
		addr += Addr(n)
		v = append(v, Addr(a))
	}
	fmt.Printf("Read %v", v)
	var vs []string
	for _, a := range v {
		s, err := t.ReadString(a, maxsize)
		if err != nil {
			fmt.Printf("Could not read string at %v", a)
			return vs, err
		}
		vs = append(vs, s)
	}
	return vs, nil
}

func (t *Tracer) Write(addr Addr, v interface{}) (int, error) {
	p := NewProcWriter(t.Pid, uintptr(addr))
	err := binary.Write(p, binary.LittleEndian, v)
	return p.(*ProcIO).bytes, err
}

func CaptureAddress(t *Tracer, addr Addr, addrlen uint32) ([]byte, error) {
	b := make([]byte, addrlen)
	if _, err := t.Read(addr, b); err != nil {
		return nil, err
	}
	return b, nil
}
