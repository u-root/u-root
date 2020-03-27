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
	"path/filepath"
)

// ErrNotRegistered is returned by Run if the given command is not registered.
var ErrNotRegistered = errors.New("command not registered")

// Noop is a noop function.
var Noop = func() {}

// ListCmds lists bb commands and verifies symlinks.
// It is by convention called when the bb command is invoked directly.
// For every command, there should be a symlink in /bbin,
// and for every symlink, there should be a command.
// Occasionally, we have bugs that result in one of these
// being false. Just running bb is an easy way to tell if something
// in your image is messed up.
func ListCmds() {
	type known struct {
		name string
		bb   string
	}
	names := map[string]*known{}
	g, err := filepath.Glob("/bbin/*")
	if err != nil {
		fmt.Printf("bb: unable to enumerate /bbin")
	}

	// First step is to assemble a list of all possible
	// names, both from /bbin/* and our built in commands.
	for _, l := range g {
		if l == "/bbin/bb" {
			continue
		}
		b := filepath.Base(l)
		names[b] = &known{name: l}
	}
	for n := range bbCmds {
		if n == "bb" {
			continue
		}
		if c, ok := names[n]; ok {
			c.bb = n
			continue
		}
		names[n] = &known{bb: n}
	}
	// Now walk the array of structs.
	// We don't sort as we don't want the
	// footprint of bringing in the package.
	// If you want it sorted, bb | sort
	var hadError bool
	for c, k := range names {
		if len(k.name) == 0 || len(k.bb) == 0 {
			hadError = true
			fmt.Printf("%s:\t", c)
			if k.name == "" {
				fmt.Printf("NO SYMLINK\t")
			} else {
				fmt.Printf("%q\t", k.name)
			}
			if k.bb == "" {
				fmt.Printf("NO COMMAND\n")
			} else {
				fmt.Printf("%s\n", k.bb)
			}
		}
	}
	if hadError {
		fmt.Println("There is at least one problem. Known causes:")
		fmt.Println("At least two initrds -- one compiled in to the kernel, a second supplied by the bootloader.")
		fmt.Println("The initrd cpio was changed after creation or merged with another one.")
		fmt.Println("When the initrd was created, files were inserted into /bbin by mistake.")
		fmt.Println("Post boot, files were added to /bbin.")
	}
}

type bbCmd struct {
	init, main func()
}

var bbCmds = map[string]bbCmd{}

var defaultCmd *bbCmd

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
		return ErrNotRegistered
	}
	cmd.init()
	cmd.main()
	os.Exit(0)
	// Unreachable.
	return nil
}
