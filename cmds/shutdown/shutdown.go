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
	"syscall"
)

var (
	dryrun  = flag.Bool("dryrun", false, "Do not do kexec system calls")
	op      = "reboot"
	opcodes = map[string]uintptr{
		"halt":    syscall.LINUX_REBOOT_CMD_POWER_OFF,
		"reboot":  syscall.LINUX_REBOOT_CMD_RESTART,
		"suspend": syscall.LINUX_REBOOT_CMD_SW_SUSPEND,
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
		log.Printf("syscall.Syscall6(0x%x, 0x%x, 0x%x, 0x%x, 0, 0, 0)", syscall.SYS_REBOOT, syscall.LINUX_REBOOT_MAGIC1, syscall.LINUX_REBOOT_MAGIC2, f)
		os.Exit(0)
	}
	if e1, e2, err := syscall.Syscall6(syscall.SYS_REBOOT, syscall.LINUX_REBOOT_MAGIC1, syscall.LINUX_REBOOT_MAGIC2, f, 0, 0, 0); err != 0 {
		log.Fatalf("a %v b %v err %v", e1, e2, err)
	}
}
