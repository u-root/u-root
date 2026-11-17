// Copyright 2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package bbmain is the package imported by all rewritten busybox
// command-packages to register themselves.
package bbmain

import (
	"errors"
	"fmt"
	"os"
	"sort"
	"sync/atomic"
	// There MUST NOT be any other dependencies here.
	//
	// It is preferred to copy minimal code necessary into this file, as
	// dependency management for this main file is... hard.
)

// ErrNotRegistered is returned by Run if the given command is not registered.
var ErrNotRegistered = errors.New("command is not present in busybox")

// Noop is a noop function.
var Noop = func() {}

// ListCmds returns all supported commands.
func ListCmds() []string {
	var cmds []string
	for c := range bbCmds {
		cmds = append(cmds, c)
	}
	sort.Strings(cmds)
	return cmds
}

type bbCmd struct {
	init, main func()
}

var bbCmds = map[string]bbCmd{}

var defaultCmd *bbCmd

// lazyInit reports whether a library's init has been deferred to the command
// that needs it, rather than run before main by the go tool. This file is
// written verbatim into the generated tree, so a busybox built that way arms
// it from a generated file beside this one, in this same package. It is false
// in every busybox built so far.
var lazyInit bool

// dispatched records which command Run has already dispatched, so that a
// second can be refused by name. It is only kept while lazyInit is set, which
// is the only time anything reads it.
//
// It is claimed with a compare-and-swap, because the refusal has to hold
// against two goroutines dispatching at once. That is the shape a plain
// read-then-write lets through: both see nothing dispatched, both proceed,
// and both run the deferred initialization, which is the thing being
// prevented. Nothing in a generated busybox does it -- main calls Run from
// one goroutine -- but a command may call bbmain.Run itself, and a guard is
// worth no more than the claim it makes.
//
// One variable rather than a flag beside a name, because the loser reads the
// winner's name to report it and two variables cannot be claimed together.
//
// Deferred initialization is not made safe to run concurrently, and
// deliberately. ADR 0003 says one command runs per process, so the answer to
// a second one is to refuse it; serialising it would be a different decision,
// and would leave the second command looking at state the first had changed.
var dispatched atomic.Pointer[string]

// Register registers an init and main function for name.
func Register(name string, init, main func()) {
	if _, ok := bbCmds[name]; ok {
		panic(fmt.Sprintf("cannot register two commands with name %q", name))
	}
	bbCmds[name] = bbCmd{
		init: init,
		main: main,
	}
}

// RegisterDefault registers a default init and main function.
func RegisterDefault(init, main func()) {
	defaultCmd = &bbCmd{
		init: init,
		main: main,
	}
}

// Run runs the command with the given name.
//
// If the command's main exits without calling os.Exit, Run will exit with exit
// code 0.
func Run(name string) error {
	var cmd *bbCmd
	if c, ok := bbCmds[name]; ok {
		cmd = &c
	} else if defaultCmd != nil {
		cmd = defaultCmd
	} else {
		return fmt.Errorf("%w: %s", ErrNotRegistered, name)
	}
	// Exactly one command runs per busybox process. The refusal sits here,
	// after a command has been selected, because the generated main program
	// calls Run twice on purpose: when argv[0] names no registered command it
	// retries with argv[1]. That first call returns above having run nothing,
	// so a refusal at the entry to Run would break dispatch itself.
	//
	// It is only load-bearing once a library's init is deferred to the command
	// that needs it. While every package is initialized before main, a second
	// command runs exactly the init Go would have run for it, so there is
	// nothing to refuse and the gate keeps this out of the way.
	if lazyInit {
		if !dispatched.CompareAndSwap(nil, &name) {
			panic(fmt.Sprintf("busybox dispatched %q and cannot also dispatch %q: one command runs per process", *dispatched.Load(), name))
		}
	}
	cmd.init()
	cmd.main()
	os.Exit(0)
	// Unreachable.
	return nil
}
