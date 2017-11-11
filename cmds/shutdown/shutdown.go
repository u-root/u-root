// Copyright 2015-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// shutdown halts or reboots.
//
// Synopsis:
//     shutdown [-dryrun] [operation]
//
// Description:
//     shutdown will either do or simulate the operation.
//     current operations are reboot and halt.
//
// Options:
//     -dryrun:   do not do really do it.
package main

import (
	"flag"
	"log"
	"os"

	"golang.org/x/sys/unix"
)

var (
	dryrun  = flag.Bool("dryrun", false, "Do not do kexec system calls")
	op      = "reboot"
	opcodes = map[string]int{
		"halt":    unix.LINUX_REBOOT_CMD_POWER_OFF,
		"reboot":  unix.LINUX_REBOOT_CMD_RESTART,
		"suspend": unix.LINUX_REBOOT_CMD_SW_SUSPEND,
	}
)

func usage() {
	log.Fatalf("shutdown [-dryrun] [halt|reboot|suspend] (defaults to reboot)")
}

func main() {
	flag.Parse()
	switch len(flag.Args()) {
	default:
		usage()
	case 0:
	case 1:
		op = flag.Args()[0]
	}

	f, ok := opcodes[op]
	if !ok {
		usage()
	}

	if *dryrun {
		log.Printf("unix.Reboot(0x%x)", f)
		os.Exit(0)
	}

	if err := unix.Reboot(f); err != nil {
		log.Fatalf(err.Error())
	}
}
