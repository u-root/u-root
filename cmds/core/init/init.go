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
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"syscall"

	"github.com/u-root/u-root/pkg/libinit"
	"github.com/u-root/u-root/pkg/ulog"
	"github.com/u-root/u-root/pkg/uroot/util"
)

var (
	verbose  = flag.Bool("v", false, "print all build commands")
	test     = flag.Bool("test", false, "Test mode: don't try to set control tty")
	debug    = func(string, ...interface{}) {}
	osInitGo = func() {}
	cmdList  []string
	cmdCount int
	envs     []string
)

func init() {
	r := util.UrootPath
	cmdList = []string{
		r("/inito"),

		r("/bbin/uinit"),
		r("/bin/uinit"),
		r("/buildbin/uinit"),

		r("/bin/defaultsh"),
		r("/bin/sh"),
	}
}

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
	// Create the root file systems.
	libinit.CreateRootfs()

	envs = os.Environ()
	debug("envs %v", envs)

	// Start background build.
	if isBgBuildEnabled() {
		go startBgBuild()
	}

	osInitGo()

	for _, v := range cmdList {
		debug("Trying to run %v", v)
		if _, err := os.Stat(v); os.IsNotExist(err) {
			debug("%v", err)
			continue
		}

		// I *love* special cases. Evaluate just the top-most symlink.
		//
		// In source mode, this would be a symlink like
		// /buildbin/defaultsh -> /buildbin/elvish ->
		// /buildbin/installcommand.
		//
		// To actually get the command to build, argv[0] has to end
		// with /elvish, so we resolve one level of symlink.
		if path.Base(v) == "defaultsh" {
			s, err := os.Readlink(v)
			if err == nil {
				v = s
			}
			debug("readlink of %v returns %v", v, s)
			// and, well, it might be a relative link.
			// We must go deeper.
			d, b := filepath.Split(v)
			d = filepath.Base(d)
			v = filepath.Join("/", os.Getenv("UROOT_ROOT"), d, b)
			debug("is now %v", v)
		}

		// inito is (optionally) created by the u-root command when the
		// u-root initramfs is merged with an existing initramfs that
		// has a /init. The name inito means "original /init" There may
		// be an inito if we are building on an existing initramfs. All
		// initos need their own pid space.
		var cloneFlags uintptr
		if v == "/inito" {
			cloneFlags = uintptr(syscall.CLONE_NEWPID)
		}

		cmdCount++
		cmd := exec.Command(v)
		cmd.Env = envs
		cmd.Stdin, cmd.Stdout, cmd.Stderr = os.Stdin, os.Stdout, os.Stderr
		if *test {
			cmd.SysProcAttr = &syscall.SysProcAttr{Cloneflags: cloneFlags}
		} else {
			cmd.SysProcAttr = &syscall.SysProcAttr{Setctty: true, Setsid: true, Cloneflags: cloneFlags}
		}
		debug("running %v", cmd)
		if err := cmd.Start(); err != nil {
			log.Printf("Error starting %v: %v", v, err)
			continue
		}
		for {
			var s syscall.WaitStatus
			var r syscall.Rusage
			if p, err := syscall.Wait4(-1, &s, 0, &r); p == cmd.Process.Pid {
				debug("Shell exited, exit status %d", s.ExitStatus())
				break
			} else if p != -1 {
				debug("Reaped PID %d, exit status %d", p, s.ExitStatus())
			} else {
				debug("Error from Wait4 for orphaned child: %v", err)
				break
			}
		}
		if err := cmd.Process.Release(); err != nil {
			log.Printf("Error releasing process %v: %v", v, err)
		}
	}
	if cmdCount == 0 {
		log.Printf("No suitable executable found in %+v", cmdList)
	}

	// We need to reap all children before exiting.
	log.Printf("Waiting for orphaned children")
	for {
		var s syscall.WaitStatus
		var r syscall.Rusage
		p, err := syscall.Wait4(-1, &s, 0, &r)
		if p == -1 {
			break
		}
		log.Printf("%v: exited with %v, status %v, rusage %v", p, err, s, r)
	}
	log.Printf("All commands exited")
	log.Printf("Syncing filesystems")
	syscall.Sync()
	log.Printf("Exiting...")
}
