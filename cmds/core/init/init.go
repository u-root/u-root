// Copyright 2012-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// assumptions
// we've been booted into a ramfs with all this stuff unpacked and ready.
// we don't need a loop device mount because it's all there.
// So we run /go/bin/go build installcommand
// and then exec /buildbin/sh

package main

import (
	"flag"
	"fmt"
	"log"
	"os/exec"
	"syscall"

	"github.com/u-root/u-root/pkg/cmdline"
	"github.com/u-root/u-root/pkg/libinit"
	"github.com/u-root/u-root/pkg/ulog"
)

var (
	verbose  = flag.Bool("v", false, "print all build commands")
	test     = flag.Bool("test", false, "Test mode: don't try to set control tty")
	debug    = func(string, ...interface{}) {}
	osInitGo = func() {}
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
	// are still accessible through dmesg.
	if !*verbose {
		// Only messages more severe than "notice" are printed.
		if err := ulog.KernelLog.SetConsoleLogLevel(ulog.KLogNotice); err != nil {
			log.Printf("Could not set log level: %v", err)
		}
	}

	libinit.SetEnv()
	libinit.CreateRootfs()
	libinit.NetInit()

	// Potentially exec systemd if we have been asked to.
	osInitGo()

	// Start background build.
	if isBgBuildEnabled() {
		go startBgBuild()
	}

	// Turn off job control when test mode is on.
	ctty := libinit.WithTTYControl(!*test)

	// Allows passing args to uinit via kernel parameters, for example:
	//
	// uroot.uinitargs="-v --foobar"
	uinitArgs := libinit.WithArguments(cmdline.GetUinitArgs()...)

	cmdList := []*exec.Cmd{
		// inito is (optionally) created by the u-root command when the
		// u-root initramfs is merged with an existing initramfs that
		// has a /init. The name inito means "original /init" There may
		// be an inito if we are building on an existing initramfs. All
		// initos need their own pid space.
		libinit.Command("/inito", libinit.WithCloneFlags(syscall.CLONE_NEWPID), ctty),

		libinit.Command("/bbin/uinit", ctty, uinitArgs),
		libinit.Command("/bin/uinit", ctty, uinitArgs),
		libinit.Command("/buildbin/uinit", ctty, uinitArgs),

		libinit.Command("/bin/defaultsh", ctty),
		libinit.Command("/bin/sh", ctty),
	}

	cmdCount := libinit.RunCommands(debug, cmdList...)
	if cmdCount == 0 {
		log.Printf("No suitable executable found in %v", cmdList)
	}

	// We need to reap all children before exiting.
	log.Printf("Waiting for orphaned children")
	libinit.WaitOrphans()
	log.Printf("All commands exited")
	log.Printf("Syncing filesystems")
	syscall.Sync()
	log.Printf("Exiting...")
}
