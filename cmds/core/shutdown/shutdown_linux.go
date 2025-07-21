// Copyright 2015-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// shutdown halts, suspends, or reboots at a specified time, or immediately.
//
// Synopsis:
//
//	shutdown [<-h|-r|-s|halt|reboot|suspend> [time [message...]]]
//
// Description:
//
//	current operations are reboot (-r), suspend, and halt [-h].
//	If no operation is specified halt is assumed.
//	If a time is given, an opcode is not optional.
//
// Options:
//
//	-r|reboot:	reboot the machine.
//	-h|halt:		halt the machine.
//	-s|suspend:	suspend the machine.
//
// Time is specified as "now", +minutes, or RFC3339 format.
// All other arguments past time are printed as a message.
// This could be used, for example, as input to goexpect.
package main

import (
	"errors"
	"log"
	"os"
	"time"

	"golang.org/x/sys/unix"
)

var errUsageMessage = errors.New("shutdown [<-h|-r|-s|halt|reboot|suspend> [time [message...]]]")

var opcodes = map[string]uint{
	"halt":    unix.LINUX_REBOOT_CMD_POWER_OFF,
	"-h":      unix.LINUX_REBOOT_CMD_POWER_OFF,
	"reboot":  unix.LINUX_REBOOT_CMD_RESTART,
	"-r":      unix.LINUX_REBOOT_CMD_RESTART,
	"suspend": unix.LINUX_REBOOT_CMD_SW_SUSPEND,
	"-s":      unix.LINUX_REBOOT_CMD_SW_SUSPEND,
}

// shutdown calls unix.Reboot, with the type of shutdown defined in args, currently
// halt, reboot, or suspend. A time may be specified as "now",
// a future time parseable by time.ParseDuration, or in
// RFC3339 format. If dryrun is chosen, shutdown returns the opcode it
// would have used and an error, if any.
func shutdown(dryrun bool, args ...string) (uint, error) {
	if len(args) == 0 {
		args = append(args, "halt")
	}
	op, ok := opcodes[args[0]]
	if !ok {
		return 0, errUsageMessage
	}
	if len(args) < 2 {
		args = append(args, "now")
	}
	when := time.Now()

	switch {
	case args[1] == "now":

	case args[1][0] == '+':
		m, err := time.ParseDuration(args[1][1:] + "m")
		if err != nil {
			return 0, err
		}
		when = when.Add(m)
	default:
		t, err := time.Parse(time.RFC3339, args[1])
		if err != nil {
			return 0, err
		}
		when = t
	}

	// TODO: broadcast args[2:]... via wall or a similar mechanism.
	if !dryrun {
		time.Sleep(time.Until(when))
	}
	if !dryrun {
		if err := unix.Reboot(int(op)); err != nil {
			return 0, err
		}
	}

	return op, nil
}

func main() {
	if _, err := shutdown(false, os.Args[1:]...); err != nil {
		log.Fatal(err)
	}
}
