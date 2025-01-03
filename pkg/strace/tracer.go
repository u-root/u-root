// Copyright 2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build (linux && arm64) || (linux && amd64) || (linux && riscv64)

// Package strace traces Linux process events.
//
// An straced process will emit events for syscalls, signals, exits, and new
// children.
package strace

import (
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"os/exec"
	"runtime"
	"sync/atomic"
	"syscall"
	"time"

	"golang.org/x/sys/unix"
)

func wait(pid int) (int, unix.WaitStatus, error) {
	var w unix.WaitStatus
	pid, err := unix.Wait4(pid, &w, 0, nil)
	return pid, w, err
}

// TraceError is returned when something failed on a specific process.
type TraceError struct {
	// PID is the process ID associated with the error.
	PID int
	Err error
}

func (t *TraceError) Error() string {
	return fmt.Sprintf("trace error on pid %d: %v", t.PID, t.Err)
}

// SyscallEvent is populated for both SyscallEnter and SyscallExit event types.
type SyscallEvent struct {
	// Regs are the process's registers as they were when the event was
	// recorded.
	Regs unix.PtraceRegs

	// Sysno is the syscall number.
	Sysno int

	// Args are the arguments to the syscall.
	Args SyscallArguments

	// Ret is the return value of the syscall. Only populated on
	// SyscallExit.
	Ret [2]SyscallArgument

	// Errno is an errno, if there was on in Ret. Only populated on
	// SyscallExit.
	Errno unix.Errno

	// Duration is the duration from enter to exit for this particular
	// syscall. Only populated on SyscallExit.
	Duration time.Duration
}

// SignalEvent is a signal that was delivered to the process.
type SignalEvent struct {
	// Signal is the signal number.
	Signal unix.Signal

	// TODO: Add other siginfo_t stuff
}

// ExitEvent is emitted when the process exits regularly using exit_group(2).
type ExitEvent struct {
	// WaitStatus is the exit status.
	WaitStatus unix.WaitStatus
}

// NewChildEvent is emitted when a clone/fork/vfork syscall is done.
type NewChildEvent struct {
	PID int
}

// TraceRecord has information about a process event.
type TraceRecord struct {
	PID   int
	Time  time.Time
	Event EventType

	// Poor man's union. One of the following five will be populated
	// depending on the Event.

	Syscall    *SyscallEvent
	SignalExit *SignalEvent
	SignalStop *SignalEvent
	Exit       *ExitEvent
	NewChild   *NewChildEvent
}

// process is a Linux thread.
type process struct {
	pid int

	// ptrace does not tell you whether a syscall-stop is a
	// syscall-enter-stop or syscall-exit-stop. You gotta keep track of
	// that shit your own self.
	lastSyscallStop *TraceRecord
	SecComp         atomic.Bool
}

// Name implements Task.Name.
func (p *process) Name() string {
	return fmt.Sprintf("[pid %d]", p.pid)
}

// Read reads from the process at Addr to the interface{}
// and returns a byte count and error.
func (p *process) Read(addr Addr, v interface{}) (int, error) {
	r := newProcReader(p.pid, uintptr(addr))
	err := binary.Read(r, binary.NativeEndian, v)
	return r.bytes, err
}

func (p *process) cont(signal unix.Signal) error {
	// Event has been processed. Restart 'em.
	if p.SecComp.Load() {
		// If seccomp is enabled, continue the process without stopping at each syscall.
		if err := unix.PtraceCont(p.pid, int(signal)); err != nil {
			return os.NewSyscallError("ptrace(PTRACE_SYSCALL)", fmt.Errorf("on pid %d: %w", p.pid, err))
		}
		return nil
	}
	// If seccomp is not enabled, continue the process and stop at each syscall.
	if err := unix.PtraceSyscall(p.pid, int(signal)); err != nil {
		return os.NewSyscallError("ptrace(PTRACE_SYSCALL)", fmt.Errorf("on pid %d: %w", p.pid, err))
	}
	return nil
}

type tracer struct {
	processes map[int]*process
	callback  []EventCallback
}

func (t *tracer) call(p *process, rec *TraceRecord) error {
	for _, c := range t.callback {
		if err := c(p, rec); err != nil {
			return err
		}
	}
	return nil
}

var traceActive uint32

// Trace traces `c` and any children c clones.
//
// Only one trace can be active per process.
//
// recordCallback is called every time a process event happens with the process
// in a stopped state.
func Trace(c *exec.Cmd, recordCallback ...EventCallback) error {
	return New(c, false, recordCallback...)
}

