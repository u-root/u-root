// Copyright 2015-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// shutdown halts, suspends, or reboots at a specified time, or immediately.
//
// Synopsis:
//     shutdown [<-h|-r|-s|halt|reboot|suspend> [time [message...]]]
//
// Description:
//     current operations are reboot (-r), suspend, and halt [-h].
//     If no operation is specified halt is assumed.
//     If a time is given, an opcode is not optional.
//
// Options:
//     -r|reboot:	reboot the machine.
//     -h|halt:		halt the machine.
//     -s|suspend:	suspend the machine.
//
// Time is specified as "now", +minutes, or RFC3339 format.
// All other arguments past time are printed as a message.
// This could be used, for example, as input to goexpect.
package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

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
	delay  = time.Sleep
)

func usage() {
	log.Fatalf("shutdown [<-h|-r|-s|halt|reboot|suspend> [time [message...]]]")
}

func main() {
	a := os.Args
	if len(a) == 1 {
		a = append(a, "halt")
	}
	op, ok := opcodes[a[1]]
	if !ok {
		usage()
	}
	if len(a) < 3 {
		a = append(a, "now")
	}
	when := time.Now()
	switch {
	case a[2] == "now":
	case a[2][0] == '+':
		m, err := time.ParseDuration(a[2][1:] + "m")
		if err != nil {
			log.Fatal(err)
		}
		when = when.Add(m)
	default:
		t, err := time.Parse(time.RFC3339, a[2])
		if err != nil {
			log.Fatal(err)
		}
		when = t
	}

	delay(when.Sub(time.Now()))
	if len(a) > 2 {
		fmt.Println(strings.Join(a[3:], " "))
	}
	if err := reboot(int(op)); err != nil {
		log.Fatal(err)
	}
}
