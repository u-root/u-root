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

// Debug is a do-nothing function which can be replaced by, e.g., log.Printf
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
	EX      EventType
	Records chan *TraceRecord
	Count   int
	Raw     bool // Set by the user, it disables pretty printing
	Name    string
	Printer func(t *Tracer, r *TraceRecord)
	Last    *TraceRecord
	// We save the output from the previous Enter so Exit handling
	// can both use and adjust it.
	output []string
}

// New returns a new Tracer.
func New() (*Tracer, error) {
	return &Tracer{Pid: -1, Records: make(chan *TraceRecord, 1), Printer: SysCall}, nil
}

// RunTracerFromCommand runs a Tracer given an exec.Cmd.
// It locks itself down with LockOSThread and will unlock itself
// when it returns, after the command and all its children exit.
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
	t.EX = Exit
	Run(t)
}

// NewTracerChild creates a tracer from a tracer.
func NewTracerChild(pid int) (*Tracer, error) {
	nt, err := New()
	if err != nil {
		return nil, err
	}
	nt.Pid = pid
	nt.Name = fmt.Sprintf("%d", pid)
	nt.EX = Exit
	return nt, nil
}

// Step steps a Tracer by issuing a PtraceSyscall to it and then doing a Wait.
// Note that Step waits for any child to return, not just the one we are stepping.
func (t *Tracer) Step(e EventType) (int, error) {
	Debug("Step %d", t.Pid)
	if err := unix.PtraceSyscall(t.Pid, 0); err != nil {
		r := &TraceRecord{Serial: t.Count, EX: e, Pid: t.Pid}
		Debug("ptracesyscall for %d gets %v", t.Pid, err)
		r.Err = fmt.Errorf("unix.PtraceSyscall: %d: %s: %v", t.Pid, t.Name, err)
		t.Records <- r
		return -1, r.Err
	}
	Debug("Stepped %d, now Wait", t.Pid)
	pid, w, err := Wait(-1)
	Debug("Wait returns (%d, %v, %v)", pid, w, err)
	if err != nil {
		r := &TraceRecord{Serial: t.Count, EX: e, Pid: t.Pid}
		r.Err = fmt.Errorf("unix.Wait: %d: %s: %v, %v", t.Pid, t.Name, w, err)
		Debug("wait4 for %d gets %v, %v", t.Pid, w, err)
		t.Records <- r
		return -1, r.Err
	}
	Debug("Step %d: back from wait", pid)
	return pid, nil
}

