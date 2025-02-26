// Copyright 2012-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build !windows

// init is u-root's standard userspace init process.
//
// init is intended to be the first process run by the kernel when it boots up.
// init does some basic initialization (mount file systems, turn on loopback)
// and then tries to execute, in order, /inito, a uinit (either in /bin, /bbin,
// or /ubin), and then a shell (/bin/defaultsh and /bin/sh).
package main

import (
	"flag"
	"fmt"
	"log"
	"os/exec"

	"github.com/u-root/u-root/pkg/libinit"
)

// initCmds has all the bits needed to continue
// the init process after some initial setup.
type initCmds struct {
	cmds []*exec.Cmd
}

var (
	verbose = flag.Bool("v", false, "Enable libinit debugging (includes showing commands that are run)")
	test    = flag.Bool("test", false, "Test mode: don't try to set control tty")
	debug   = func(string, ...interface{}) {}
)

func main() {
	flag.Parse()

	log.Printf("Welcome to u-root!")
	fmt.Println(`                              _`)
	fmt.Println(`   _   _      _ __ ___   ___ | |_`)
	fmt.Println(`  | | | |____| '__/ _ \ / _ \| __|`)
	fmt.Println(`  | |_| |____| | | (_) | (_) | |_`)
	fmt.Println(`   \__,_|    |_|  \___/ \___/ \__|`)
	fmt.Println()

	log.SetPrefix("init: ")

	if *verbose {
		debug = log.Printf
	}

	// Before entering an interactive shell, decrease the loglevel because
	// spamming non-critical logs onto the shell frustrates users. The logs
	// are still accessible through kernel logs buffers (on most kernels).
	quiet()

	libinit.SetEnv()
	libinit.CreateRootfs()
	libinit.NetInit()

	// osInitGo wraps all the kernel-specific (i.e. non-portable) stuff.
	// It returns an initCmds struct derived from kernel-specific information
	// to be used in the rest of init.
	ic := osInitGo()

	cmdCount := libinit.RunCommands(debug, ic.cmds...)
	if cmdCount == 0 {
		log.Printf("No suitable executable found in %v", ic.cmds)
	}

	// We need to reap all children before exiting.
	log.Printf("Waiting for orphaned children")
	libinit.WaitOrphans()
	log.Printf("All commands exited")
	log.Printf("Syncing filesystems")
	if err := quiesce(); err != nil {
		log.Printf("%v", err)
	}
	log.Printf("Exiting...")
}
