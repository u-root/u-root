// Copyright 2012-2022 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build (linux && arm64) || (linux && amd64) || (linux && riscv64)

package syscallfilter

import (
	"context"
	"fmt"
	"io"
	"log"
	"os/exec"
	"regexp"
	"strconv"
	"strings"

	"github.com/u-root/u-root/pkg/strace"
)

const cmdUsage = "Usage: strace [-o <outputfile>] <command> [args...]"

type event struct {
	Name   string
	Pat    string
	Action string
	Value  int
}

// Cmd contains a command and a filter.
// The filter is allowed to be empty.
type Cmd struct {
	*exec.Cmd
	// if Log is non-nil, actions take are written to it.
	Log io.Writer
	// events is a simple array.
	// It was a map in earlier versions, but:
	// o it is usually going to be short
	// o we want to support regexp matching (part of keeping it short and convenient!)
	// So we iterate it on every system call. The cost of doing the strace is so high that
	// iterating over this array is not significant.
	events []*event
	// The strace package has its limitations and we don't expect to use it for long.
	// It's not capable of returning a value to the process from a system call,
	// and it may not be appropriate to extend it in that way.
	// We keep this function here to allow us to kill the child on error.
	cancel func()
}

var eventNames = map[string]strace.EventType{
	"SyscallEnter": strace.SyscallEnter,
	"SyscallExit":  strace.SyscallExit,
	"SignalExit":   strace.SignalExit,
	"Exit":         strace.Exit,
	"SignalStop":   strace.SignalStop,
	"NewChild":     strace.NewChild,
}

var allActions = map[string]interface{}{
	// Error will end the strace, expeditiously, unless the value is 0
	"error": nil,
	// log will log the record
	"log": nil,
}

func findEvent(events []*event, n string) *event {
	for _, e := range events {
		m, err := regexp.MatchString(e.Pat, n)
		if err != nil {
			log.Fatalf("Can't happen: %q is bad: %v", e.Pat, err)
		}
		if m {
			return e
		}
	}
	return nil
}

// eventName returns an event name. There should never be an event name
// we do not known and, if we encounter one, we panic.
func eventName(r *strace.TraceRecord) string {
	// form up a reasonable name for a system call.
	// If there is no name, then it will be Exxxx or Xxxxx, where x
	// is the system call number as %04x.
	// Note that users can specify this: E0x0000, for example
	var sysname string

	switch r.Event {
	case strace.SyscallEnter, strace.SyscallExit:
		var err error
		if sysname, err = strace.ByNumber(uintptr(r.Syscall.Sysno)); err != nil {
			sysname = fmt.Sprintf("%04x", r.Syscall.Sysno)
		}
	}
	switch r.Event {
	case strace.SyscallEnter:
		return "E" + sysname
	case strace.SyscallExit:
		return "X" + sysname
	case strace.SignalExit:
		return "SignalExit"
	case strace.Exit:
		return "Exit"
	case strace.SignalStop:
		return "SignalStop"
	case strace.NewChild:
		return "NewChild"
	default:
		log.Panicf("Unknown event %#x from record %v", r.Event, r)
	}
	return ""
}

func (c *Cmd) handleEvent(t strace.Task, r *strace.TraceRecord, e []*event) error {
	// All attempts to use defer for printing got ... weird.
	var ret error
	n := eventName(r)
	act := findEvent(e, n)
	if act == nil {
		return nil
	}
	switch act.Action {
	case "error":
		ret = fmt.Errorf("%v", act.Value)
		if c.Log != nil {
			fmt.Fprintf(c.Log, "%v act %v error %v\n", n, act.Name, ret)
		}
		// And, since we can't really return an error to the process, but only to strace
		// and, since this is more about killing it than anything else ...
		// bye bye
		c.cancel()

	// Actually, for now, log is the same as error,0 but we are keeping this
	// name as we expect it to evolve.
	case "log":
		if c.Log != nil {
			fmt.Fprintf(c.Log, "%v act %v error %v\n", n, act.Name, ret)
		}
	}
	return ret
}

// Run implements cmd.Run, filtering strace events
// as created by AddActions. The slice can be empty, in which case the command
// runs as normal.
func (c *Cmd) Run() error {
	// This wait may or may not be needed, since the process
	// can end normally or be stopped by a filter. Hence,
	// we will not check for an error.
	defer c.Wait()
	return strace.Trace(c.Cmd, func(t strace.Task, r *strace.TraceRecord) error {
		return c.handleEvent(t, r, c.events)
	})
}

// AddActions creates an []event as defined by a possibly empty set of actions, and
// installs it in the Cmd.
// The action is a string, with comma-separated fields,
// designed to be convenient to be used on a command line.
//
// Actions have 3 fields.
//
// The first is a regexp to match against the system call name, e.g.,
// E.*imeofday. The first character, E or X,
// indicates Entry or eXit.
//
// The second parameter indications an action to take. Currently there
// are only two: error, and log.
//
// The third value indicates a value to return. It can be empty.
// For the error action, the third parameter
// indicates the error to return. For the log action, the third
// parameter is currently unused.
func (c *Cmd) AddActions(actions ...string) error {
	var events []*event
	for i, a := range actions {
		f := strings.Split(a, ",")
		if len(f) != 3 {
			return fmt.Errorf("%d of actions: %q needs to have 3 fields, has %d(%v)", i, a, len(f), f)
		}
		if _, ok := allActions[f[1]]; !ok {
			return fmt.Errorf("%q of actions %v: unknown action, not one of %q", f[0], f, allActions)
		}
		if check := findEvent(events, f[0]); check != nil {
			return fmt.Errorf("%q of actions %v: repeat action, already %q", f[0], f, check)
		}
		var value int
		if len(f[2]) > 0 {
			v, err := strconv.ParseInt(f[2], 0, 32)
			if err != nil {
				return err
			}
			value = int(v)
		}

		events = append(events, &event{Name: a, Pat: f[0], Action: f[1], Value: value})
	}
	c.events = events
	return nil
}

// Command creates a new Cmd, with a context, embedding an exec.Cmd, with an empty set of events.
func Command(name string, args ...string) *Cmd {
	ctx, cancel := context.WithCancel(context.Background())
	return &Cmd{Cmd: exec.CommandContext(ctx, name, args...), cancel: cancel}
}