// Run runs a set of processes as defined by a Tracer. Because of Unix restrictions
// around which processes which can trace other processes, Run gets a tad involved.
// It is implemented as a simple loop, driving events via ptrace commands to processes;
// and responding to events returned by a Wait.
// It has to handle a few events specially:
// o if a wait fails, the process has exited, and must no longer be commanded
//   this is indicated by a wait followed by an error on PtraceGetRegs
// o if a process forks successfully, we must add it to our set of traced processes.
//   We attach that process, wait for it, then issue a ptrace system call command
//   to it. We don't use the Linux SEIZE command as we can do this in a more generic
//   Unix way.
// We create a map of our traced processes and run until it is empty.
// The initial value of the map is just the one process we start with.
func Run(root *Tracer) error {
	var nextEX EventType
	var procs = map[int]*Tracer{
		root.Pid: root,
	}
	Debug("procs %v", procs)
	var tm time.Time
	var a SyscallArguments
	var sysno = syscall.SYS_EXECVE
	Debug("Run %v", root.Pid)
	pid := root.Pid
	var err error
	var count int
	for len(procs) > 0 {
		t := procs[pid]
		t.Count++
		count++
		Debug("Get regs for %d", pid)
		x := &TraceRecord{Serial: t.Count, EX: t.EX, Pid: pid, Args: a}
		if err := unix.PtraceGetRegs(pid, &x.Regs); err != nil {
			Debug("ptracegetregs for %d gets %v", pid, err)
			x.Err = fmt.Errorf("ptracegetregs for %d gets %v", pid, err)
			t.Records <- x
			delete(procs, pid)
			pid, _, _ = Wait(-1)
			continue
		}
		Debug("GOT regs for %d", pid)
		x.FillArgs()
		if t.EX == Exit {
			x.FillRet()
			x.Sysno = sysno
			x.Time = time.Since(tm)
			nextEX = Enter
		} else {
			tm = time.Now()
			x.FillArgs()
			a = x.Args
			sysno = x.Sysno
			t.Last = x
			nextEX = Exit
		}

		if !t.Raw {
			SysCall(t, x)
		}
		Debug("Push %v", x)
		t.Records <- x
		// Was there a clone? Capture the child. Don't forget the child has an exit
		// record for the clone too, so don't get confused.
		p := int(x.Ret[0].Int())
		Debug("Check for new pid: tracer pid %d, ret %d", t.Pid, p)
		if x.Sysno == unix.SYS_CLONE && x.EX == Exit && p > 0 && p != t.Pid {
			nt, err := NewTracerChild(int(x.Ret[0].Int()))
			if err != nil {
				Debug("Setting up child: %v", err)
			} else {
				nt.Records = t.Records
				Debug("New child: %v", nt)
				// The result of the attach gets picked up by the wait()
				if err := unix.PtraceAttach(nt.Pid); err != nil {
					r := &TraceRecord{Serial: nt.Count, EX: Enter, Pid: nt.Pid}
					Debug("RunTracerChild: attach for %d gets %v", nt.Pid, err)
					nt.Records <- r
				}
				pid, w, err := Wait(nt.Pid)
				Debug("Wait returns (%d, %v, %v)", pid, w, err)
				if err != nil || pid != nt.Pid {
					r := &TraceRecord{Serial: nt.Count, EX: Exit, Pid: nt.Pid}
					r.Err = fmt.Errorf("unix.Wait: %d: %s: %v, %v", nt.Pid, nt.Name, w, err)
					Debug("wait4 for %d gets %v, %v", nt.Pid, w, err)
					t.Records <- r
				}
				Debug("Step %d", nt.Pid)
				if err := unix.PtraceSyscall(nt.Pid, 0); err != nil {
					r := &TraceRecord{Serial: nt.Count, EX: Exit, Pid: nt.Pid}
					Debug("ptracesyscall for %d gets %v", nt.Pid, err)
					r.Err = fmt.Errorf("unix.PtraceSyscall: %d: %s: %v", nt.Pid, nt.Name, err)
					t.Records <- r
				}
				procs[nt.Pid] = nt
			}
		}

		Debug("Step after exit")
		t.EX = nextEX
		if pid, err = t.Step(t.EX); err != nil {
			return err
		}
	}

	Debug("Pushed %d records", count)
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

// A ProcIO is used to implement io.Reader and io.Writer.
// it contains a pid, which is unchanging; and an
// addr and byte count which change as IO proceeds.
type ProcIO struct {
	pid   int
	addr  uintptr
	bytes int
}

// NewProcReader returns an io.Reader for a ProcIO.
func NewProcReader(pid int, addr uintptr) io.Reader {
	return &ProcIO{pid: pid, addr: addr}
}

// Read implements io.Read for a ProcIO.
func (p *ProcIO) Read(b []byte) (int, error) {
	n, err := unix.PtracePeekData(p.pid, p.addr, b)
	if err != nil {
		return n, err
	}
	p.addr += uintptr(n)
	p.bytes += n
	return n, nil
}

// NewProcWriter returns an io.Writer for a ProcIO.
func NewProcWriter(pid int, addr uintptr) io.Writer {
	return &ProcIO{pid: pid, addr: addr}
}

// Write implements io.Write for a ProcIO.
func (p *ProcIO) Write(b []byte) (int, error) {
	n, err := unix.PtracePokeData(p.pid, p.addr, b)
	if err != nil {
		return n, err
	}
	p.addr += uintptr(n)
	p.bytes += n
	return n, nil
}

// Read reads from the process at Addr to the interface{}
// and returns a byte count and error.
func (t *Tracer) Read(addr Addr, v interface{}) (int, error) {
	p := NewProcReader(t.Pid, uintptr(addr))
	err := binary.Read(p, binary.LittleEndian, v)
	return p.(*ProcIO).bytes, err
}

// ReadString reads a null-terminated string from the process
// at Addr and any errors.
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

// ReadStringVector takes an address, max string size, and max number of string to read,
// and returns a string slice or error.
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

// Write writes to the process address sapce and returns a count and error.
func (t *Tracer) Write(addr Addr, v interface{}) (int, error) {
	p := NewProcWriter(t.Pid, uintptr(addr))
	err := binary.Write(p, binary.LittleEndian, v)
	return p.(*ProcIO).bytes, err
}

// CaptureAddress pulls a socket address from the process as a byte slice.
// It returns any errors.
func CaptureAddress(t *Tracer, addr Addr, addrlen uint32) ([]byte, error) {
	b := make([]byte, addrlen)
	if _, err := t.Read(addr, b); err != nil {
		return nil, err
	}
	return b, nil
}