// New traces `c` and any children c clones with the option to enable seccomp.
func New(c *exec.Cmd, secComp bool, recordCallback ...EventCallback) error {
	if !atomic.CompareAndSwapUint32(&traceActive, 0, 1) {
		return fmt.Errorf("a process trace is already active in this process")
	}
	defer func() {
		atomic.StoreUint32(&traceActive, 0)
	}()

	if c.SysProcAttr == nil {
		c.SysProcAttr = &syscall.SysProcAttr{}
	}
	c.SysProcAttr.Ptrace = true

	// Because the go runtime forks traced processes with PTRACE_TRACEME
	// we need to maintain the parent-child relationship for ptrace to work.
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	if err := c.Start(); err != nil {
		return err
	}

	tracer := &tracer{
		processes: make(map[int]*process),
		callback:  recordCallback,
	}

	// Start will fork, set PTRACE_TRACEME, and then execve. Once that
	// happens, we should be stopped at the execve "exit". This wait will
	// return at that exit point.
	//
	// The new task image has been loaded at this point, with us just about
	// to jump into _start.
	//
	// It'd make sense to assume, but this stop is NOT a syscall-exit-stop
	// of the execve. It is a signal-stop triggered at the end of execve,
	// within the confines of the new task image.  This means the execve
	// syscall args are not in their registers, and we can't print the
	// exit.
	//
	// NOTE(chrisko): we could make it such that we can read the args of
	// the execve. If we were to signal ourselves between PTRACE_TRACEME
	// and execve, we'd stop before the execve and catch execve as a
	// syscall-stop after. To do so, we have 3 options: (1) write a copy of
	// stdlib exec.Cmd.Start/os.StartProcess with the change, or (2)
	// upstreaming a change that would make it into the next Go version, or
	// (3) use something other than *exec.Cmd as the API.
	//
	// A copy of the StartProcess logic would be tedious, an upstream
	// change would take a while to get into Go, and we want this API to be
	// easily usable. I think it's ok to sacrifice the execve for now.
	if _, ws, err := wait(c.Process.Pid); err != nil {
		return err
	} else if ws.TrapCause() != 0 {
		return fmt.Errorf("wait(pid=%d): got %v, want stopped process", c.Process.Pid, ws)
	}
	tracer.addProcess(c.Process.Pid, SyscallExit, secComp)

	if err := unix.PtraceSetOptions(c.Process.Pid,
		// Tells ptrace to generate a SIGTRAP signal immediately before a new program is executed with the execve system call.
		unix.PTRACE_O_TRACEEXEC|
			// Tells ptrace to generate a SIGTRAP signal for seccomp events.
			unix.PTRACE_O_TRACESECCOMP|
			// Make it easy to distinguish syscall-stops from other SIGTRAPS.
			unix.PTRACE_O_TRACESYSGOOD|
			// Kill tracee if tracer exits.
			unix.PTRACE_O_EXITKILL|
			// Automatically trace fork(2)'d, clone(2)'d, and vfork(2)'d children.
			unix.PTRACE_O_TRACECLONE|unix.PTRACE_O_TRACEFORK|unix.PTRACE_O_TRACEVFORK); err != nil {
		return &TraceError{
			PID: c.Process.Pid,
			Err: os.NewSyscallError("ptrace(PTRACE_SETOPTIONS)", err),
		}
	}

	// Start the process back up.
	if err := unix.PtraceSyscall(c.Process.Pid, 0); err != nil {
		return &TraceError{
			PID: c.Process.Pid,
			Err: fmt.Errorf("failed to resume: %w", err),
		}
	}

	return tracer.runLoop()
}

func (t *tracer) addProcess(pid int, event EventType, secComp bool) {
	t.processes[pid] = &process{
		pid:     pid,
		SecComp: atomic.Bool{},
		lastSyscallStop: &TraceRecord{
			Event: event,
			Time:  time.Now(),
		},
	}
	if secComp {
		t.processes[pid].SecComp.Store(true) // Use the atomic method to set the value
	}
}

func (t *TraceRecord) syscallStop(p *process) error {
	t.Syscall = &SyscallEvent{}

	if err := unix.PtraceGetRegs(p.pid, &t.Syscall.Regs); err != nil {
		return &TraceError{
			PID: p.pid,
			Err: os.NewSyscallError("ptrace(PTRACE_GETREGS)", err),
		}
	}

	t.Syscall.FillArgs()

	// TODO: the ptrace man page mentions that seccomp can inject a
	// syscall-exit-stop without a preceding syscall-enter-stop. Detect
	// that here, however you'd detect it...
	if p.lastSyscallStop.Event == SyscallEnter {
		t.Event = SyscallExit
		t.Syscall.FillRet()
		t.Syscall.Duration = time.Since(p.lastSyscallStop.Time)
	} else {
		t.Event = SyscallEnter
	}
	p.lastSyscallStop = t
	return nil
}

