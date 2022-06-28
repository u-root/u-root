// Copyright 2022 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package syscallfilter supports filtering child process events
// by system call or strace event.  Strace events include
// SyscallEnter, SyscallExit, SignalExit, Exit, SignalStop, and
// NewChild.  System call events are for an explicit system call, and
// can be a regexp. They follow the naming pattern of the u-root
// strace package (which came from gvisor, which came from the plan 9
// "ratrace" command, described in
// https://www.osti.gov/biblio/1028390), i.e.  Esyscallname means
// system call entry; Xsyscallname means system call exit.  Thus you
// can, if you wish, allow the system call to run but force an error
// return by specifying Xsyscall, or force the entry to fail
// immediately by specifying Esyscall.
//
// For each event, an action can be specified.  These actions can in
// principle be set up at any time, although, currently, the package
// only supports setting them up before the process is started.
//
// Actions are specified as a triple: event,action,value.
// For example, if we wish to block any kind of fork, we would specify the
// action:
//
//	NewChild,error,-1
//
// which would cause a -1 to be returned from any action that causes a
// new child to be created.  If we wish to block all reads, we
// might specify
//
// E.*read,error,-1
package syscallfilter
