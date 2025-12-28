// Copyright 2025 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package builtin defines how packages can be used
// as command builtins in bare metal or other embedded environments.
//
// The common mechanism for using a command as a package is to use
// the gobusybox tool. That said, we have occasional requests to make
// individual commands available as packages; the gobusybox
// model, in some cases, brings in more than people want.
// In tightly space-constrained environments, such as the Pico2,
// every byte counts.
//
// This package defines a type, Cmd, and an interface,
// Runner. Packages must define their own Cmd type
// which may be as simple as:
//
//	type Cmd struct {
//	  *builtin.Cmd
//	  ...
//	  }
//
// If a package defines builtin.Cmd, it should make an
// assertion:
// var _ builtin.Runner = &Cmd{}
// For this assertion to work, a package must satisfy the
// builtin.Runner interface:
// func (c *Cmd) Run() error
//
// Packages must define a Command function, e.g.
// func Command(path string, args ...string) *Cmd
// The non-optional path arg should be thought of as the
// arg0 provided by the kernel to a process on exec(2).
// Among other things, this provides consistency with
// os/exec. That may seem a foolish consistency, but
// we have already had to fix one bug that
// resulted from a Command method which was not consistent
// with os/exec.Command().
//
// The most common pattern for a user is to call a function,
// Command, which returns a package defined type which
// must implement the Runner interface.
//
// With the returned struct, users can call the Run method.
//
// To see a usage, look at builtin_test.go in the forth package.
// It uses the forth package to run simple forth scripts.
package builtin