func (t *tracer) runLoop() error {
	for {
		// TODO: we cannot have any other children. I'm not sure this
		// is actually solvable: if we used a session or process group,
		// a tracee process's usage of them would mess up our accounting.
		//
		// If we just ignored wait's of processes that we're not
		// tracing, we'll be messing up other stuff in this program
		// waiting on those.
		//
		// To actually encapsulate this library in a packge, we could
		// do one of two things:
		//
		//   1) fork from the parent in order to be able to trace
		//      children correctly. Then, a user of this library could
		//      actually independently trace two different processes.
		//      I don't know if that's worth doing.
		//   2) have one goroutine per process, and call wait4
		//      individually on each process we expect. We gotta check
		//      if each has to be tied to an OS thread or not.
		//
		// The latter option seems much nicer.
		pid, status, err := wait(-1)
		if err == unix.ECHILD {
			// All our children are gone.
			return nil
		} else if err != nil {
			return os.NewSyscallError("wait4", err)
		}

		// Which process was stopped?
		p, ok := t.processes[pid]
		if !ok {
			continue
		}

		rec := &TraceRecord{
			PID:  p.pid,
			Time: time.Now(),
		}

		var injectSignal unix.Signal
		if status.Exited() {
			rec.Event = Exit
			rec.Exit = &ExitEvent{
				WaitStatus: status,
			}
		} else if status.Signaled() {
			rec.Event = SignalExit
			rec.SignalExit = &SignalEvent{
				Signal: status.Signal(),
			}
		} else if status.Stopped() {
			// Ptrace stops kinds.
			switch signal := status.StopSignal(); signal {
			// Syscall-stop.
			//
			// Setting PTRACE_O_TRACESYSGOOD means StopSignal ==
			// SIGTRAP|0x80 (0x85) for syscall-stops.
			//
			// It allows us to distinguish syscall-stops from regular
			// SIGTRAPs (e.g. sent by tkill(2)).
			case syscall.SIGTRAP | 0x80:
				if err := rec.syscallStop(p); err != nil {
					return err
				}

			// Group-stop, but also a special stop: first stop after
			// fork/clone/vforking a new task.
			//
			// TODO: is that different than a group-stop, or the same?
			case syscall.SIGSTOP:
				// TODO: have a list of expected children SIGSTOPs, and
				// make events only for all the unexpected ones.
				fallthrough

			// Group-stop.
			//
			// TODO: do something.
			case syscall.SIGTSTP, syscall.SIGTTOU, syscall.SIGTTIN:
				rec.Event = SignalStop
				injectSignal = signal
				rec.SignalStop = &SignalEvent{
					Signal: signal,
				}

				// TODO: Do we have to use PTRACE_LISTEN to
				// restart the task in order to keep the task
				// in stopped state, as expected by whomever
				// sent the stop signal?

			// Either a regular signal-delivery-stop, or a PTRACE_EVENT stop.
			case syscall.SIGTRAP:
				switch tc := status.TrapCause(); tc {
				case unix.PTRACE_EVENT_SECCOMP:
					// Handle seccomp event by continuing the syscall.
					if err := syscall.PtraceSyscall(pid, 0); err != nil {
						return os.NewSyscallError("ptrace(PTRACE_SYSCALL)", fmt.Errorf("on pid %d: %w", p.pid, err))
					}
					continue
				// This is a PTRACE_EVENT stop.
				case unix.PTRACE_EVENT_CLONE, unix.PTRACE_EVENT_FORK, unix.PTRACE_EVENT_VFORK:
					childPID, err := unix.PtraceGetEventMsg(pid)
					if err != nil {
						return &TraceError{
							PID: pid,
							Err: os.NewSyscallError("ptrace(PTRACE_GETEVENTMSG)", err),
						}
					}
					// The first event will be an Enter syscall, so
					// set the last event to an exit.
					t.addProcess(int(childPID), SyscallExit, p.SecComp.Load())

					rec.Event = NewChild
					rec.NewChild = &NewChildEvent{
						PID: int(childPID),
					}

				// Regular signal-delivery-stop.
				default:
					rec.Event = SignalStop
					rec.SignalStop = &SignalEvent{
						Signal: signal,
					}
					injectSignal = signal
				}

			// Signal-delivery-stop.
			default:
				rec.Event = SignalStop
				rec.SignalStop = &SignalEvent{
					Signal: signal,
				}
				injectSignal = signal
			}
		} else {
			rec.Event = Unknown
		}

		if err := t.call(p, rec); err != nil {
			return err
		}

		if rec.Event == SignalExit || rec.Event == Exit {
			delete(t.processes, pid)
			continue
		}

		if err := p.cont(injectSignal); err != nil {
			return err
		}
	}
}

// EventCallback is a function called on each event while the subject process
// is stopped.
type EventCallback func(t Task, record *TraceRecord) error

