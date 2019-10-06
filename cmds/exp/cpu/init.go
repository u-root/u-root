// Copyright 2018-2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// This is init code for the case that cpu finds itself as pid 1.
// This is duplicative of the real init, but we're implementing it
// as a duplicate so we can get some idea of:
// what an init package should have
// what an init interface should have
// So we take a bit of duplication now to better understand these
// things. We also assume for now this is a busybox environment.
// It is unusual (I guess?) for cpu to be an init in anything else.
// So far, the case for an init pkg is not as strong as I thought
// it might be.
package main

import (
	"log"
	"syscall"

	"github.com/u-root/u-root/pkg/libinit"
)

func cpuSetup() error {
	log.Printf("Welcome to Plan 9(tm)!")
	libinit.SetEnv()
	libinit.CreateRootfs()
	libinit.NetInit()
	return nil
}

func cpuDone(c chan uint) {
	// We need to reap all children before exiting.
	log.Printf("init: Waiting for orphaned children")
	procs := libinit.WaitOrphans()
	log.Printf("cpu: All commands exited")
	log.Printf("cpu: Syncing filesystems")
	syscall.Sync()
	c <- procs
}
