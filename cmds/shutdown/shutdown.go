// Copyright 2015-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// shutdown halts or reboots.
//
// Synopsis:
//     shutdown [-dryrun] [-k] [-h|-r|operation]
//
// Description:
//     shutdown will either do or simulate the operation.
//     current operations are reboot and halt.
//
// Options:
//     -dryrun:	do not do really do it.
//     -k:	do not do really do it.
//     -r:	reboot the machine.
//     -h:	halt the machine.
package main

import (
	"errors"
	"flag"
	"log"
	"fmt"
	"os"

	"golang.org/x/sys/unix"
)

type options struct {
	k      bool
	h      bool
	r      bool
}

var (
	flags   options
	op      string
	opcodes = map[string]uint32{
		"halt":    unix.LINUX_REBOOT_CMD_POWER_OFF,
		"reboot":  unix.LINUX_REBOOT_CMD_RESTART,
		"suspend": unix.LINUX_REBOOT_CMD_SW_SUSPEND,
	}
)

func usage() {
	log.Fatalf("shutdown [-dryrun] [-k] [-h|-r|halt|reboot|suspend] (defaults to reboot)")
}

func init() {
	flag.BoolVar(&flags.k, "dryrun", false, "Do not do kexec system calls")
	flag.BoolVar(&flags.k, "k", false, "Do not do kexec system calls")
	flag.BoolVar(&flags.h, "h", false, "Halt the machine")
	flag.BoolVar(&flags.r, "r", false, "Reboot the machine")
}

func ascertainOperation() (uint32, error) {
	switch len(flag.Args()) {
	default:
		usage()
	case 0:
	case 1:
		op = flag.Args()[0]
	}

	if op == "" {
		if flags.h && flags.r {
			return 0, errors.New("Multiple flags detected")
		} else if flags.h {
			op = "halt"
		} else {
			op = "reboot"
		}

	} else {
		if flags.h || flags.r {
			return 0, errors.New("Please only provide flags or operation codes.")
		}
	}

	c, ok := opcodes[op]
	if !ok {
		return 0, errors.New(fmt.Sprintf("Invalid operation %s", op))
	}
	return c, nil
}

func main() {
	flag.Parse()

	f, err := ascertainOperation()
	if err != nil {
		log.Println(err)
		usage()
	}

	if flags.k {
		log.Printf("unix.Reboot(0x%x)", f)
		os.Exit(0)
	}

	if err := unix.Reboot(int(f)); err != nil {
		log.Fatalf(err.Error())
	}
}
