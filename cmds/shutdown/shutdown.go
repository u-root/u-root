// Copyright 2015-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// shutdown halts, suspends, or reboots.
//
// Synopsis:
//     shutdown [-h|-r|halt|reboot|suspend]
//
// Description:
//     current operations are reboot (-r), suspend, and halt [-h].
//
// Options:
//     -r|reboot:	reboot the machine.
//     -h|halt:		halt the machine.
//     -s|suspend:	suspend the machine.
package main

import (
	"log"
	"os"

	"golang.org/x/sys/unix"
)

var (
	opcodes = map[string]uint{
		"halt":    unix.LINUX_REBOOT_CMD_POWER_OFF,
		"-h":      unix.LINUX_REBOOT_CMD_POWER_OFF,
		"reboot":  unix.LINUX_REBOOT_CMD_RESTART,
		"-r":      unix.LINUX_REBOOT_CMD_RESTART,
		"suspend": unix.LINUX_REBOOT_CMD_SW_SUSPEND,
		"-s":      unix.LINUX_REBOOT_CMD_SW_SUSPEND,
	}
	reboot = unix.Reboot
)

func usage() {
	log.Fatalf("shutdown [-h|-r|-s|halt|reboot|suspend] (defaults to halt)")
}

func main() {
	if len(os.Args) == 1 {
		os.Args = append(os.Args, "halt")
	}
	op, ok := opcodes[os.Args[1]]
	if !ok || len(os.Args) > 2 {
		usage()
	}
	if err := reboot(int(op)); err != nil {
		log.Fatalf(err.Error())
	}
}