// RecordTraces sends each event on c.
func RecordTraces(c chan<- *TraceRecord) EventCallback {
	return func(t Task, record *TraceRecord) error {
		c <- record
		return nil
	}
}

func signalString(s unix.Signal) string {
	if 0 <= s && int(s) < len(signals) {
		return fmt.Sprintf("%s (%d)", signals[s], int(s))
	}
	return fmt.Sprintf("signal %d", int(s))
}

// PrintTraces prints every trace event to w.
func PrintTraces(w io.Writer) EventCallback {
	return func(t Task, record *TraceRecord) error {
		switch record.Event {
		case SyscallEnter:
			fmt.Fprintln(w, SysCallEnter(t, record.Syscall))
		case SyscallExit:
			fmt.Fprintln(w, SysCallExit(t, record.Syscall))
		case SignalExit:
			fmt.Fprintf(w, "PID %d exited from signal %s\n", record.PID, signalString(record.SignalExit.Signal))
		case Exit:
			fmt.Fprintf(w, "PID %d exited from exit status %d (code = %d)\n", record.PID, record.Exit.WaitStatus, record.Exit.WaitStatus.ExitStatus())
		case SignalStop:
			fmt.Fprintf(w, "PID %d got signal %s\n", record.PID, signalString(record.SignalStop.Signal))
		case NewChild:
			fmt.Fprintf(w, "PID %d spawned new child %d\n", record.PID, record.NewChild.PID)
		}
		return nil
	}
}

// Strace traces and prints process events for `c` and its children to `out`.
func Strace(c *exec.Cmd, out io.Writer) error {
	return Trace(c, PrintTraces(out))
}

// EventType describes a process event.
type EventType int

const (
	// Unknown is for events we do not know how to interpret.
	Unknown EventType = 0x0

	// SyscallEnter is the event for a process calling a syscall.  Event
	// Args will contain the arguments sent by the userspace process.
	//
	// ptrace calls this a syscall-enter-stop.
	SyscallEnter EventType = 0x2

	// SyscallExit is the event for the kernel returning a syscall. Args
	// will contain the arguments as returned by the kernel.
	//
	// ptrace calls this a syscall-exit-stop.
	SyscallExit EventType = 0x3

	// SignalExit means the process has been terminated by a signal.
	SignalExit EventType = 0x4

	// Exit means the process has exited with an exit code.
	Exit EventType = 0x5

	// SignalStop means the process was stopped by a signal.
	//
	// ptrace calls this a signal-delivery-stop.
	SignalStop EventType = 0x6

	// NewChild means the process created a new child thread or child
	// process via fork, clone, or vfork.
	//
	// ptrace calls this a PTRACE_EVENT_(FORK|CLONE|VFORK).
	NewChild EventType = 0x7
)

// A procIO is used to implement io.Reader and io.Writer.
// it contains a pid, which is unchanging; and an
// addr and byte count which change as IO proceeds.
type procIO struct {
	pid   int
	addr  uintptr
	bytes int
}

// newProcReader returns an io.Reader for a procIO.
func newProcReader(pid int, addr uintptr) *procIO {
	return &procIO{pid: pid, addr: addr}
}

// Read implements io.Read for a procIO.
func (p *procIO) Read(b []byte) (int, error) {
	n, err := unix.PtracePeekData(p.pid, p.addr, b)
	if err != nil {
		return n, err
	}
	p.addr += uintptr(n)
	p.bytes += n
	return n, nil
}

// ReadString reads a null-terminated string from the process
// at Addr and any errors.
func ReadString(t Task, addr Addr, maximum int) (string, error) {
	if addr == 0 {
		return "<nil>", nil
	}
	var s string
	var b [1]byte
	for len(s) < maximum {
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
func ReadStringVector(t Task, addr Addr, maxsize, maxno int) ([]string, error) {
	var v []Addr
	if addr == 0 {
		return []string{}, nil
	}

	// Read in a maximum of maxno addresses
	for len(v) < maxno {
		var a uint64
		n, err := t.Read(addr, &a)
		if err != nil {
			return nil, fmt.Errorf("could not read vector element at %#x: %w", addr, err)
		}
		if a == 0 {
			break
		}
		addr += Addr(n)
		v = append(v, Addr(a))
	}
	var vs []string
	for _, a := range v {
		s, err := ReadString(t, a, maxsize)
		if err != nil {
			return vs, fmt.Errorf("could not read string at %#x: %w", a, err)
		}
		vs = append(vs, s)
	}
	return vs, nil
}

// CaptureAddress pulls a socket address from the process as a byte slice.
// It returns any errors.
func CaptureAddress(t Task, addr Addr, addrlen uint32) ([]byte, error) {
	b := make([]byte, addrlen)
	if _, err := t.Read(addr, b); err != nil {
		return nil, err
	}
	return b, nil
}
