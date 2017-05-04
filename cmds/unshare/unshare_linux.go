// Copyright 2016-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//
// Options:
//     -ipc:           Unshare the IPC namespace
//     -mount:         Unshare the mount namespace
//     -pid:           Unshare the pid namespace
//     -net:           Unshare the net namespace
//     -uts:           Unshare the uts namespace
//     -user:          Unshare the user namespace
//     -map-root-user  Map current uid to root. Not working
//     -chroot         Chroot to the place specified by the argument
package main

import (
	"flag"
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
	chroot  = flag.String("chroot", "", "Chroot to the argument before running")
)

func setSysProcAttr(c *exec.Cmd) {
	c.SysProcAttr = &syscall.SysProcAttr{}
	if *mount {
		c.SysProcAttr.Unshareflags |= syscall.CLONE_NEWNS
	}
	if *uts {
		c.SysProcAttr.Unshareflags |= syscall.CLONE_NEWUTS
	}
	if *ipc {
		c.SysProcAttr.Unshareflags |= syscall.CLONE_NEWIPC
	}
	if *net {
		c.SysProcAttr.Unshareflags |= syscall.CLONE_NEWNET
	}
	if *pid {
		c.SysProcAttr.Unshareflags |= syscall.CLONE_NEWPID
	}
	if *user {
		c.SysProcAttr.Unshareflags |= syscall.CLONE_NEWUSER
	}

	if *chroot != "" {
		c.SysProcAttr.Chroot = *chroot
	}
}
