// Copyright 2016 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Disassociate parts of the process execution context
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
		a = []string{"/bin/bash", "bash"}
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
