// Copyright 2016-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Disassociate parts of the process execution context.
//
// Synopsis:
//     unshare [OPTIONS] [PROGRAM [ARGS]...]
//
// Description:
//     Go applications use multiple processes, and the Go user level scheduler
//     schedules goroutines onto those processes. For this reason, it is not
//     possible to use syscall.Unshare. A goroutine can call `syscall.Unshare`
//     from process m and the scheduler can resume that goroutine in process n,
//     which has not had the unshare operation! This is a known problem with
//     any system call that modifies the name space or file system context of
//     only one process as opposed to the entire Go application, i.e. all of
//     its processes. Examples include chroot and unshare. There has been
//     lively discussion of this problem but no resolution as of yet. In sum:
//     it is not possible to use `syscall.Unshare` from Go with any reasonable
//     expectation of success.
//
//     If PROGRAM is not specified, unshare defaults to /ubin/rush.
//
// Options:
//     -ipc:           Unshare the IPC namespace
//     -mount:         Unshare the mount namespace
//     -pid:           Unshare the pid namespace
//     -net:           Unshare the net namespace
//     -uts:           Unshare the uts namespace
//     -user:          Unshare the user namespace
//     -map-root-user: Map current uid to root. Not working
package main

import (
	"flag"
	"log"
	"os"
	"os/exec"
	"syscall"
)

var (
	ipc     = flag.Bool("ipc", false, "Unshare the IPC namespace")
	mount   = flag.Bool("mount", false, "Unshare the mount namespace")
	pid     = flag.Bool("pid", false, "Unshare the pid namespace")
	net     = flag.Bool("net", false, "Unshare the net namespace")
	uts     = flag.Bool("uts", false, "Unshare the uts namespace")
	user    = flag.Bool("user", false, "Unshare the user namespace")
	maproot = flag.Bool("map-root-user", false, "Map current uid to root. Not working")
)

func main() {
	flag.Parse()

	a := flag.Args()
	if len(a) == 0 {
		a = []string{"/ubin/rush", "rush"}
	}

	c := exec.Command(a[0], a[1:]...)
	c.SysProcAttr = &syscall.SysProcAttr{}
	if *mount {
		c.SysProcAttr.Cloneflags |= syscall.CLONE_NEWNS
	}
	if *uts {
		c.SysProcAttr.Cloneflags |= syscall.CLONE_NEWUTS
	}
	if *ipc {
		c.SysProcAttr.Cloneflags |= syscall.CLONE_NEWIPC
	}
	if *net {
		c.SysProcAttr.Cloneflags |= syscall.CLONE_NEWNET
	}
	if *pid {
		c.SysProcAttr.Cloneflags |= syscall.CLONE_NEWPID
	}
	if *user {
		c.SysProcAttr.Cloneflags |= syscall.CLONE_NEWUSER
	}

	c.Stdin = os.Stdin
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr

	if err := c.Run(); err != nil {
		log.Printf(err.Error())
	}
}
